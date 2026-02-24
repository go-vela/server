// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/invopop/jsonschema"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

type (
	// UlimitSlice is the yaml representation of
	// the ulimits block for a step in a pipeline.
	UlimitSlice []*Ulimit

	// Ulimit is the yaml representation of a ulimit
	// from the ulimits block for a step in a pipeline.
	Ulimit struct {
		Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"required,minLength=1,description=Unique name of the user limit.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ulimits-key"`
		Soft int64  `yaml:"soft,omitempty" json:"soft,omitempty" jsonschema:"description=Set the soft limit.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ulimits-key"`
		Hard int64  `yaml:"hard,omitempty" json:"hard,omitempty" jsonschema:"description=Set the hard limit.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ulimits-key"`
	}
)

// ToPipeline converts the UlimitSlice type
// to a pipeline UlimitSlice type.
func (u *UlimitSlice) ToPipeline() *pipeline.UlimitSlice {
	// ulimit slice we want to return
	ulimitSlice := new(pipeline.UlimitSlice)

	// iterate through each element in the ulimit slice
	for _, ulimit := range *u {
		// append the element to the pipeline ulimit slice
		*ulimitSlice = append(*ulimitSlice, &pipeline.Ulimit{
			Name: ulimit.Name,
			Soft: ulimit.Soft,
			Hard: ulimit.Hard,
		})
	}

	return ulimitSlice
}

// UnmarshalYAML implements the Unmarshaler interface for the UlimitSlice type.
func (u *UlimitSlice) UnmarshalYAML(unmarshal func(any) error) error {
	// string slice we try unmarshalling to
	stringSlice := new(raw.StringSlice)

	// attempt to unmarshal as a string slice type
	err := unmarshal(stringSlice)
	if err == nil {
		// iterate through each element in the string slice
		for _, ulimit := range *stringSlice {
			// split each slice element into key/value pairs
			parts := strings.Split(ulimit, "=")
			if len(parts) != 2 {
				return fmt.Errorf("ulimit %s must contain 1 `=` (equal)", ulimit)
			}

			// split each value into soft and hard limits
			limitParts := strings.Split(parts[1], ":")

			switch {
			case len(limitParts) == 1:
				// capture value for soft and hard limit
				value, err := strconv.ParseInt(limitParts[0], 10, 64)
				if err != nil {
					return err
				}

				// append the element to the ulimit slice
				*u = append(*u, &Ulimit{
					Name: parts[0],
					Soft: value,
					Hard: value,
				})

				continue
			case len(limitParts) == 2:
				// capture value for soft limit
				firstValue, err := strconv.ParseInt(limitParts[0], 10, 64)
				if err != nil {
					return err
				}

				// capture value for hard limit
				secondValue, err := strconv.ParseInt(limitParts[1], 10, 64)
				if err != nil {
					return err
				}

				// append the element to the ulimit slice
				*u = append(*u, &Ulimit{
					Name: parts[0],
					Soft: firstValue,
					Hard: secondValue,
				})

				continue
			default:
				return fmt.Errorf("ulimit %s can only contain 1 `:` (colon)", ulimit)
			}
		}

		return nil
	}

	// ulimit slice we try unmarshalling to
	ulimits := new([]*Ulimit)

	// attempt to unmarshal as a ulimit slice type
	err = unmarshal(ulimits)
	if err != nil {
		return err
	}

	// iterate through each element in the volume slice
	for _, ulimit := range *ulimits {
		// implicitly set `hard` field if empty
		if ulimit.Hard == 0 {
			ulimit.Hard = ulimit.Soft
		}
	}

	// overwrite existing UlimitSlice
	*u = UlimitSlice(*ulimits)

	return nil
}

// JSONSchemaExtend handles some overrides that need to be in place
// for this type for the jsonschema generation.
//
// Without these changes it would only allow an object per the struct,
// but we do some special handling to allow specially formatted strings.
func (Ulimit) JSONSchemaExtend(schema *jsonschema.Schema) {
	oldAddProps := schema.AdditionalProperties
	oldProps := schema.Properties
	oldReq := schema.Required

	schema.AdditionalProperties = nil
	schema.OneOf = []*jsonschema.Schema{
		{
			Type:                 "string",
			Pattern:              "[a-z]+=[0-9]+:[0-9]+",
			AdditionalProperties: oldAddProps,
		},
		{
			Type:                 "object",
			Properties:           oldProps,
			Required:             oldReq,
			AdditionalProperties: oldAddProps,
		},
	}
	schema.Properties = nil
	schema.Required = nil
	schema.Type = ""
}
