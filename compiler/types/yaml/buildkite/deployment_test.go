// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Deployment_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		name       string
		deployment *Deployment
		want       *pipeline.Deployment
	}{
		{
			name: "deployment with template",
			deployment: &Deployment{
				Template: StepTemplate{Name: "foo"},
			},
			want: &pipeline.Deployment{},
		},
		{
			name: "deployment with targets and parameters",
			deployment: &Deployment{
				Targets: []string{"foo"},
				Parameters: ParameterMap{
					"foo": {
						Description: "bar",
						Type:        "string",
						Required:    true,
						Options:     []string{"baz"},
					},
					"bar": {
						Description: "baz",
						Type:        "string",
						Required:    false,
					},
				},
			},
			want: &pipeline.Deployment{
				Targets: []string{"foo"},
				Parameters: pipeline.ParameterMap{
					"foo": {
						Description: "bar",
						Type:        "string",
						Required:    true,
						Options:     []string{"baz"},
					},
					"bar": {
						Description: "baz",
						Type:        "string",
						Required:    false,
					},
				},
			},
		},
		{
			name:       "empty deployment config",
			deployment: &Deployment{},
			want:       &pipeline.Deployment{},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.deployment.ToPipeline()

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("ToPipeline for %s does not match: -want +got):\n%s", test.name, diff)
		}
	}
}
