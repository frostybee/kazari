package main

import (
	"fmt"
	"io"
	"runtime/debug"
)

// runVersion prints the module version from build info. Installs via go run
// or go install carry the resolved module version; source builds report
// devel.
func runVersion(stdout io.Writer) int {
	version := "devel"
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		version = bi.Main.Version
	}
	fmt.Fprintf(stdout, "kazari %s\n", version)
	return 0
}
