// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
)

// Git is the yaml representation of git configurations for a pipeline.
type Git struct {
	Token `yaml:"token,omitempty" json:"token" jsonschema:"description=Provide the git token specifications, primarily used for cloning.\nReference: https://go-vela.github.io/docs/reference/yaml/git/#token"`
}

// Token is the yaml representation of the git token.
// Only applies when using GitHub App installations.
type Token struct {
	Repositories []string          `yaml:"repositories,omitempty" json:"repositories,omitempty" jsonschema:"description=Provide a list of repositories to clone.\nReference: https://go-vela.github.io/docs/reference/yaml/git/#repositories"`
	Permissions  map[string]string `yaml:"permissions,omitempty"  json:"permissions,omitempty"  jsonschema:"description=Provide a list of repository permissions to apply to the git token.\nReference: https://go-vela.github.io/docs/reference/yaml/git/#permissions"`
}

// ToPipeline converts the Git type
// to a pipeline Git type.
func (g *Git) ToPipeline() *pipeline.Git {
	for i, repo := range g.Repositories {
		split := strings.Split(repo, "/")
		if len(split) == 2 {
			g.Repositories[i] = split[1]
		}
	}

	return &pipeline.Git{
		Token: &pipeline.Token{
			Repositories: g.Repositories,
			Permissions:  g.Permissions,
		},
	}
}
