package sourcer

import "github.com/mikkeloscar/maze-ci/repository"

// Sourcer can get PKGBUILDs from different sources
type Sourcer interface {
	// Get sources
	Get(pkg string, repo repository.Repo) ([]*SrcPkg, error)
}
