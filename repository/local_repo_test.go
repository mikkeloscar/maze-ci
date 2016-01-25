package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/mikkeloscar/gopkgbuild"
	"github.com/stretchr/testify/assert"
)

const (
	repo1Path = "test_repo/repo1.db.tar.gz"
	repo2Path = "test_repo/repo2.db"
)

var (
	repo1 = LocalRepo{
		RepoName: "repo1",
		Db:       repo1Path,
		Deps:     nil,
	}
	repo2 = LocalRepo{
		RepoName: "repo2",
		Db:       fmt.Sprintf("%s.tar.gz", repo2Path),
		Deps:     nil,
	}
)

func TestIsNew(t *testing.T) {
	pkg := "ca-certificates"

	// Check if existing package is new
	version, _ := pkgbuild.NewCompleteVersion("20150402-1")
	found, new, err := repo1.IsNew(pkg, *version)
	assert.NoError(t, err, "should not fail")
	assert.True(t, found, "should be true")
	assert.False(t, new, "should be false")

	// Check if new package is new
	version, _ = pkgbuild.NewCompleteVersion("20150402-2")
	found, new, err = repo1.IsNew(pkg, *version)
	assert.NoError(t, err, "should not fail")
	assert.True(t, found, "should be true")
	assert.True(t, new, "should be true")

	// Check if old package is new
	version, _ = pkgbuild.NewCompleteVersion("20150401-1")
	found, new, err = repo1.IsNew(pkg, *version)
	assert.NoError(t, err, "should not fail")
	assert.True(t, found, "should be true")
	assert.False(t, new, "should be false")

	// Check if existing package is new (repo is empty)
	version, _ = pkgbuild.NewCompleteVersion("20150402-1")
	found, new, err = repo2.IsNew(pkg, *version)
	assert.NoError(t, err, "should not fail")
	assert.False(t, found, "should be false")
	assert.True(t, new, "should be true")
}

func TestAdd(t *testing.T) {
	pkg := "test_repo/ca-certificates-20150402-1-any.pkg.tar.xz"

	// Add to empty repo
	err := repo2.Add(pkg, 0)
	assert.NoError(t, err, "should not fail")

	// Add existing pkg to repo
	err = repo2.Add(pkg, 0)
	assert.NoError(t, err, "should not fail")

	// cleanup
	err = os.Remove(repo2.Db)
	assert.NoError(t, err, "should not fail")
	err = os.Remove(repo2Path)
	assert.NoError(t, err, "should not fail")
	err = os.Remove(repo2.Db + ".old")
	assert.NoError(t, err, "should not fail")
}
