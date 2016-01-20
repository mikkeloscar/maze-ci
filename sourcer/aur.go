package sourcer

import (
	"fmt"
	"path"

	"github.com/mikkeloscar/aur"
	"github.com/mikkeloscar/gopkgbuild"
	"github.com/mikkeloscar/maze-ci/repository"
	"github.com/mikkeloscar/maze-ci/workspace"
)

const AURCloneURL = "https://aur.archlinux.org/%s.git"

// AUR defines a sourcer for Arch User Repository
type AUR struct {
	// Workspace
	Workspace *workspace.Workspace
	// Pkgs      []string
}

// Get PKGBUILDs from AUR
// func (a AUR) Get(pkgs map[string]*pkg.RepoPkg) ([]*SrcPkg, error) {
// 	updates := make(map[string]struct{})
// 	err := a.getUpdates(a.Pkgs, pkgs, updates)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = a.getSourceRepos(updates)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var srcPkg *SrcPkg
// 	var filePath string

// 	srcPkgs := make([]*SrcPkg, 0, len(updates))

// 	// get a list of PKGBUILDs/SrcPkgs
// 	for u, _ := range updates {
// 		filePath = path.Join(a.Workspace.SrcDir, u, ".SRCINFO")

// 		pkgb, err := pkgbuild.ParseSRCINFO(filePath)
// 		if err != nil {
// 			return nil, err
// 		}
// 		srcPkg = &SrcPkg{
// 			PKGBUILD: pkgb,
// 			Path:     path.Join(a.Workspace.SrcDir, u),
// 		}
// 		srcPkgs = append(srcPkgs, srcPkg)
// 	}

// 	return srcPkgs, nil
// }

// Get PKGBUILDs from AUR
func (a AUR) Get(pkg string, repo repository.Repo) ([]*SrcPkg, error) {
	updates := make(map[string]struct{})
	err := a.getUpdates([]string{pkg}, repo, updates)
	if err != nil {
		return nil, err
	}

	err = a.getSourceRepos(updates)
	if err != nil {
		return nil, err
	}

	var srcPkg *SrcPkg
	var filePath string

	srcPkgs := make([]*SrcPkg, 0, len(updates))

	// get a list of PKGBUILDs/SrcPkgs
	for u, _ := range updates {
		filePath = path.Join(a.Workspace.SrcDir, u, ".SRCINFO")

		pkgb, err := pkgbuild.ParseSRCINFO(filePath)
		if err != nil {
			fmt.Println("HELLO")
			return nil, err
		}
		srcPkg = &SrcPkg{
			PKGBUILD: pkgb,
			Path:     path.Join(a.Workspace.SrcDir, u),
		}
		srcPkgs = append(srcPkgs, srcPkg)
	}

	return srcPkgs, nil
}

// GetSource fetches source from AUR and returns pkg info.
func (a AUR) GetSource(pkg string) (*SrcPkg, error) {
	err := a.updatePkgSrc(pkg)
	if err != nil {
		return nil, err
	}

	filePath := path.Join(a.Workspace.SrcDir, pkg, ".SRCINFO")

	pkgb, err := pkgbuild.ParseSRCINFO(filePath)
	if err != nil {
		return nil, err
	}
	srcPkg := &SrcPkg{
		PKGBUILD: pkgb,
		Path:     path.Join(a.Workspace.SrcDir, pkg),
	}

	return srcPkg, nil
}

// func (a AUR) GetDependencies(pkg string) ([]*BuildPkg, error) {
// 	pkgs := make(map[string]struct{})
// 	err := a.getDependencies([]string{pkg}, pkgs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	buildPkgs := make([]*BuildPkg, 0, len(pkgs))

// 	for pkg, _ := range pkgs {
// 		buildPkg := &BuildPkg{
// 			Name:    pkg,
// 			Sourcer: a,
// 		}
// 		buildPkgs = append(buildPkgs, buildPkg)
// 	}

// 	return buildPkgs, nil
// }

// func (a *AUR) getDependencies(pkgs []string, allPkgs map[string]struct{}) error {
// 	pkgsInfo, err := aur.Multiinfo(pkgs)
// 	if err != nil {
// 		return err
// 	}

// 	for _, pkg := range pkgsInfo {
// 		allPkgs[pkg.Name]

// 		// TODO handle checkdepends and maybe optdepends
// 		depends := make([]string, 0, len(pkg.Depends)+len(pkg.MakeDepends))
// 		depends = append(depends, pkg.Depends...)
// 		depends = append(depends, pkg.MakeDepends...)

// 		a.getDependencies(depends, allPkgs)
// 	}

// 	return nil
// }

// get list of packages with a new version in AUR including possible
// dependencies
// func (a *AUR) getUpdates(pkgs []string, currPkgs map[string]*pkg.RepoPkg, updates map[string]struct{}) error {
// 	pkgsInfo, err := aur.Multiinfo(pkgs)
// 	if err != nil {
// 		return err
// 	}

// 	for _, pkg := range pkgsInfo {
// 		if p, ok := currPkgs[pkg.Name]; ok {
// 			if p.Version.Older(pkg.Version) {
// 				updates[pkg.Name] = struct{}{}
// 			}
// 		} else {
// 			updates[pkg.Name] = struct{}{}
// 		}

// 		// TODO handle checkdepends and maybe optdepends
// 		depends := make([]string, 0, len(pkg.Depends)+len(pkg.MakeDepends))
// 		depends = append(depends, pkg.Depends...)
// 		depends = append(depends, pkg.MakeDepends...)

// 		a.getUpdates(depends, currPkgs, updates)
// 	}

// 	return nil
// }

// get list of packages with a new version in AUR including possible
// dependencies
func (a *AUR) getUpdates(pkgs []string, repo repository.Repo, updates map[string]struct{}) error {
	pkgsInfo, err := aur.Multiinfo(pkgs)
	if err != nil {
		return err
	}

	for _, pkg := range pkgsInfo {
		updates[pkg.Name] = struct{}{}

		// TODO handle checkdepends and maybe optdepends
		depends := make([]string, 0, len(pkg.Depends)+len(pkg.MakeDepends))
		depends = append(depends, pkg.Depends...)
		depends = append(depends, pkg.MakeDepends...)

		a.getUpdates(depends, repo, updates)
	}

	return nil
}

// get source repos from set of package names
func (a *AUR) getSourceRepos(pkgs map[string]struct{}) error {
	clone := make(chan error)

	for pkg, _ := range pkgs {
		go a.updateRepo(pkg, clone)
	}

	for range pkgs {
		// TODO grab error response
		<-clone
	}

	return nil
}

// update (clone or pull) AUR package repo
func (a *AUR) updatePkgSrc(pkg string) error {
	url := fmt.Sprintf(AURCloneURL, pkg)

	// TODO implement version that can pull instead of clone
	err := gitClone(url, a.Workspace.SrcDir)
	if err != nil {
		return err
	}

	return nil
}

// update (clone or pull) AUR package repo
func (a *AUR) updateRepo(pkg string, c chan<- error) {
	url := fmt.Sprintf(AURCloneURL, pkg)

	// TODO implement version that can pull instead of clone
	err := gitClone(url, path.Join(a.Workspace.SrcDir, pkg))
	c <- err
}
