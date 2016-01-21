package sourcer

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mikkeloscar/maze-ci/workspace"
)

// Git sourcer
type Git struct {
	// Workspace
	Workspace *workspace.Workspace
	// URL is the repository url
	URL string
	// Name is the name of the repository
	Name string
	// TODO implement authentication for private repositories
	// auth      *authentication.Auth
}

// get repo path
func (g *Git) src() string {
	return path.Join(g.Workspace.SrcDir, g.Name)
}

// Recursive search from repoPath to find directories with a PKGBUILD file
// The search is limited to a directory depth of 10
func findPKGBUILDs(repoPath string, repoName string, depth int) ([]*SrcPkg, error) {
	var srcPkg *SrcPkg
	var filePath string

	pkgs := []*SrcPkg{}

	files, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(files))

	for _, file := range files {
		filePath = path.Join(repoPath, file.Name())
		if file.IsDir() && file.Name()[0] != '.' {
			dirs = append(dirs, filePath)
		}

		if file.Name() == "PKGBUILD" {
			if depth < 2 {
				pkgb, err := pkgbuild.ParsePKGBUILD(filePath)
				if err != nil {
					return nil, err
				}
				srcPkg = &SrcPkg{
					PKGBUILD: pkgb,
					Path:     path.Dir(filePath),
				}
				pkgs = append(pkgs, srcPkg)
				return pkgs, nil
			}
		}
	}

	for _, dir := range dirs {
		pkgbs, err := findPKGBUILDs(dir, repoName, depth+1)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, pkgbs...)
	}

	return pkgs, nil
}

// Get PkgSrcs from git repo
// func (g *Git) Get(pkgs map[string]*pkg.RepoPkg) ([]*SrcPkg, error) {
// 	err := g.getRepo()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return findPKGBUILDs(g.src(), g.Name, 0)
// }

// TODO implement a version that can pull instead of just clean+clone
func (g *Git) getRepo() error {
	err := g.Workspace.CleanSrc()
	if err != nil {
		return err
	}

	err = gitClone(g.URL, g.src())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

// Clone git repository at url to dst
func gitClone(url, dst string) error {
	cmd := exec.Command("git", "clone", url, dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	return nil
}
