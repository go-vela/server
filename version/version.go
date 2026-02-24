// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
	"github.com/sirupsen/logrus"
)

const versionFormat = `{
  Canonical: %s,
  Major: %d,
  Minor: %d,
  Patch: %d,
  PreRelease: %s,
  Metadata: {
    Architecture: %s,
    BuildDate: %s,
    Compiler: %s,
    GitCommit: %s,
    GoVersion: %s,
    OperatingSystem: %s,
  }
}`

// Version represents application information that
// follows semantic version guidelines from
// https://semver.org/.
//
// swagger:model Version
type Version struct {
	// Canonical represents a canonical semantic version for the application.
	Canonical string `json:"canonical"`
	// Major represents incompatible API changes.
	Major uint64 `json:"major"`
	// Minor represents added functionality in a backwards compatible manner.
	Minor uint64 `json:"minor"`
	// Patch represents backwards compatible bug fixes.
	Patch uint64 `json:"patch"`
	// PreRelease represents unstable changes that might not be compatible.
	PreRelease string `json:"pre_release,omitempty"`
	// Metadata represents extra information surrounding the application version.
	Metadata Metadata `json:"metadata,omitempty"`
}

// Meta implements a formatted string containing only metadata for the Version type.
func (v *Version) Meta() string {
	return v.Metadata.String()
}

// Semantic implements a formatted string containing a formal semantic version for the Version type.
func (v *Version) Semantic() string {
	return v.Canonical
}

// String implements the Stringer interface for the Version type.
func (v *Version) String() string {
	return fmt.Sprintf(
		versionFormat,
		v.Canonical,
		v.Major,
		v.Minor,
		v.Patch,
		v.PreRelease,
		v.Metadata.Architecture,
		v.Metadata.BuildDate,
		v.Metadata.Compiler,
		v.Metadata.GitCommit,
		v.Metadata.GoVersion,
		v.Metadata.OperatingSystem,
	)
}

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
func New() *Version {
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

	return &Version{
		Canonical:  Tag,
		Major:      v.Major(),
		Minor:      v.Minor(),
		Patch:      v.Patch(),
		PreRelease: v.Prerelease(),
		Metadata: Metadata{
			Architecture:    Arch,
			BuildDate:       Date,
			Compiler:        Compiler,
			GitCommit:       Commit,
			GoVersion:       Go,
			OperatingSystem: OS,
		},
	}
}
