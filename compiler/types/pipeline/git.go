// SPDX-License-Identifier: Apache-2.0

package pipeline

// Git is the pipeline representation of git configurations for a pipeline.
//
// swagger:model PipelineGit
type Git struct {
	Token *Token `json:"token,omitempty" yaml:"token,omitempty"`
}

// Token is the pipeline representation of git token access configurations for a pipeline.
//
// swagger:model PipelineGitToken
type Token struct {
	Repositories []string          `json:"repositories,omitempty" yaml:"repositories,omitempty"`
	Permissions  map[string]string `json:"permissions,omitempty"  yaml:"permissions,omitempty"`
}

// Empty returns true if the provided struct is empty.
func (g *Git) Empty() bool {
	// return true if every field is empty
	if g.Token != nil {
		if g.Token.Repositories != nil {
			return false
		}

		if g.Token.Permissions != nil {
			return false
		}
	}

	// return false if any of the fields are provided
	return true
}
