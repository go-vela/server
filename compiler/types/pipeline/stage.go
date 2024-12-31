// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"

	"github.com/go-vela/server/constants"
)

type (
	// StageSlice is the pipeline representation
	// of the stages block for a pipeline.
	//
	// swagger:model PipelineStageSlice
	StageSlice []*Stage

	// Stage is the pipeline representation
	// of a stage in a pipeline.
	//
	// swagger:model PipelineStage
	Stage struct {
		Done        chan error        `json:"-"                     yaml:"-"`
		Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
		Name        string            `json:"name,omitempty"        yaml:"name,omitempty"`
		Needs       []string          `json:"needs,omitempty"       yaml:"needs,omitempty"`
		Independent bool              `json:"independent,omitempty" yaml:"independent,omitempty"`
		Steps       ContainerSlice    `json:"steps,omitempty"       yaml:"steps,omitempty"`
	}
)

// Purge removes the steps, from the stages, that have
// a ruleset that do not match the provided ruledata.
// If all steps from a stage are removed, then the
// entire stage is removed from the pipeline.
func (s *StageSlice) Purge(r *RuleData) (*StageSlice, error) {
	counter := 1
	stages := new(StageSlice)

	// iterate through each stage for the pipeline
	for _, stage := range *s {
		containers := new(ContainerSlice)

		// iterate through each step for the stage in the pipeline
		for _, step := range stage.Steps {
			match, err := step.Ruleset.Match(r, step.Environment)
			if err != nil {
				return nil, fmt.Errorf("unable to process ruleset for step %s: %w", step.Name, err)
			}

			// verify ruleset matches
			if match {
				// overwrite the step number with the step counter
				step.Number = counter

				// increment step counter
				counter = counter + 1

				// append the step to the new slice of containers
				*containers = append(*containers, step)
			}
		}

		// no steps for the stage so we continue processing to the next stage
		if len(*containers) == 0 {
			continue
		}

		// overwrite the steps for the stage with the new slice of steps
		stage.Steps = *containers

		// append the stage to the new slice of stages
		*stages = append(*stages, stage)
	}

	// return the new slice of stages
	return stages, nil
}

// Sanitize cleans the fields for every step in each stage so they
// can be safely executed on the worker. The fields are sanitized
// based off of the provided runtime driver which is setup on every
// worker. Currently, this function supports the following runtimes:
//
//   - Docker
//   - Kubernetes
func (s *StageSlice) Sanitize(driver string) *StageSlice {
	stages := new(StageSlice)

	switch driver {
	// sanitize container for Docker
	case constants.DriverDocker:
		for _, stage := range *s {
			stage.Steps.Sanitize(driver)

			*stages = append(*stages, stage)
		}

		return stages
	// sanitize container for Kubernetes
	case constants.DriverKubernetes:
		for _, stage := range *s {
			stage.Steps.Sanitize(driver)

			*stages = append(*stages, stage)
		}

		return stages
	// unrecognized driver
	default:
		// log here?
		return nil
	}
}

// Empty returns true if the provided stage is empty.
func (s *Stage) Empty() bool {
	// return true if the stage is nil
	if s == nil {
		return true
	}

	// return true if every stage field is empty
	if len(s.Name) == 0 &&
		len(s.Needs) == 0 &&
		len(s.Steps) == 0 &&
		len(s.Environment) == 0 {
		return true
	}

	// return false if any of the stage fields are not empty
	return false
}

// MergeEnv takes a list of environment variables and attempts
// to set them in the stage environment. If the environment
// variable already exists in the stage, then this will
// overwrite the existing environment variable.
func (s *Stage) MergeEnv(environment map[string]string) error {
	// check if the stage is empty
	if s.Empty() {
		// TODO: evaluate if we should error here
		//
		// immediately return and do nothing
		//
		// treated as a no-op
		return nil
	}

	// check if the environment provided is empty
	if environment == nil {
		return fmt.Errorf("empty environment provided for stage %s", s.Name)
	}

	// iterate through all environment variables provided
	for key, value := range environment {
		// set or update the stage environment variable
		s.Environment[key] = value
	}

	return nil
}
