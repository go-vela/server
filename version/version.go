// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/version"
)

var (
	// Arch represents the architecture information for the package.
	Arch = runtime.GOARCH
	// Commit represents the git commit information for the package.
	Commit string
	// Compiler represents the compiler information for the package.
	Compiler = runtime.Compiler
	// Date represents the build date information for the package.
	Date string
	// Go represents the golang version information for the package.
	Go = runtime.Version()
	// OS represents the operating system information for the package.
	OS = runtime.GOOS
	// Tag represents the git tag information for the package.
	Tag string
)

// New creates a new version object for Vela that is used throughout the application.
func New() *version.Version {
	// check if a semantic tag was provided
	if len(Tag) == 0 {
		logrus.Warning("no semantic tag provided - defaulting to v0.0.0")

		// set a fallback default for the tag
		Tag = "v0.0.0"
	}

	v, err := semver.NewVersion(Tag)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to parse semantic version for %s: %w", Tag, err))
	}

	return &version.Version{
		Canonical:  Tag,
		Major:      v.Major(),
		Minor:      v.Minor(),
		Patch:      v.Patch(),
		PreRelease: v.Prerelease(),
		Metadata: version.Metadata{
			Architecture:    Arch,
			BuildDate:       Date,
			Compiler:        Compiler,
			GitCommit:       Commit,
			GoVersion:       Go,
			OperatingSystem: OS,
		},
	}
}
