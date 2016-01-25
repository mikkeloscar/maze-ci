package sourcer

import (
	"github.com/mikkeloscar/maze-ci/repository"
	"github.com/mikkeloscar/maze-ci/workspace"
)

// Sourcer can get PKGBUILDs from different sources
type Sourcer interface {
	// Get sources
	Get(pkg string, repo repository.Repo, ws *workspace.Workspace) ([]*SrcPkg, error)
}
