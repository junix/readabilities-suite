// readabilities-suite is an independent black-box readability harness.
package main

import (
	"os"

	"github.com/junix/readabilities-suite/internal/cli"
)

var (
	version    = "0.1.0"
	sourcePath = "(unknown)"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, cli.BuildInfo{
		Version: version, SourcePath: sourcePath,
	}))
}
