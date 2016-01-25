package repository

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/mikkeloscar/gopkgbuild"
)

type LocalRepo struct {
	// Name of repo.
	RepoName string
	// Path to repo db.
	Db string
	// repos which this repo depend on.
	Deps []Repo
}

// Url gets repo url/path to db file.
func (l *LocalRepo) Url() string {
	return fmt.Sprintf("file://%s", l.Db)
}

// Name gets name of repo.
func (l *LocalRepo) Name() string {
	return l.RepoName
}

// Dependencies gets list of repo dependencies.
func (l *LocalRepo) Dependencies() []Repo {
	return l.Deps
}

// Add adds package to repo.
// TODO: this is a very simple version, should be better.
func (l *LocalRepo) Add(pkgPath string, buildId uint32) error {
	repoPath, _ := path.Split(l.Db)
	pkgBasePath, pkg := path.Split(pkgPath)

	if repoPath != pkgBasePath {
		// copy pkg to repo path.
		err := copyFile(path.Join(repoPath, pkg), pkgPath)
		if err != nil {
			return err
		}
	}

	err := exec.Command("/usr/bin/repo-add", "-R", l.Db, pkgPath).Run()
	if err != nil {
		return err
	}

	return nil
}

// TODO implement
func (l *LocalRepo) Remove(name string) error {
	return nil
}

// IsNew returns true if pkg is a newer version than what's in the repo.
func (l *LocalRepo) IsNew(pkg string, version pkgbuild.CompleteVersion) (bool, bool, error) {
	for _, dep := range l.Deps {
		found, new, err := dep.IsNew(pkg, version)
		if err != nil {
			return false, false, err
		}

		if found && new {
			return found, new, nil
		}
	}

	return l.isNew(pkg, version)
}

// Check if package exists in repo and if the version is newer.
func (l *LocalRepo) isNew(pkg string, version pkgbuild.CompleteVersion) (bool, bool, error) {
	f, err := os.Open(l.Db)
	if err != nil {
		if os.IsNotExist(err) {
			return false, true, nil
		}
		return false, false, err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return false, false, err
	}

	tarR := tar.NewReader(gzf)

	for {
		header, err := tarR.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return false, false, err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			n, v := splitNameVersion(header.Name)
			if n == pkg {
				if version.Newer(v) {
					return true, true, nil
				}
				return true, false, nil
			}
		case tar.TypeReg:
			continue
		}
	}

	return false, true, nil
}

// turn "zlib-1.2.8-4/" into ("zlib", "1.2.8-4").
func splitNameVersion(str string) (string, string) {
	chars := strings.Split(str[:len(str)-1], "-")
	name := chars[:len(chars)-2]
	version := chars[len(chars)-2:]

	return strings.Join(name, "-"), strings.Join(version, "-")
}

// copy file from src to dst.
func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
