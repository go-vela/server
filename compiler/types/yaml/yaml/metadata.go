// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"slices"

	"github.com/go-vela/server/compiler/types/pipeline"
)

type (
	// Metadata is the yaml representation of
	// the metadata block for a pipeline.
	Metadata struct {
		Template     bool           `yaml:"template,omitempty"      json:"template,omitempty"      jsonschema:"description=Enables compiling the pipeline as a template.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-template-key"`
		RenderInline bool           `yaml:"render_inline,omitempty" json:"render_inline,omitempty" jsonschema:"description=Enables inline compiling for the pipeline templates.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-render-inline-key"`
		Clone        *bool          `yaml:"clone,omitempty"         json:"clone,omitempty"         jsonschema:"default=true,description=Enables injecting the default clone process.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-clone-key"`
		Environment  []string       `yaml:"environment,omitempty"   json:"environment,omitempty"   jsonschema:"description=Controls which containers processes can have global env injected.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-environment-key"`
		AutoCancel   *CancelOptions `yaml:"auto_cancel,omitempty"   json:"auto_cancel,omitempty"   jsonschema:"description=Enables auto canceling of queued or running pipelines that become stale due to new push.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-auto-cancel-key"`
	}

	// CancelOptions is the yaml representation of
	// the auto_cancel block for a pipeline.
	CancelOptions struct {
		Running       *bool `yaml:"running,omitempty"        json:"running,omitempty"        jsonschema:"description=Enables auto canceling of running pipelines that become stale due to new push.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-auto-cancel-key"`
		Pending       *bool `yaml:"pending,omitempty"        json:"pending,omitempty"        jsonschema:"description=Enables auto canceling of queued pipelines that become stale due to new push.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-auto-cancel-key"`
		DefaultBranch *bool `yaml:"default_branch,omitempty" json:"default_branch,omitempty" jsonschema:"description=Enables auto canceling of queued or running pipelines that become stale due to new push to default branch.\nReference: https://go-vela.github.io/docs/reference/yaml/metadata/#the-auto-cancel-key"`
	}
)

// ToPipeline converts the Metadata type
// to a pipeline Metadata type.
func (m *Metadata) ToPipeline() *pipeline.Metadata {
	var clone bool
	if m.Clone == nil {
		clone = true
	} else {
		clone = *m.Clone
	}

	autoCancel := new(pipeline.CancelOptions)

	// default to false for all fields if block isn't found
	if m.AutoCancel == nil {
		autoCancel.Pending = false
		autoCancel.Running = false
		autoCancel.DefaultBranch = false
	} else {
		// if block is found but pending field isn't, default to true
		if m.AutoCancel.Pending != nil {
			autoCancel.Pending = *m.AutoCancel.Pending
		} else {
			autoCancel.Pending = true
		}

		if m.AutoCancel.Running != nil {
			autoCancel.Running = *m.AutoCancel.Running
		}

		if m.AutoCancel.DefaultBranch != nil {
			autoCancel.DefaultBranch = *m.AutoCancel.DefaultBranch
		}
	}

	return &pipeline.Metadata{
		Template:    m.Template,
		Clone:       clone,
		Environment: m.Environment,
		AutoCancel:  autoCancel,
	}
}

// HasEnvironment checks if the container type
// is contained within the environment list.
func (m *Metadata) HasEnvironment(container string) bool {
	return slices.Contains(m.Environment, container)
}
