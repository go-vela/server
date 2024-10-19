// SPDX-License-Identifier: Apache-2.0

package version

import "fmt"

const metaFormat = `{
  Architecture: %s,
  BuildDate: %s,
  Compiler: %s,
  GitCommit: %s,
  GoVersion: %s,
  OperatingSystem: %s,
}`

// Metadata represents extra information surrounding the application version.
type Metadata struct {
	// Architecture represents the architecture information for the application.
	Architecture string `json:"architecture,omitempty"`
	// BuildDate represents the build date information for the application.
	BuildDate string `json:"build_date,omitempty"`
	// Compiler represents the compiler information for the application.
	Compiler string `json:"compiler,omitempty"`
	// GitCommit represents the git commit information for the application.
	GitCommit string `json:"git_commit,omitempty"`
	// GoVersion represents the golang version information for the application.
	GoVersion string `json:"go_version,omitempty"`
	// OperatingSystem represents the operating system information for the application.
	OperatingSystem string `json:"operating_system,omitempty"`
}

// String implements the Stringer interface for the Metadata type.
func (m *Metadata) String() string {
	return fmt.Sprintf(
		metaFormat,
		m.Architecture,
		m.BuildDate,
		m.Compiler,
		m.GitCommit,
		m.GoVersion,
		m.OperatingSystem,
	)
}
