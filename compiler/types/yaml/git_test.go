// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Git_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		git  *Git
		want *pipeline.Git
	}{
		{
			git: &Git{
				Token: Token{
					Repositories: []string{"foo", "bar"},
				},
			},
			want: &pipeline.Git{
				Token: &pipeline.Token{
					Repositories: []string{"foo", "bar"},
				},
			},
		},
		{
			git: &Git{
				Token: Token{
					Permissions: map[string]string{"foo": "bar"},
				},
			},
			want: &pipeline.Git{
				Token: &pipeline.Token{
					Permissions: map[string]string{"foo": "bar"},
				},
			},
		},
		{
			git: &Git{
				Token: Token{
					Repositories: []string{"foo/bar"},
				},
			},
			want: &pipeline.Git{
				Token: &pipeline.Token{
					Repositories: []string{"bar"},
				},
			},
		},
		{
			git: &Git{
				Token: Token{},
			},
			want: &pipeline.Git{
				Token: &pipeline.Token{},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.git.ToPipeline()

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("ToPipeline is %s", diff)
		}
	}
}
