// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
)

// Validate verifies the yaml configuration is valid.
func (c *client) ValidateYAML(p *yaml.Build) error {
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

	return result
}

func (c *client) ValidatePipeline(p *pipeline.Build) error {
	var result error

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
	err = validateSteps(p.Steps, make(map[string]bool), "")
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// validateServices is a helper function that verifies the
// services block in the yaml configuration is valid.
func validateServices(s pipeline.ContainerSlice) error {
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
func validateStages(s pipeline.StageSlice) error {
	nameMap := make(map[string]bool)

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

		err := validateSteps(stage.Steps, nameMap, stage.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateSteps is a helper function that verifies the
// steps block in the yaml configuration is valid.
func validateSteps(s pipeline.ContainerSlice, nameMap map[string]bool, stageName string) error {
	reportCount := 0

	reportMap := make(map[string]string)

	for _, step := range s {
		if len(step.Name) == 0 {
			return fmt.Errorf("no name provided for step")
		}

		if len(step.Image) == 0 {
			return fmt.Errorf("no image provided for step %s", step.Name)
		}

		if step.Name == "clone" || step.Name == "init" {
			continue
		}

		if _, ok := nameMap[stageName+"_"+step.Name]; ok {
			return fmt.Errorf("step `%s` is already defined", step.Name)
		}

		nameMap[stageName+"_"+step.Name] = true

		if s, ok := reportMap[step.ReportAs]; ok {
			return fmt.Errorf("report_as to %s for step %s is already targeted by step %s", step.ReportAs, step.Name, s)
		}

		if len(step.ReportAs) > 0 {
			reportMap[step.ReportAs] = step.Name
			reportCount++
		}

		if len(step.Commands) == 0 && len(step.Environment) == 0 &&
			len(step.Secrets) == 0 && !step.Detach {
			return fmt.Errorf("no commands, environment, or secrets provided for step %s", step.Name)
		}
	}

	if reportCount > constants.ReportStepStatusLimit {
		return fmt.Errorf("report_as is limited to %d steps, counted %d", constants.ReportStepStatusLimit, reportCount)
	}

	return nil
}
