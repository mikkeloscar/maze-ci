package builder

import (
	"fmt"
	"testing"

	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mikkeloscar/maze-ci/notifier"
	"github.com/mikkeloscar/maze-ci/repository"
	"github.com/mikkeloscar/maze-ci/sourcer"
	"github.com/mikkeloscar/maze-ci/workspace"
	"github.com/stretchr/testify/assert"
)

const (
	containerPath      = "/home/vagrant/arch-maze"
	templateContainer  = "/var/lib/maze/nspawn_template"
	workdir            = "/home/vagrant/maze_workdir"
	pacmanConfTemplate = "../../contrib/pacman.conf"
	pacmanConfPath     = "/home/vagrant/pacman_conf"
)

// func TestBuildPkgs(t *testing.T) {
// 	builder := NSpawn{
// 		arch:      pkgbuild.X8664,
// 		template:  templateContainer,
// 		path:      containerPath,
// 		buildPath: "/mnt/build",
// 		workspace: &workspace.Workspace{
// 			SrcDir: workdir,
// 		},
// 		user: "maze",
// 	}

// 	pkgs := []*pkg.SrcPkg{
// 		&pkg.SrcPkg{
// 			Path: workdir + "/" + "libtermkey-bzr",
// 		},
// 		&pkg.SrcPkg{
// 			Path: workdir + "/" + "libvterm-bzr",
// 		},
// 		&pkg.SrcPkg{
// 			Path: workdir + "/" + "neovim-git",
// 		},
// 	}

// 	packages, err := builder.Build(pkgs)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	assert.NoError(t, err, "should not fail")

// 	fmt.Println(packages)
// }

func TestBuildPkgs(t *testing.T) {
	repo := repository.LocalRepo{
		RepoName: "TestRepo",
		Db:       "none",
		Deps:     nil,
	}

	repoPkg := &sourcer.RepoPkg{
		Repo: &repo,
		Pkg: sourcer.BuildPkg{
			Name: "neovim-git",
			Sourcer: sourcer.AUR{
				Workspace: &workspace.Workspace{
					SrcDir: workdir,
				},
			},
		},
	}

	builder := NSpawn{
		id:                 1,
		arch:               pkgbuild.X8664,
		template:           templateContainer,
		path:               containerPath,
		buildPath:          "/mnt/build",
		pacmanConfTemplate: pacmanConfTemplate,
		currentPacmanConf:  pacmanConfPath,
		// workspace: &workspace.Workspace{
		// 	SrcDir: workdir,
		// },
		user:     "maze",
		notifier: notifier.StdoutNotifier(struct{}{}),
	}

	packages, err := builder.syncBuild(repoPkg)
	assert.NoError(t, err, "should not fail")

	fmt.Println(packages)
}
