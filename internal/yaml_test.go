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
				Parameters: map[string]interface{}{
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
		name          string
		file          string
		wantBuild     *yaml.Build
		wantWarnings  []string
		warningPrefix string
		wantErr       bool
	}{
		{
			name:      "go-yaml",
			file:      "testdata/go-yaml.yml",
			wantBuild: wantBuild,
		},
		{
			name:         "top level anchors",
			file:         "testdata/top_level_anchor.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{`6:duplicate << keys in single YAML map`},
		},
		{
			name:         "top level anchors legacy",
			file:         "testdata/top_level_anchor_legacy.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{`using legacy version - address any incompatibilities and use "1" instead`},
		},
		{
			name:         "buildkite legacy",
			file:         "testdata/buildkite.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{`using legacy version - address any incompatibilities and use "1" instead`},
		},
		{
			name:         "anchor collapse",
			file:         "testdata/buildkite_new_version.yml",
			wantBuild:    wantBuild,
			wantWarnings: []string{"16:duplicate << keys in single YAML map"},
		},
		{
			name:          "anchor collapse - warning prefix",
			file:          "testdata/buildkite_new_version.yml",
			wantBuild:     wantBuild,
			wantWarnings:  []string{"prefix:16:duplicate << keys in single YAML map"},
			warningPrefix: "prefix",
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

		gotBuild, gotWarnings, err := ParseYAML(bytes, test.warningPrefix)
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
