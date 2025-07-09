// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"fmt"

	"github.com/invopop/jsonschema"
	"go.yaml.in/yaml/v3"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

type (
	// StageSlice is the yaml representation
	// of the stages block for a pipeline.
	StageSlice []*Stage

	// Stage is the yaml representation
	// of a stage in a pipeline.
	Stage struct {
		Environment raw.StringSliceMap `yaml:"environment,omitempty" json:"environment,omitempty" jsonschema:"description=Provide environment variables injected into the container environment.\nReference: https://go-vela.github.io/docs/reference/yaml/stages/#the-environment-key"`
		Name        string             `yaml:"name,omitempty"        json:"name,omitempty"        jsonschema:"minLength=1,description=Unique identifier for the stage in the pipeline.\nReference: https://go-vela.github.io/docs/reference/yaml/stages/#the-name-key"`
		Needs       raw.StringSlice    `yaml:"needs,omitempty,flow"  json:"needs,omitempty"       jsonschema:"description=Stages that must complete before starting the current one.\nReference: https://go-vela.github.io/docs/reference/yaml/stages/#the-needs-key"`
		Independent bool               `yaml:"independent,omitempty" json:"independent,omitempty" jsonschema:"description=Stage will continue executing if other stage fails"`
		Steps       StepSlice          `yaml:"steps,omitempty"       json:"steps,omitempty"       jsonschema:"required,description=Sequential execution instructions for the stage.\nReference: https://go-vela.github.io/docs/reference/yaml/stages/#the-steps-key"`
	}
)

// ToPipeline converts the StageSlice type
// to a pipeline StageSlice type.
func (s *StageSlice) ToPipeline() *pipeline.StageSlice {
	// stage slice we want to return
	stageSlice := new(pipeline.StageSlice)

	// iterate through each element in the stage slice
	for _, stage := range *s {
		// append the element to the pipeline stage slice
		*stageSlice = append(*stageSlice, &pipeline.Stage{
			Done:        make(chan error, 1),
			Environment: stage.Environment,
			Name:        stage.Name,
			Needs:       stage.Needs,
			Independent: stage.Independent,
			Steps:       *stage.Steps.ToPipeline(),
		})
	}

	return stageSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the StageSlice type.
func (s *StageSlice) UnmarshalYAML(v *yaml.Node) error {
	if v.Kind != yaml.MappingNode {
		return fmt.Errorf("invalid yaml: expected map node for stage")
	}

	// iterate through each element in the map slice
	for i := 0; i < len(v.Content); i += 2 {
		key := v.Content[i]
		value := v.Content[i+1]

		stage := new(Stage)

		// unmarshal value into stage
		err := value.Decode(stage)
		if err != nil {
			return err
		}

		// implicitly set stage `name` if empty
		if len(stage.Name) == 0 {
			stage.Name = fmt.Sprintf("%v", key.Value)
		}

		// implicitly set the stage `needs`
		if stage.Name != "clone" && stage.Name != "init" {
			// add clone if not present
			stage.Needs = func(needs []string) []string {
				for _, s := range needs {
					if s == "clone" {
						return needs
					}
				}

				return append(needs, "clone")
			}(stage.Needs)
		}
		// append stage to stage slice
		*s = append(*s, stage)
	}

	return nil
}

// MarshalYAML implements the marshaler interface for the StageSlice type.
func (s StageSlice) MarshalYAML() (interface{}, error) {
	output := new(yaml.Node)
	output.Kind = yaml.MappingNode

	for _, inputStage := range s {
		n := new(yaml.Node)

		// create new stage with existing properties
		outputStage := &Stage{
			Name:        inputStage.Name,
			Needs:       inputStage.Needs,
			Independent: inputStage.Independent,
			Steps:       inputStage.Steps,
		}

		err := n.Encode(outputStage)
		if err != nil {
			return nil, err
		}

		// append stage to map output
		output.Content = append(output.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: inputStage.Name})
		output.Content = append(output.Content, n)
	}

	return output, nil
}

// JSONSchemaExtend handles some overrides that need to be in place
// for this type for the jsonschema generation.
//
// Stages are not really a slice of stages to the user. This change
// supports the map they really are.
func (StageSlice) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.AdditionalProperties = jsonschema.FalseSchema
	schema.Items = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		".*": {
			Ref: "#/$defs/Stage",
		},
	}
	schema.Type = "object"
}

// MergeEnv takes a list of environment variables and attempts
// to set them in the stage environment. If the environment
// variable already exists in the stage, than this will
// overwrite the existing environment variable.
func (s *Stage) MergeEnv(environment map[string]string) error {
	// check if the stage is empty
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
		return fmt.Errorf("empty environment provided for stage %s", s.Name)
	}

	// iterate through all environment variables provided
	for key, value := range environment {
		// set or update the stage environment variable
		s.Environment[key] = value
	}

	return nil
}
