package notifier

import (
	"fmt"

	"github.com/mikkeloscar/maze-ci/sourcer"
)

type Notifier interface {
	Done(buildId uint32)
	PkgDone(buildId uint32, pkg *sourcer.SrcPkg)
	Failed(buildId uint32, err error)
	BuildOutput(buildId uint32, pkg *sourcer.SrcPkg, output string)
	AddPkgFailed(buildId uint32, pkg string, err error)
}

type StdoutNotifier struct{}

func (s StdoutNotifier) Done(buildId uint32) {
	fmt.Printf("Build: %d - Done\n", buildId)
}

func (s StdoutNotifier) PkgDone(buildId uint32, pkg *sourcer.SrcPkg) {
	fmt.Printf("Build: %d - Pkg build done for: %s\n", buildId, pkg.PKGBUILD.Pkgbase)
}

func (s StdoutNotifier) Failed(buildId uint32, err error) {
	fmt.Printf("Build: %d - Error: %s\n", buildId, err.Error())
}

func (s StdoutNotifier) BuildOutput(buildId uint32, pkg *sourcer.SrcPkg, output string) {
	fmt.Printf("Build: %d - %s: %s\n", buildId, pkg.PKGBUILD.Pkgbase, output)
}

func (s StdoutNotifier) AddPkgFailed(buildId uint32, pkg string, err error) {
	fmt.Printf("Build: %d - Upload of %s failed: %s\n", buildId, pkg, err.Error())
}
