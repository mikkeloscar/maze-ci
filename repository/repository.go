package repository

import "github.com/mikkeloscar/gopkgbuild"

type Repo interface {
	// Check if package is newer than the equivalent package in the repo.
	// returns found=false, new=true if no matching package was found in
	// repo.
	IsNew(pkg string, version pkgbuild.CompleteVersion) (bool, bool, error)
	// Add package to repo.
	Add(path string) error
	// Remove package from repo.
	Remove(pkg string) error
	// Get repo Name.
	Name() string
	// Get repo url.
	Url() string
	// Get repo dependencies.
	Dependencies() []Repo
}

// RepoPkg describes a built package in a repository
type RepoPkg struct {
	Archs   pkgbuild.Arch
	Version *pkgbuild.CompleteVersion
	Name    string
}

func ConfEntry(repo Repo) string {
	// entry := ""
	// for _, dep := range repo.Dependencies() {
	// 	entry += ConfEntry(dep)
	// }
	// return fmt.Sprintf("[%s]\nServer = %s\n\n", repo.Name(), repo.Url())

	return "\n"
}

// // A list of Repos
// type Repos []Repo

// // Check if package is a new version of a package in a list of repos.
// func (r Repos) IsNew(pkg, version string) (bool, error) {
// 	var err error
// 	new, found := false

// 	for _, repo := range r {
// 		found, err = repo.InRepo(pkg.Name)
// 		if err != nil {
// 			return err
// 		}

// 		if !found {
// 			continue
// 		}

// 		new, err = repo.IsNew(pkg.Name, pkg.Version)
// 		if err != nil {
// 			return err
// 		}

// 		if !new {
// 			continue
// 		}
// 	}

// 	return new || !found
// }

// // Repository is the central part of maze-ci. It knows how to handle package
// // sources with the help of sourcers as well as how to build packages and
// // serve them
// type Repository struct {
// 	//Workspace *sourcer.Workspace // TODO should not be part of sourcer pkg
// 	Sourcers []*sourcer.Sourcer
// 	// Packages is a map of name -> RepoPkg
// 	Packages map[string]*pkg.RepoPkg
// 	Name     string
// 	// Owner     *User
// 	Private bool
// 	Archs   map[string]pkgbuild.Arch
// }

// // TODO better error handling
// // update sources concurrently
// func (r *Repository) updateSources() ([]*pkg.SrcPkg, error) {
// 	var pkgs []*pkg.SrcPkg
// 	var tmpPkgs []*pkg.SrcPkg
// 	u := make(chan []*pkg.SrcPkg)

// 	for _, s := range r.Sourcers {
// 		go r.conSourceUpdate(s, u)
// 	}

// 	for range r.Sourcers {
// 		tmpPkgs = <-u
// 		if tmpPkgs == nil {
// 			return nil, fmt.Errorf("unable to get sources")
// 		}
// 		pkgs = append(pkgs, tmpPkgs...)
// 	}

// 	return pkgs, nil
// }

// // helper function for making updateSources() concurrent
// func (r *Repository) conSourceUpdate(sourcer *sourcer.Sourcer, update chan<- []*pkg.SrcPkg) {
// 	pkgs, err := (*sourcer).Get(r.Packages)
// 	if err != nil {
// 		update <- nil
// 	}

// 	update <- pkgs
// }

// func (r *Repository) getUpdates() ([]*sourcer.PkgSrc, error) {
// 	for _, s := range r.Sourcers {
// 		s.Get(r.)
// 	}
// }
