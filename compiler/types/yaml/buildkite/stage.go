// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"fmt"
	bkYaml "github.com/buildkite/yaml"
	"github.com/invopop/jsonschema"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
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
func (s *StageSlice) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// map slice we try unmarshalling to
	mapSlice := new(bkYaml.MapSlice)

	// attempt to unmarshal as a map slice type
	err := unmarshal(mapSlice)
	if err != nil {
		return err
	}

	// iterate through each element in the map slice
	for _, v := range *mapSlice {
		// stage we try unmarshalling to
		stage := new(Stage)

		// marshal interface value from ordered map
		out, _ := bkYaml.Marshal(v.Value)

		// unmarshal interface value as stage
		err = bkYaml.Unmarshal(out, stage)
		if err != nil {
			return err
		}

		// implicitly set stage `name` if empty
		if len(stage.Name) == 0 {
			stage.Name = fmt.Sprintf("%v", v.Key)
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
	// map slice to return as marshaled output
	var output bkYaml.MapSlice

	// loop over the input stages
	for _, inputStage := range s {
		// create a new stage
		outputStage := new(Stage)

		// add the existing needs to the new stage
		outputStage.Needs = inputStage.Needs

		// add the existing dependent tag to the new stage
		outputStage.Independent = inputStage.Independent

		// add the existing steps to the new stage
		outputStage.Steps = inputStage.Steps

		// append stage to MapSlice
		output = append(output, bkYaml.MapItem{Key: inputStage.Name, Value: outputStage})
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

func (s *Stage) ToYAML() *yaml.Stage {
	if s == nil {
		return nil
	}

	return &yaml.Stage{
		Environment: s.Environment,
		Name:        s.Name,
		Needs:       s.Needs,
		Independent: s.Independent,
		Steps:       *s.Steps.ToYAML(),
	}
}

func (s *StageSlice) ToYAML() *yaml.StageSlice {
	// stage slice we want to return
	stageSlice := new(yaml.StageSlice)

	// iterate through each element in the stage slice
	for _, stage := range *s {
		// append the element to the yaml stage slice
		*stageSlice = append(*stageSlice, stage.ToYAML())
	}

	return stageSlice
}
