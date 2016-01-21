package builder

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/kr/pty"
	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mikkeloscar/maze-ci/notifier"
	"github.com/mikkeloscar/maze-ci/repository"
	"github.com/mikkeloscar/maze-ci/sourcer"
)

// NSpawn is a systemd-nspawn based builder.
type NSpawn struct {
	// current build ID.
	id uint32
	// Arch is the architecture for which this builder can build.
	arch pkgbuild.Arch
	// path to container template.
	template string
	// path to container (TODO better container management).
	path string
	// path to pacman config template.
	pacmanConfTemplate string
	// path to possible modified copy of pacman config.
	currentPacmanConf string
	// build path in container.
	buildPath string
	// build user.
	user string

	// current pkg being build.
	pkg *sourcer.SrcPkg

	notifier notifier.Notifier

	channels struct {
		stop     chan struct{}
		killResp chan error
	}
}

// Arch return the build architecture supported by the builder.
func (n *NSpawn) Arch() pkgbuild.Arch {
	return n.arch
}

// Ready returns true if the builder is ready to start a build.
func (n *NSpawn) Ready() bool {
	return n.id == 0
}

// SetNotifier sets the notifier to use for notifying about build status.
func (n *NSpawn) SetNotifier(notifier notifier.Notifier) {
	n.notifier = notifier
}

// update nspawn container.
func (n *NSpawn) update(repo repository.Repo) error {
	err := n.setupPacmanConf(repo)
	if err != nil {
		return err
	}

	dir, file := path.Split(n.currentPacmanConf)

	// replace default pacman.conf
	err = n.run(fmt.Sprintf("sudo cp %s/%s /etc/pacman.conf", n.buildPath, file), dir, false)
	if err != nil {
		return err
	}

	err = n.run("sudo pacman -Syu --noconfirm", dir, false)
	if err != nil {
		return err
	}

	return nil
}

// takedown container.
func (n *NSpawn) takedown() error {
	n.id = 0

	// remove temp pacman.conf
	err := os.Remove(n.currentPacmanConf)
	if err != nil {
		return err
	}

	cmd := exec.Command("sudo", "btrfs", "subvolume", "delete", n.path)
	return cmd.Run()
}

