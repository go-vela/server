// SPDX-License-Identifier: Apache-2.0

package constants

// Pipeline types.
const (
	// PipelineStages defines the type for a pipeline with stages.
	PipelineStage = "stages"

	// PipelineStep defines the type for a pipeline with steps.
	PipelineStep = "steps"

	// PipelineTemplate defines the type for a pipeline as a template.
	PipelineTemplate = "template"

	// InitName defines the name for the initialization step.
	InitName = "init"

	// CloneName defines the name for the clone step.
	CloneName = "clone"

	// DefaultHomeDir defines the default home directory for steps.
	DefaultHomeDir = "/root"

	// DefaultShell defines the default shell for steps.
	DefaultShell = "/bin/sh"

	// PipelineIDPattern defines the string pattern for the pipeline ID
	// format: `<org>_<repo>_<build number>`
	PipelineIDPattern = "%s_%s_%d"

	// StageIDPattern defines the string pattern for the stage ID
	// format: `<org name>_<repo name>_<build number>_<stage name>_<step name>`
	StageIDPattern = "%s_%s_%d_%s_%s"

	// StepIDPattern defines the string pattern for the step ID
	// format: `step_<org name>_<repo name>_<build number>_<step name>`
	StepIDPattern = "step_%s_%s_%d_%s"

	// ServiceIDPattern defines the string pattern for the service ID
	// format: `service_<org name>_<repo name>_<build number>_<service name>`
	ServiceIDPattern = "service_%s_%s_%d_%s"

	// SecretIDPattern defines the string pattern for the secret ID
	// format: `secret_<org name>_<repo name>_<build number>_<secret name>`
	//
	//nolint:gosec // ignore gosec keying off of secret as no credentials are hardcoded
	SecretIDPattern = "secret_%s_%s_%d_%s"
)
