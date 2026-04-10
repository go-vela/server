// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func boolPtr(b bool) *bool { return &b }

func TestYaml_Artifacts_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		name      string
		artifacts *Artifacts
		want      *pipeline.Artifacts
	}{
		{
			name:      "nil Secured defaults to true",
			artifacts: &Artifacts{Paths: []string{"test-results.xml"}},
			want:      &pipeline.Artifacts{Paths: []string{"test-results.xml"}, Secured: true},
		},
		{
			name:      "Secured explicitly true",
			artifacts: &Artifacts{Paths: []string{"coverage.html"}, Secured: boolPtr(true)},
			want:      &pipeline.Artifacts{Paths: []string{"coverage.html"}, Secured: true},
		},
		{
			name:      "Secured explicitly false",
			artifacts: &Artifacts{Paths: []string{"junit-report.json"}, Secured: boolPtr(false)},
			want:      &pipeline.Artifacts{Paths: []string{"junit-report.json"}, Secured: false},
		},
		{
			name:      "no paths",
			artifacts: &Artifacts{},
			want:      &pipeline.Artifacts{Paths: nil, Secured: true},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.artifacts.ToPipeline()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ToPipeline mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
