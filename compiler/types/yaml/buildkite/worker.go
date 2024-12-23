// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

// Worker is the yaml representation of a worker
// from a worker block in a pipeline.
type Worker struct {
	Flavor   string `yaml:"flavor,omitempty"   json:"flavor,omitempty"   jsonschema:"minLength=1,description=Flavor identifier for worker.\nReference: https://go-vela.github.io/docs/reference/yaml/worker/#the-flavor-key,example=large"`
	Platform string `yaml:"platform,omitempty" json:"platform,omitempty" jsonschema:"minLength=1,description=Platform identifier for the worker.\nReference: https://go-vela.github.io/docs/reference/yaml/worker/#the-platform-key,example=kubernetes"`
}

// ToPipeline converts the Worker type
// to a pipeline Worker type.
func (w *Worker) ToPipeline() *pipeline.Worker {
	return &pipeline.Worker{
		Flavor:   w.Flavor,
		Platform: w.Platform,
	}
}

func (w *Worker) ToYAML() *yaml.Worker {
	return &yaml.Worker{
		Flavor:   w.Flavor,
		Platform: w.Platform,
	}
}
