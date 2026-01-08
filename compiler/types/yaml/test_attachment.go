// SPDX-License-Identifier: Apache-2.0

package yaml

import "github.com/go-vela/server/compiler/types/pipeline"

// Artifact represents the structure for artifact configuration.
type Artifact struct {
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// ToPipeline converts the Artifact type
// to a pipeline Artifact type.
func (a *Artifact) ToPipeline() *pipeline.Artifact {
	return &pipeline.Artifact{
		Paths: a.Paths,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the Artifact type.
func (a *Artifact) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// artifact we try unmarshalling to
	artifact := new(struct {
		Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	})

	// attempt to unmarshal artifact type
	err := unmarshal(artifact)
	if err != nil {
		return err
	}

	// set the paths field
	a.Paths = artifact.Paths

	return nil
}
