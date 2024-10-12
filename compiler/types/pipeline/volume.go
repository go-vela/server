// SPDX-License-Identifier: Apache-2.0

package pipeline

type (
	// VolumeSlice is the pipeline representation of
	// the volumes block for a step in a pipeline.
	//
	// swagger:model PipelineVolumeSlice
	VolumeSlice []*Volume

	// Volume is the pipeline representation of a volume
	// from a volumes block for a step in a pipeline.
	//
	// swagger:model PipelineVolume
	Volume struct {
		Source      string `json:"source,omitempty"      yaml:"source,omitempty"`
		Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
		AccessMode  string `json:"access_mode,omitempty" yaml:"access_mode,omitempty"`
	}
)
