package workspace

import (
	"os"
	"testing"
)

const wPath = "test_workspace"

// Test creating a workspace
func TestNewWorkspace(t *testing.T) {
	_, err := NewWorkspace(wPath)
	if err != nil {
		t.Error(err)
	}

	_, err = NewWorkspace(wPath)
	if err != nil {
		t.Error(err)
	}

	cleanup(t)

	// should fail because of permissions
	_, err = NewWorkspace("/")
	if err == nil {
		t.Errorf("expected error creating workspace at /")
	}
}

// Test cleaning directories
func TestClean(t *testing.T) {
	w, _ := NewWorkspace(wPath)
	err := w.Clean()
	if err != nil {
		t.Error(err)
	}

	cleanup(t)
}

// cleanup
func cleanup(t *testing.T) {
	err := os.RemoveAll(wPath)
	if err != nil {
		t.Error(err)
	}
}
