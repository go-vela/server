// SPDX-License-Identifier: Apache-2.0

package constants

// Repo get pipeline types.
const (
	// PipelineTypeYAML defines the pipeline type for allowing users
	// in Vela to control their pipeline being compiled as yaml.
	PipelineTypeYAML = "yaml"

	// PipelineTypeGo defines the pipeline type for allowing users
	// in Vela to control their pipeline being compiled as Go templates.
	PipelineTypeGo = "go"

	// PipelineTypeStarlark defines the pipeline type for allowing users
	// in Vela to control their pipeline being compiled as Starlark templates.
	PipelineTypeStarlark = "starlark"
)

// Repo ApproveBuild types.
const (
	// ApproveForkAlways defines the CI strategy of having a repo administrator approve
	// all builds triggered from a forked PR.
	ApproveForkAlways = "fork-always"

	// ApproveForkNoWrite defines the CI strategy of having a repo administrator approve
	// all builds triggered from a forked PR where the author does not have write access.
	ApproveForkNoWrite = "fork-no-write"

	// ApproveOnce defines the CI strategy of having a repo administrator approve
	// all builds triggered from an outside contributor if this is their first time contributing.
	ApproveOnce = "first-time"

	// ApproveNever defines the CI strategy of never having to approve CI builds from outside contributors.
	ApproveNever = "never"
)
