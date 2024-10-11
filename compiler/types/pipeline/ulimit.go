// SPDX-License-Identifier: Apache-2.0

package pipeline

type (
	// UlimitSlice is the pipeline representation of
	// the ulimits block for a step in a pipeline.
	//
	// swagger:model PipelineUlimitSlice
	UlimitSlice []*Ulimit

	// Ulimit is the pipeline representation of a ulimit
	// from the ulimits block for a step in a pipeline.
	//
	// swagger:model PipelineUlimit
	Ulimit struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`
		Soft int64  `json:"soft,omitempty" yaml:"soft,omitempty"`
		Hard int64  `json:"hard,omitempty" yaml:"hard,omitempty"`
	}
)
