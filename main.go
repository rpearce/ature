// Command ature converts temperatures between Celsius, Fahrenheit, and
// Kelvin.
package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rpearce/ature/internal/cli"
)

// version is set at release time by goreleaser via
// -ldflags "-X main.version=...".
var version string

func main() {
	if err := cli.NewRootCmd(buildVersion()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errorLine(err))
		os.Exit(1)
	}
}

// buildVersion reports this binary's version from the ldflags value and the
// build information the Go toolchain embeds.
func buildVersion() string {
	info, _ := debug.ReadBuildInfo()
	return resolveVersion(version, info)
}

// resolveVersion picks the ldflags value, then the module version stamped by
// the Go toolchain for "go install ...@vX.Y.Z" builds, then "dev".
func resolveVersion(ldflag string, info *debug.BuildInfo) string {
	if ldflag != "" {
		return ldflag
	}
	if info != nil && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

// errorLine formats err for stderr. Cobra's unknown-command message already
// ends in a newline when it carries "Did you mean" suggestions; trimming it
// avoids printing a blank line after the error.
func errorLine(err error) string {
	return "ature: " + strings.TrimRight(err.Error(), "\n")
}
