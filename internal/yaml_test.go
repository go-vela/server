// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

func TestInternal_ParseYAML(t *testing.T) {
	// wantBuild
	wantBuild := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Name:  "example",
				Image: "alpine:latest",
				Environment: map[string]string{
					"REGION": "dev",
				},
				Pull: "not_present",
				Commands: []string{
					"echo $REGION",
				},
			},
		},
	}

	// set up tests
	tests := []struct {
		file    string
		want    *yaml.Build
		wantErr bool
	}{
		{
			file: "testdata/go-yaml.yml",
			want: wantBuild,
		},
		{
			file: "testdata/buildkite.yml",
			want: wantBuild,
		},
		{
			file: "testdata/no_version.yml",
			want: wantBuild,
		},
		{
			file:    "testdata/invalid.yml",
			want:    nil,
			wantErr: true,
		},
	}

	// run tests
	for _, test := range tests {
		bytes, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
		}

		gotBuild, err := ParseYAML(bytes)
		if err != nil && !test.wantErr {
			t.Errorf("ParseYAML returned err: %v", err)
		}

		if err == nil && test.wantErr {
			t.Errorf("ParseYAML returned nil error")
		}

		if err != nil && test.wantErr {
			continue
		}

		// different versions expected
		wantBuild.Version = gotBuild.Version

		if diff := cmp.Diff(gotBuild, test.want); diff != "" {
			t.Errorf("ParseYAML returned diff (-got +want):\n%s", diff)
		}
	}
}
