// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

type (
	// Deployment is the yaml representation of a
	// deployment block in a pipeline.
	Deployment struct {
		Targets    raw.StringSlice `yaml:"targets,omitempty" json:"targets,omitempty" jsonschema:"description=List of deployment targets for the deployment block.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-targets-key"`
		Parameters ParameterMap    `yaml:"parameters,omitempty" json:"parameters,omitempty" jsonschema:"description=List of parameters for the deployment block.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Template   StepTemplate    `yaml:"template,omitempty" json:"template,omitempty" jsonschema:"description=Name of template to expand in the deployment block.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-template-key"`
	}

	// ParameterMap is the yaml representation
	// of the parameters block for a deployment block of a pipeline.
	ParameterMap map[string]*Parameter

	// Parameters is the yaml representation of the deploy parameters
	// from a deployment block in a pipeline.
	Parameter struct {
		Description string          `yaml:"description,omitempty" json:"description,omitempty" jsonschema:"description=Description of the parameter.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Type        string          `yaml:"type,omitempty" json:"type,omitempty" jsonschema:"description=Type of the parameter.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Required    bool            `yaml:"required,omitempty" json:"required,omitempty" jsonschema:"description=Flag indicating if the parameter is required.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Options     raw.StringSlice `yaml:"options,omitempty" json:"options,omitempty" jsonschema:"description=List of options for the parameter.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Min         int             `yaml:"min,omitempty" json:"min,omitempty" jsonschema:"description=Minimum value for the parameter.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
		Max         int             `yaml:"max,omitempty" json:"max,omitempty" jsonschema:"description=Maximum value for the parameter.\nReference: https://go-vela.github.io/docs/reference/yaml/deployments/#the-parameters-key"`
	}
)

// ToPipeline converts the Deployment type
// to a pipeline Deployment type.
func (d *Deployment) ToPipeline() *pipeline.Deployment {
	return &pipeline.Deployment{
		Targets:    d.Targets,
		Parameters: d.Parameters.ToPipeline(),
	}
}

// ToPipeline converts the Parameters type
// to a pipeline Parameters type.
func (p *ParameterMap) ToPipeline() pipeline.ParameterMap {
	if len(*p) == 0 {
		return nil
	}

	// parameter map we want to return
	parameterMap := make(pipeline.ParameterMap)

	// iterate through each element in the parameter map
	for k, v := range *p {
		// add the element to the pipeline parameter map
		parameterMap[k] = &pipeline.Parameter{
			Description: v.Description,
			Type:        v.Type,
			Required:    v.Required,
			Options:     v.Options,
			Min:         v.Min,
			Max:         v.Max,
		}
	}

	return parameterMap
}
