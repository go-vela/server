// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

// Artifacts Artifact represents the structure for artifacts configuration.
type Artifacts struct {
	Paths raw.StringSlice `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// ToPipeline converts the Artifact type
// to a pipeline Artifact type.
func (a *Artifacts) ToPipeline() *pipeline.Artifacts {
	return &pipeline.Artifacts{
		Paths: a.Paths,
	}
}
