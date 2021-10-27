// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/yaml"
)

// Validate verifies the the yaml configuration is valid.
func (c *client) Validate(p *yaml.Build) error {
	// check a version is provided
	if len(p.Version) == 0 {
		return fmt.Errorf("no version provided")
	}

	// check that stages or steps are provided
	if len(p.Stages) == 0 && len(p.Steps) == 0 {
		return fmt.Errorf("no stages or steps provided")
	}

	// check that stages and steps aren't provided
	if len(p.Stages) > 0 && len(p.Steps) > 0 {
		return fmt.Errorf("stages and steps provided")
	}

	// validate the services block provided
	err := validateServices(p.Services)
	if err != nil {
		return err
	}

	// validate the stages block provided
	err = validateStages(p.Stages)
	if err != nil {
		return err
	}

	// validate the steps block provided
	err = validateSteps(p.Steps)
	if err != nil {
		return err
	}

	return nil
}

// validateServices is a helper function that verifies the
// services block in the yaml configuration is valid.
func validateServices(s yaml.ServiceSlice) error {
	for _, service := range s {
		if len(service.Name) == 0 {
			return fmt.Errorf("no name provided for service")
		}

		if len(service.Image) == 0 {
			return fmt.Errorf("no image provided for service %s", service.Name)
		}
	}

	return nil
}

// validateStages is a helper function that verifies the
// stages block in the yaml configuration is valid.
func validateStages(s yaml.StageSlice) error {
	for _, stage := range s {
		if len(stage.Name) == 0 {
			return fmt.Errorf("no name provided for stage")
		}

		// validate that a stage is not referencing itself in needs
		for _, need := range stage.Needs {
			if stage.Name == need {
				return fmt.Errorf("stage %s references itself in 'needs' declaration", stage.Name)
			}
		}

		for _, step := range stage.Steps {
			if len(step.Name) == 0 {
				return fmt.Errorf("no name provided for step for stage %s", stage.Name)
			}

			// nolint: lll // ignore simplification here
			if len(step.Image) == 0 && len(step.Template.Name) == 0 {
				return fmt.Errorf("no image or template provided for step %s for stage %s", step.Name, stage.Name)
			}

			if step.Name == "clone" || step.Name == "init" {
				continue
			}

			// nolint: lll // ignore simplification here
			if len(step.Commands) == 0 && len(step.Environment) == 0 &&
				len(step.Parameters) == 0 && len(step.Secrets) == 0 &&
				len(step.Template.Name) == 0 && !step.Detach {
				return fmt.Errorf("no commands, environment, parameters, secrets or template provided for step %s for stage %s", step.Name, stage.Name)
			}
		}
	}

	return nil
}

// validateSteps is a helper function that verifies the
// steps block in the yaml configuration is valid.
func validateSteps(s yaml.StepSlice) error {
	for _, step := range s {
		if len(step.Name) == 0 {
			return fmt.Errorf("no name provided for step")
		}

		if len(step.Image) == 0 && len(step.Template.Name) == 0 {
			return fmt.Errorf("no image or template provided for step %s", step.Name)
		}

		if step.Name == "clone" || step.Name == "init" {
			continue
		}

		// nolint: lll // ignore simplification here
		if len(step.Commands) == 0 && len(step.Environment) == 0 &&
			len(step.Parameters) == 0 && len(step.Secrets) == 0 &&
			len(step.Template.Name) == 0 && !step.Detach {
			return fmt.Errorf("no commands, environment, parameters, secrets or template provided for step %s", step.Name)
		}
	}

	return nil
}
