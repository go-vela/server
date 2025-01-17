// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/constants"
)

type (
	// StepSlice is the yaml representation
	// of the steps block for a pipeline.
	StepSlice []*Step

	// Step is the yaml representation of a step
	// from the steps block for a pipeline.
	Step struct {
		Ruleset     Ruleset                `yaml:"ruleset,omitempty"     json:"ruleset,omitempty"     jsonschema:"description=Conditions to limit the execution of the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Commands    raw.StringSlice        `yaml:"commands,omitempty"    json:"commands,omitempty"    jsonschema:"description=Execution instructions to run inside the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-commands-key"`
		Entrypoint  raw.StringSlice        `yaml:"entrypoint,omitempty"  json:"entrypoint,omitempty"  jsonschema:"description=Command to execute inside the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-entrypoint-key"`
		Secrets     StepSecretSlice        `yaml:"secrets,omitempty"     json:"secrets,omitempty"     jsonschema:"description=Sensitive variables injected into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-secrets-key"`
		Template    StepTemplate           `yaml:"template,omitempty"    json:"template,omitempty"    jsonschema:"oneof_required=template,description=Name of template to expand in the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-template-key"`
		TestReport  TestReport             `yaml:"test_report,omitempty" json:"test_report,omitempty" jsonschema:"description=Test report configuration for the step.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-test_report-key"`
		Ulimits     UlimitSlice            `yaml:"ulimits,omitempty"     json:"ulimits,omitempty"     jsonschema:"description=Set the user limits for the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ulimits-key"`
		Volumes     VolumeSlice            `yaml:"volumes,omitempty"     json:"volumes,omitempty"     jsonschema:"description=Mount volumes for the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-volume-key"`
		Image       string                 `yaml:"image,omitempty"       json:"image,omitempty"       jsonschema:"oneof_required=image,minLength=1,description=Docker image to use to create the ephemeral container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-image-key"`
		Name        string                 `yaml:"name,omitempty"        json:"name,omitempty"        jsonschema:"required,minLength=1,description=Unique name for the step.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-name-key"`
		Pull        string                 `yaml:"pull,omitempty"        json:"pull,omitempty"        jsonschema:"enum=always,enum=not_present,enum=on_start,enum=never,default=not_present,description=Declaration to configure if and when the Docker image is pulled.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-pull-key"`
		Environment raw.StringSliceMap     `yaml:"environment,omitempty" json:"environment,omitempty" jsonschema:"description=Provide environment variables injected into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-environment-key"`
		Parameters  map[string]interface{} `yaml:"parameters,omitempty"  json:"parameters,omitempty"  jsonschema:"description=Extra configuration variables for a plugin.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-parameters-key"`
		Detach      bool                   `yaml:"detach,omitempty"      json:"detach,omitempty"      jsonschema:"description=Run the container in a detached (headless) state.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-detach-key"`
		Privileged  bool                   `yaml:"privileged,omitempty"  json:"privileged,omitempty"  jsonschema:"description=Run the container with extra privileges.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-privileged-key"`
		User        string                 `yaml:"user,omitempty"        json:"user,omitempty"        jsonschema:"description=Set the user for the container.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-user-key"`
		ReportAs    string                 `yaml:"report_as,omitempty"   json:"report_as,omitempty"   jsonschema:"description=Set the name of the step to report as.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-report_as-key"`
		IDRequest   string                 `yaml:"id_request,omitempty"  json:"id_request,omitempty"  jsonschema:"description=Request ID Request Token for the step.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-id_request-key"`
	}
)

// ToPipeline converts the StepSlice type
// to a pipeline ContainerSlice type.
func (s *StepSlice) ToPipeline() *pipeline.ContainerSlice {
	// step slice we want to return
	stepSlice := new(pipeline.ContainerSlice)

	// iterate through each element in the step slice
	for _, step := range *s {
		// append the element to the pipeline container slice
		*stepSlice = append(*stepSlice, &pipeline.Container{
			Commands:    step.Commands,
			Detach:      step.Detach,
			Entrypoint:  step.Entrypoint,
			Environment: step.Environment,
			Image:       step.Image,
			Name:        step.Name,
			Privileged:  step.Privileged,
			Pull:        step.Pull,
			Ruleset:     *step.Ruleset.ToPipeline(),
			Secrets:     *step.Secrets.ToPipeline(),
			Ulimits:     *step.Ulimits.ToPipeline(),
			Volumes:     *step.Volumes.ToPipeline(),
			User:        step.User,
			ReportAs:    step.ReportAs,
			IDRequest:   step.IDRequest,
			TestReport:  *step.TestReport.ToPipeline(),
		})
	}

	return stepSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the StepSlice type.
func (s *StepSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// step slice we try unmarshalling to
	stepSlice := new([]*Step)

	// attempt to unmarshal as a step slice type
	err := unmarshal(stepSlice)
	if err != nil {
		return err
	}

	// iterate through each element in the step slice
	for _, step := range *stepSlice {
		// handle nil step to avoid panic
		if step == nil {
			return fmt.Errorf("invalid step with nil content found")
		}

		// implicitly set `pull` field if empty
		if len(step.Pull) == 0 {
			step.Pull = constants.PullNotPresent
		}

		// TODO: remove this in a future release
		//
		// handle true deprecated pull policy
		//
		// a `true` pull policy equates to `always`
		if strings.EqualFold(step.Pull, "true") {
			step.Pull = constants.PullAlways
		}

		// TODO: remove this in a future release
		//
		// handle false deprecated pull policy
		//
		// a `false` pull policy equates to `not_present`
		if strings.EqualFold(step.Pull, "false") {
			step.Pull = constants.PullNotPresent
		}
	}

	// overwrite existing StepSlice
	*s = StepSlice(*stepSlice)

	return nil
}

// MergeEnv takes a list of environment variables and attempts
// to set them in the step environment. If the environment
// variable already exists in the step, than this will
// overwrite the existing environment variable.
func (s *Step) MergeEnv(environment map[string]string) error {
	// check if the step container is empty
	if s == nil || s.Environment == nil {
		// TODO: evaluate if we should error here
		//
		// immediately return and do nothing
		//
		// treated as a no-op
		return nil
	}

	// check if the environment provided is empty
	if environment == nil {
		return fmt.Errorf("empty environment provided for step %s", s.Name)
	}

	// iterate through all environment variables provided
	for key, value := range environment {
		// set or update the step environment variable
		s.Environment[key] = value
	}

	return nil
}
