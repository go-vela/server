// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"slices"
	"strconv"
)

type (
	// Deployment is the pipeline representation of the deployment block for a pipeline.
	//
	// swagger:model PipelineDeployment
	Deployment struct {
		Targets    []string     `json:"targets,omitempty"    yaml:"targets,omitempty"`
		Parameters ParameterMap `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	}

	ParameterMap map[string]*Parameter

	// Parameters is the pipeline representation of the deploy parameters
	// from a deployment block in a pipeline.
	//
	// swagger:model PipelineParameters
	Parameter struct {
		Description string   `json:"description,omitempty" yaml:"description,omitempty"`
		Type        string   `json:"type,omitempty"        yaml:"type,omitempty"`
		Required    bool     `json:"required,omitempty"    yaml:"required,omitempty"`
		Options     []string `json:"options,omitempty"     yaml:"options,omitempty"`
		Min         int      `json:"min,omitempty"         yaml:"min,omitempty"`
		Max         int      `json:"max,omitempty"         yaml:"max,omitempty"`
	}
)

// Empty returns true if the provided deployment is empty.
func (d *Deployment) Empty() bool {
	// return true if deployment is nil
	if d == nil {
		return true
	}

	// return true if every deployment field is empty
	if len(d.Targets) == 0 &&
		len(d.Parameters) == 0 {
		return true
	}

	// return false if any of the deployment fields are provided
	return false
}

// Validate checks the build ruledata and parameters against the deployment configuration.
func (d *Deployment) Validate(target string, inputParams map[string]string) error {
	if d.Empty() {
		return nil
	}

	// validate targets
	if len(d.Targets) > 0 && !slices.Contains(d.Targets, target) {
		return fmt.Errorf("deployment target `%s` not found in deployment config targets", target)
	}

	// validate params
	for kConfig, vConfig := range d.Parameters {
		var (
			inputStr string
			ok       bool
		)
		// check if the parameter is required
		if vConfig.Required {
			// check if the parameter is provided
			if inputStr, ok = inputParams[kConfig]; !ok {
				return fmt.Errorf("deployment parameter %s is required", kConfig)
			}
		} else {
			// check if the parameter is provided
			if inputStr, ok = inputParams[kConfig]; !ok {
				continue
			}
		}

		// check if the parameter is an option
		if len(vConfig.Options) > 0 && len(inputStr) > 0 {
			if !slices.Contains(vConfig.Options, inputStr) {
				return fmt.Errorf("deployment parameter %s is not a valid option", kConfig)
			}
		}

		// check if the parameter is the correct type
		if len(vConfig.Type) > 0 && len(inputStr) > 0 {
			switch vConfig.Type {
			case "integer", "int", "number":
				val, err := strconv.Atoi(inputStr)
				if err != nil {
					return fmt.Errorf("deployment parameter %s is not an integer", kConfig)
				}

				if vConfig.Max != 0 && val < vConfig.Min {
					return fmt.Errorf("deployment parameter %s is less than the minimum value", kConfig)
				}

				if vConfig.Max != 0 && val > vConfig.Max {
					return fmt.Errorf("deployment parameter %s is greater than the maximum value", kConfig)
				}
			case "boolean", "bool":
				if _, err := strconv.ParseBool(inputStr); err != nil {
					return fmt.Errorf("deployment parameter %s is not a boolean", kConfig)
				}
			default:
				continue
			}
		}
	}

	return nil
}