// Setup custom pacman.conf file.
func (n *NSpawn) setupPacmanConf(repo repository.Repo) error {
	content, err := ioutil.ReadFile(n.pacmanConfTemplate)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(n.currentPacmanConf, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "# :INSERT_REPO:" {
			w.WriteString(repository.ConfEntry(repo))
			continue
		}

		_, err = w.WriteString(line)
		if err != nil {
			return err
		}

		err = w.WriteByte('\n')
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func (n *NSpawn) Stop() error {
	n.channels.stop <- struct{}{}

	// TODO timeout
	err := <-n.channels.killResp

	return err
}

// Get a sorted list of packages to build.
func (n *NSpawn) getBuildPkgs(pkg *sourcer.RepoPkg) ([]*sourcer.SrcPkg, error) {
	pkgSrcs, err := pkg.Pkg.Sourcer.Get(pkg.Pkg.Name, pkg.Repo)
	if err != nil {
		return nil, err
	}

	// Get a list of devel packages (-{bzr,git,svn,hg}) where an extra
	// version check is needed.
	updates := make([]*sourcer.SrcPkg, 0, len(pkgSrcs))

	for _, pkgSrc := range pkgSrcs {
		if pkgSrc.PKGBUILD.IsDevel() {
			updates = append(updates, pkgSrc)
		}
	}

	err = n.updatePkgSrcs(updates)
	if err != nil {
		return nil, err
	}

	// TODO: update pkgs done
	n.notifier.Done(n.id)

	return pkg.GetUpdated(pkgSrcs)
}

func (n *NSpawn) updatePkgSrcs(pkgs []*sourcer.SrcPkg) error {
	for _, pkg := range pkgs {
		n.pkg = pkg // set current package
		_, err := n.updatePkgSrc(pkg)
		if err != nil {
			return err
		}
		n.notifier.PkgDone(n.id, n.pkg)
	}

	return nil
}

func (n *NSpawn) updatePkgSrc(pkg *sourcer.SrcPkg) (*sourcer.SrcPkg, error) {
	err := n.run(fmt.Sprintf("cd %s && makepkg -os --noconfirm && mksrcinfo", n.buildPath), pkg.Path, false)
	if err != nil {
		return nil, err
	}

	// update pkgbuild info
	filePath := path.Join(pkg.Path, ".SRCINFO")

	pkgb, err := pkgbuild.ParseSRCINFO(filePath)
	if err != nil {
		return nil, err
	}

	pkg.PKGBUILD = pkgb

	return pkg, nil
}

func (n *NSpawn) syncBuild(pkg *sourcer.RepoPkg) ([]string, error) {
	// update environment
	err := n.update(pkg.Repo)
	if err != nil {
		// TODO mark as internal failure and log to local builder logs
		return nil, err
		// n.notifier.Failed(n.id, err)
		// goto takedown
	}

	pkgs, err := n.buildPkgs(pkg)
	if err != nil {
		return nil, err
		// n.notifier.Failed(n.id, err)
		// goto takedown
	}

	err = n.takedown()
	if err != nil {
		// TODO log error in takedown
		return nil, err
	}

	return pkgs, nil
}

func (n *NSpawn) Build(pkg *sourcer.RepoPkg) {
	_, err := n.syncBuild(pkg)
	if err != nil {
		// TODO: better handling of internal vs build error
		n.notifier.Failed(n.id, err)
		return
	}

	// TODO: Upload packages
	// err = n.repo.Upload(n.id, buildPkgs)
	// if err != nil {
	// 	n.notifier.UploadFailed(n.id, err)
	// }
}

func (n *NSpawn) buildPkgs(pkg *sourcer.RepoPkg) ([]string, error) {
	pkgs, err := n.getBuildPkgs(pkg)
	if err != nil {
		return nil, err
	}

	buildPkgs := make([]string, 0, len(pkgs))

	for _, pkg := range pkgs {
		n.pkg = pkg // set current package
		pkgPaths, err := n.build(pkg)
		if err != nil {
			return nil, err
		}
		n.notifier.PkgDone(n.id, n.pkg)

		buildPkgs = append(buildPkgs, pkgPaths...)
	}

	n.notifier.Done(n.id)

	return buildPkgs, nil
}

func (n *NSpawn) build(pkg *sourcer.SrcPkg) ([]string, error) {
	err := n.run(fmt.Sprintf("cd %s && makepkg -is --noconfirm", n.buildPath), pkg.Path, true)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(pkg.Path)
	if err != nil {
		return nil, err
	}

	pkgs := make([]string, 0, 1)

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "pkg.tar.xz") {
			pkgPath := path.Join(pkg.Path, f.Name())
			pkgs = append(pkgs, pkgPath)
		}
	}

	return pkgs, nil
}

// run command in container as build user while the workdir is mounted under
// buildPath.
// TODO handle stdout optional
func (n *NSpawn) run(command string, workdir string, output bool) error {
	cmd := exec.Command("sudo", "systemd-nspawn",
		"--quiet",
		fmt.Sprintf("--template=%s", n.template),
		"-u", n.user,
		fmt.Sprintf("--bind=%s:%s", workdir, n.buildPath),
		"-D", n.path,
		"sh", "-c", command)

	tty, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(tty)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 	select {
			// 	case <-n.channels.stop:
			// 		err := cmd.Process.Kill()
			// 		// send result of killing process
			// 		n.channels.killResp <- err
			// 		return
			// 	default:
			// if output {
			if scanner.Scan() {
				stdout := scanner.Text()
				if output {
					n.notifier.BuildOutput(n.id, n.pkg, stdout)
				} else {
					fmt.Printf("container: %s\n", stdout)
				}
			} else {
				// no more output
				return
			}
		}
		// }
		// }
	}()

	// if err != nil {
	// 	if exitErr, ok := err.(*exec.ExitError); ok {
	// 		if _, ok := exitErr.Sys().(syscall.WaitStatus); ok {
	// 			return err
	// 		}
	// 	}
	// 	return err
	// }

	err = cmd.Wait()
	if err != nil {
		return err
	}

	wg.Wait() // wait until all output has been sent

	return nil
}
