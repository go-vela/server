// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-multierror"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
)

// ValidateYAML verifies the yaml configuration is valid.
func (c *Client) ValidateYAML(p *yaml.Build) error {
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
	err := validateYAMLServices(p.Services)
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the stages block provided
	err = validateYAMLStages(p.Stages)
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the steps block provided
	err = validateYAMLSteps(p.Steps)
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// validateYAMLStages is a helper function that verifies the
// stages block in the yaml configuration is valid.
func validateYAMLStages(s yaml.StageSlice) error {
	for _, stage := range s {
		if len(stage.Name) == 0 {
			return fmt.Errorf("no name provided for stage")
		}

		// validate that a stage is not referencing itself in needs
		if slices.Contains(stage.Needs, stage.Name) {
			return fmt.Errorf("stage %s references itself in 'needs' declaration", stage.Name)
		}

		err := validateYAMLSteps(stage.Steps)
		if err != nil {
			return err
		}
	}

	return nil
}

// validateYAMLSteps is a helper function that verifies the
// steps block in the yaml configuration is valid.
func validateYAMLSteps(s yaml.StepSlice) error {
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

		if len(step.Commands) == 0 && len(step.Environment) == 0 &&
			len(step.Parameters) == 0 && len(step.Secrets) == 0 &&
			len(step.Template.Name) == 0 && !step.Detach {
			return fmt.Errorf("no commands, environment, parameters, secrets or template provided for step %s", step.Name)
		}
	}

	return nil
}

// validateYAMLServices is a helper function that verifies the
// services block in the yaml configuration is valid.
func validateYAMLServices(s yaml.ServiceSlice) error {
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

// ValidatePipeline verifies the final pipeline configuration is valid.
func (c *Client) ValidatePipeline(p *pipeline.Build) error {
	var result error

	// report count for custom report containers
	reportCount := 0

	// validate the services block provided
	err := validatePipelineContainers(p.Services, &reportCount, make(map[string]string), make(map[string]bool), "")
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the stages block provided
	err = validatePipelineStages(p.Stages)
	if err != nil {
		result = multierror.Append(result, err)
	}

	// validate the steps block provided
	err = validatePipelineContainers(p.Steps, &reportCount, make(map[string]string), make(map[string]bool), "")
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// validatePipelineStages is a helper function that verifies the
// stages block in the final pipeline configuration is valid.
func validatePipelineStages(s pipeline.StageSlice) error {
	reportMap := make(map[string]string)
	reportCount := 0

	nameMap := make(map[string]bool)

	for _, stage := range s {
		err := validatePipelineContainers(stage.Steps, &reportCount, reportMap, nameMap, stage.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// validatePipelineContainers is a helper function that
// ensures custom report containers do not exceed the limit
// and that the container names are unique.
func validatePipelineContainers(s pipeline.ContainerSlice, reportCount *int, reportMap map[string]string, nameMap map[string]bool, stageName string) error {
	for _, ctn := range s {
		if ctn.Name == "clone" || ctn.Name == "init" {
			continue
		}

		if _, ok := nameMap[stageName+"_"+ctn.Name]; ok {
			return fmt.Errorf("step `%s` is already defined", ctn.Name)
		}

		nameMap[stageName+"_"+ctn.Name] = true

		if s, ok := reportMap[ctn.ReportAs]; ok {
			return fmt.Errorf("report_as to %s for step %s is already targeted by step %s", ctn.ReportAs, ctn.Name, s)
		}

		if len(ctn.ReportAs) > 0 {
			reportMap[ctn.ReportAs] = ctn.Name
			*reportCount++
		}
	}

	if *reportCount > constants.ReportStepStatusLimit {
		return fmt.Errorf("report_as is limited to %d steps, counted %d", constants.ReportStepStatusLimit, reportCount)
	}

	return nil
}
