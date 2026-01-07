// SPDX-License-Identifier: Apache-2.0

package yaml

import "github.com/go-vela/server/compiler/types/pipeline"

// Artifacts represents the structure for artifacts configuration.
type Artifacts struct {
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// ToPipeline converts the Artifacts type
// to a pipeline Artifacts type.
func (t *Artifacts) ToPipeline() *pipeline.Artifacts {
	return &pipeline.Artifacts{
		Paths: t.Paths,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the Artifacts type.
func (t *Artifacts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Artifacts we try unmarshalling to
	artifacts := new(struct {
		Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	})

	// attempt to unmarshal Artifacts type
	err := unmarshal(artifacts)
	if err != nil {
		return err
	}

	// set the results field
	t.Paths = artifacts.Paths

	return nil
}
