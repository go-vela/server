// SPDX-License-Identifier: Apache-2.0

package pipeline

// Git is the pipeline representation of the git block for a pipeline.
//
// swagger:model PipelineGit
type Git struct {
	Access *Access `json:"access,omitempty" yaml:"access,omitempty"`
}

type Access struct {
	Repositories []string `json:"repositories,omitempty" yaml:"repositories,omitempty"`
}
