package builder

import (
	"os"
	"path"
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
	dbPath             = "/home/vagrant/maze_repo/test.db.tar.gz"
	mirror             = "https://ftp.myrveln.se/pub/linux/archlinux/$repo/os/$arch"
)

func TestBuildPkgs(t *testing.T) {
	// setup
	err := os.MkdirAll(workdir, 0755)
	assert.NoError(t, err, "should not fail")
	err = os.MkdirAll(path.Dir(dbPath), 0755)
	assert.NoError(t, err, "should not fail")

	repo := repository.LocalRepo{
		RepoName: "TestRepo",
		Db:       dbPath,
		Deps:     nil,
	}

	repoPkg := &sourcer.RepoPkg{
		Repo: &repo,
		Pkg: sourcer.BuildPkg{
			Name:    "neovim-git",
			Sourcer: sourcer.AUR{},
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
		workspace: &workspace.Workspace{
			SrcDir: workdir,
		},
		user:     "maze",
		mirror:   mirror,
		notifier: notifier.StdoutNotifier(struct{}{}),
	}

	err = builder.Build(repoPkg)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.RemoveAll(workdir)
	assert.NoError(t, err, "should not fail")
	err = os.RemoveAll(path.Dir(dbPath))
	assert.NoError(t, err, "should not fail")
}
