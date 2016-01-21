package workspace

import (
	"io/ioutil"
	"os"
	"path"
)

// Workspace defines a workspace for keeping sources
type Workspace struct {
	// SrcDir is the source directory of the workspace
	SrcDir string
	// TmpDir is the temporary directory of the workspace
	TmpDir string
}

// NewWorkspace creates a new workspace struct including the underlying file
// structure
func NewWorkspace(wPath string) (*Workspace, error) {
	w := Workspace{
		SrcDir: path.Join(wPath, "src"),
		TmpDir: path.Join(wPath, "tmp"),
	}

	err := os.MkdirAll(w.SrcDir, 0755)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(w.TmpDir, 0755)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

// Clean entire Workspace
func (w *Workspace) Clean() error {
	err := w.CleanSrc()
	if err != nil {
		return err
	}

	return w.CleanTmp()
}

// CleanSrc removes anything inside SrcDir but not SrcDir itself
func (w *Workspace) CleanSrc() error {
	return w.clean(w.SrcDir)
}

// CleanTmp removes anything inside TmpDir but not TmpDir itself
func (w *Workspace) CleanTmp() error {
	return w.clean(w.TmpDir)
}

// cleans everything within a directory excluding the directory itself
func (w *Workspace) clean(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		err = os.Remove(path.Join(dir, f.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}
