// SPDX-License-Identifier: Apache-2.0

package yaml

import "github.com/go-vela/server/compiler/types/pipeline"

// Artifact represents the structure for artifacts configuration.
type Artifacts struct {
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// ToPipeline converts the Artifact type
// to a pipeline Artifact type.
func (a *Artifacts) ToPipeline() *pipeline.Artifacts {
	return &pipeline.Artifacts{
		Paths: a.Paths,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the Artifact type.
func (a *Artifacts) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// artifacts we try unmarshalling to
	artifacts := new(struct {
		Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	})

	// attempt to unmarshal artifacts type
	err := unmarshal(artifacts)
	if err != nil {
		return err
	}

	// set the paths field
	a.Paths = artifacts.Paths

	return nil
}
