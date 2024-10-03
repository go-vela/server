// SPDX-License-Identifier: Apache-2.0

package pipeline

// Metadata is the pipeline representation of the metadata block for a pipeline.
//
// swagger:model PipelineMetadata
type Metadata struct {
	Template    bool           `json:"template,omitempty" yaml:"template,omitempty"`
	Clone       bool           `json:"clone,omitempty" yaml:"clone,omitempty"`
	Environment []string       `json:"environment,omitempty" yaml:"environment,omitempty"`
	AutoCancel  *CancelOptions `json:"auto_cancel,omitempty" yaml:"auto_cancel,omitempty"`
}

// CancelOptions is the pipeline representation of the auto_cancel block for a pipeline.
type CancelOptions struct {
	Running       bool `yaml:"running,omitempty" json:"running,omitempty"`
	Pending       bool `yaml:"pending,omitempty" json:"pending,omitempty"`
	DefaultBranch bool `yaml:"default_branch,omitempty" json:"default_branch,omitempty"`
}
