// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/go-vela/types/yaml"
)

// Validate verifies the yaml configuration is valid.
func (c *client) Validate(p *yaml.Build) error {
	var result error
	// check a version is provided
	if len(p.Version) == 0 {
		result = multierror.Append(result, fmt.Errorf("no \"version:\" YAML property provided"))
	}

	// check that stages or steps are provided
	if len(p.Stages) == 0 && len(p.Steps) == 0 && (!p.Metadata.RenderInline && len(p.Templates) == 0) {
		result = multierror.Append(result, fmt.Errorf("no stages, steps or templates provided"))
	}

	// check that stages and steps aren't provided
	if len(p.Stages) > 0 && len(p.Steps) > 0 {
		result = multierror.Append(result, fmt.Errorf("stages and steps provided"))
	}

	if p.Metadata.RenderInline {
		for _, step := range p.Steps {
			if step.Template.Name != "" {
				result = multierror.Append(result, fmt.Errorf("step %s: cannot combine render_inline and a step that references a template", step.Name))
			}
		}

		for _, stage := range p.Stages {
			for _, step := range stage.Steps {
				if step.Template.Name != "" {
					result = multierror.Append(result, fmt.Errorf("step %s.%s: cannot combine render_inline and a step that references a template", stage.Name, step.Name))
				}
			}
		}
	}

	// validate the services block provided
	err := validateServices(p.Services)
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the stages block provided
	err = validateStages(p.Stages)
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the steps block provided
	err = validateSteps(p.Steps)
	if err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		c.log.AppendData([]byte(fmt.Sprintf("pipeline Validate error: %v\n", result)))
	}

	return result
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

			if len(step.Image) == 0 && len(step.Template.Name) == 0 {
				return fmt.Errorf("no image or template provided for step %s for stage %s", step.Name, stage.Name)
			}

			if step.Name == "clone" || step.Name == "init" {
				continue
			}

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

		if len(step.Commands) == 0 && len(step.Environment) == 0 &&
			len(step.Parameters) == 0 && len(step.Secrets) == 0 &&
			len(step.Template.Name) == 0 && !step.Detach {
			return fmt.Errorf("no commands, environment, parameters, secrets or template provided for step %s", step.Name)
		}
	}

	return nil
}
