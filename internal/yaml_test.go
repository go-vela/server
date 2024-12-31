// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"os"
	"reflect"
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
		name         string
		file         string
		wantBuild    *yaml.Build
		wantWarnings []string
		wantErr      bool
	}{
		{
			name:      "go-yaml",
			file:      "testdata/go-yaml.yml",
			wantBuild: wantBuild,
		},
		{
			name:         "buildkite legacy",
			file:         "testdata/buildkite.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{"using legacy version. Upgrade to go-yaml v3"},
		},
		{
			name:         "anchor collapse",
			file:         "testdata/buildkite_new_version.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{"16:duplicate << keys in single YAML map"},
		},
		{
			name:      "no version",
			file:      "testdata/no_version.yml",
			wantBuild: wantBuild,
		},
		{
			name:      "invalid yaml",
			file:      "testdata/invalid.yml",
			wantBuild: nil,
			wantErr:   true,
		},
	}

	// run tests
	for _, test := range tests {
		bytes, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file for test %s: %v", test.name, err)
		}

		gotBuild, gotWarnings, err := ParseYAML(bytes)
		if err != nil && !test.wantErr {
			t.Errorf("ParseYAML for test %s returned err: %v", test.name, err)
		}

		if err == nil && test.wantErr {
			t.Errorf("ParseYAML for test %s returned nil error", test.name)
		}

		if err != nil && test.wantErr {
			continue
		}

		// different versions expected
		wantBuild.Version = gotBuild.Version

		if diff := cmp.Diff(gotBuild, test.wantBuild); diff != "" {
			t.Errorf("ParseYAML for test %s returned diff (-got +want):\n%s", test.name, diff)
		}

		if !reflect.DeepEqual(gotWarnings, test.wantWarnings) {
			t.Errorf("ParseYAML for test %s returned warnings %v, want %v", test.name, gotWarnings, test.wantWarnings)
		}
	}
}