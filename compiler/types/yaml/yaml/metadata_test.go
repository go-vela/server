// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Metadata_ToPipeline(t *testing.T) {
	tBool := true
	fBool := false
	// setup tests
	tests := []struct {
		metadata *Metadata
		want     *pipeline.Metadata
	}{
		{
			metadata: &Metadata{
				Template:    false,
				Clone:       &fBool,
				Environment: []string{"steps", "services", "secrets"},
				AutoCancel: &CancelOptions{
					Pending:       &tBool,
					Running:       &tBool,
					DefaultBranch: &fBool,
				},
			},
			want: &pipeline.Metadata{
				Template:    false,
				Clone:       false,
				Environment: []string{"steps", "services", "secrets"},
				AutoCancel: &pipeline.CancelOptions{
					Pending:       true,
					Running:       true,
					DefaultBranch: false,
				},
			},
		},
		{
			metadata: &Metadata{
				Template:    false,
				Clone:       &tBool,
				Environment: []string{"steps", "services"},
			},
			want: &pipeline.Metadata{
				Template:    false,
				Clone:       true,
				Environment: []string{"steps", "services"},
				AutoCancel: &pipeline.CancelOptions{
					Pending:       false,
					Running:       false,
					DefaultBranch: false,
				},
			},
		},
		{
			metadata: &Metadata{
				Template:    false,
				Clone:       nil,
				Environment: []string{"steps"},
				AutoCancel: &CancelOptions{
					Running:       &tBool,
					DefaultBranch: &tBool,
				},
			},
			want: &pipeline.Metadata{
				Template:    false,
				Clone:       true,
				Environment: []string{"steps"},
				AutoCancel: &pipeline.CancelOptions{
					Pending:       true,
					Running:       true,
					DefaultBranch: true,
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.metadata.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_Metadata_HasEnvironment(t *testing.T) {
	// setup tests
	tests := []struct {
		metadata  *Metadata
		container string
		want      bool
	}{
		{
			metadata: &Metadata{
				Environment: []string{"steps", "services", "secrets"},
			},
			container: "steps",
			want:      true,
		},
		{
			metadata: &Metadata{
				Environment: []string{"services", "secrets"},
			},
			container: "services",
			want:      true,
		},
		{
			metadata: &Metadata{
				Environment: []string{"steps", "services", "secrets"},
			},
			container: "notacontainer",
			want:      false,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.metadata.HasEnvironment(test.container)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}
