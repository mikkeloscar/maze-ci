package sourcer

import (
	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mikkeloscar/maze-ci/repository"
)

// RepoPkg describes a built package in a repository
// type RepoPkg struct {
// 	Archs   []pkgbuild.Arch
// 	Version *pkgbuild.CompleteVersion
// 	Name    string
// }

// SrcPkg describes a package source including its pkgbuild and the path of the
// surrounding directory holding the PKGBUILD file and possible other files
type SrcPkg struct {
	PKGBUILD *pkgbuild.PKGBUILD
	Path     string
}

type BuildPkg struct {
	Name    string
	Sourcer Sourcer
}

type RepoPkg struct {
	Repo repository.Repo
	Pkg  BuildPkg
}

func (rp *RepoPkg) GetUpdated(pkgs []*SrcPkg) ([]*SrcPkg, error) {
	updated := make([]*SrcPkg, 0, len(pkgs))
	for _, pkg := range pkgs {
		_, new, err := rp.Repo.IsNew(pkg.PKGBUILD.Pkgbase, pkg.PKGBUILD.CompleteVersion())
		if err != nil {
			return nil, err
		}

		if new {
			updated = append(updated, pkg)
		}
	}

	return TopologicalSort(updated)
}
