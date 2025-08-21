// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"os"
	"reflect"
	"testing"

	"github.com/buildkite/yaml"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Ruleset_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		ruleset *Ruleset
		want    *pipeline.Ruleset
	}{
		{
			ruleset: &Ruleset{
				If: Rules{
					Branch:   []string{"main"},
					Comment:  []string{"test comment"},
					Event:    []string{"push", "pull_request:labeled"},
					Path:     []string{"foo.txt"},
					Repo:     []string{"github/octocat"},
					Sender:   []string{"octocat"},
					Status:   []string{"success"},
					Tag:      []string{"v0.1.0"},
					Target:   []string{"production"},
					Label:    []string{"enhancement"},
					Instance: []string{"http://localhost:8080"},
					Matcher:  "filepath",
					Operator: "and",
				},
				Unless: Rules{
					Branch:   []string{"main"},
					Comment:  []string{"real comment"},
					Event:    []string{"pull_request"},
					Path:     []string{"bar.txt"},
					Repo:     []string{"github/octocat"},
					Sender:   []string{"octokitty"},
					Status:   []string{"failure"},
					Tag:      []string{"v0.2.0"},
					Target:   []string{"production"},
					Instance: []string{"http://localhost:8080"},
					Matcher:  "regexp",
					Operator: "or",
				},
				Continue: false,
			},
			want: &pipeline.Ruleset{
				If: pipeline.Rules{
					Branch:   []string{"main"},
					Comment:  []string{"test comment"},
					Event:    []string{"push", "pull_request:labeled"},
					Path:     []string{"foo.txt"},
					Repo:     []string{"github/octocat"},
					Sender:   []string{"octocat"},
					Status:   []string{"success"},
					Tag:      []string{"v0.1.0"},
					Target:   []string{"production"},
					Label:    []string{"enhancement"},
					Instance: []string{"http://localhost:8080"},
					Matcher:  "filepath",
					Operator: "and",
				},
				Unless: pipeline.Rules{
					Branch:   []string{"main"},
					Comment:  []string{"real comment"},
					Event:    []string{"pull_request"},
					Path:     []string{"bar.txt"},
					Repo:     []string{"github/octocat"},
					Sender:   []string{"octokitty"},
					Status:   []string{"failure"},
					Tag:      []string{"v0.2.0"},
					Target:   []string{"production"},
					Instance: []string{"http://localhost:8080"},
					Matcher:  "regexp",
					Operator: "or",
				},
				Continue: false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.ruleset.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_Ruleset_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		file string
		want *Ruleset
	}{
		{
			file: "testdata/ruleset_simple.yml",
			want: &Ruleset{
				If: Rules{
					Branch:   []string{"main"},
					Comment:  []string{"test comment"},
					Event:    []string{"push"},
					Instance: []string{"vela-server"},
					Label:    []string{"bug"},
					Path:     []string{"foo.txt"},
					Repo:     []string{"github/octocat"},
					Sender:   []string{"octocat"},
					Status:   []string{"success"},
					Tag:      []string{"v0.1.0"},
					Target:   []string{"production"},
					Matcher:  "filepath",
					Operator: "and",
				},
				Continue: true,
			},
		},
		{
			file: "testdata/ruleset_advanced.yml",
			want: &Ruleset{
				If: Rules{
					Branch:   []string{"main"},
					Event:    []string{"push"},
					Tag:      []string{"^refs/tags/(\\d+\\.)+\\d+$"},
					Matcher:  "regexp",
					Operator: "or",
				},
				Unless: Rules{
					Event:    []string{"deployment:created", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened", "comment:created", "comment:edited", "schedule"},
					Path:     []string{"foo.txt", "/foo/bar.txt"},
					Matcher:  "regexp",
					Operator: "or",
				},
				Continue: true,
			},
		},
		{
			file: "testdata/ruleset_op_match.yml",
			want: &Ruleset{
				If: Rules{
					Branch:   []string{"main"},
					Event:    []string{"push"},
					Tag:      []string{"^refs/tags/(\\d+\\.)+\\d+$"},
					Matcher:  "regexp",
					Operator: "and",
				},
				Unless: Rules{
					Event:    []string{"deployment:created", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened", "comment:created", "comment:edited", "schedule"},
					Path:     []string{"foo.txt", "/foo/bar.txt"},
					Matcher:  "filepath",
					Operator: "or",
				},
				Continue: true,
			},
		},
		{
			file: "testdata/ruleset_regex.yml",
			want: &Ruleset{
				If: Rules{
					Branch:   []string{"main"},
					Event:    []string{"tag"},
					Tag:      []string{"^refs/tags/(\\d+\\.)+\\d+$"},
					Operator: "and",
					Matcher:  "regex",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := new(Ruleset)

		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
		}

		err = yaml.Unmarshal(b, got)
		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("UnmarshalYAML mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestYaml_Rules_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		rules *Rules
		want  *pipeline.Rules
	}{
		{
			rules: &Rules{
				Branch:   []string{"main"},
				Comment:  []string{"test comment"},
				Event:    []string{"push", "pull_request:labeled"},
				Instance: []string{"vela-server"},
				Path:     []string{"foo.txt"},
				Repo:     []string{"github/octocat"},
				Sender:   []string{"octocat"},
				Status:   []string{"success"},
				Tag:      []string{"v0.1.0"},
				Target:   []string{"production"},
				Label:    []string{"enhancement"},
			},
			want: &pipeline.Rules{
				Branch:   []string{"main"},
				Comment:  []string{"test comment"},
				Event:    []string{"push", "pull_request:labeled"},
				Instance: []string{"vela-server"},
				Path:     []string{"foo.txt"},
				Repo:     []string{"github/octocat"},
				Sender:   []string{"octocat"},
				Status:   []string{"success"},
				Tag:      []string{"v0.1.0"},
				Target:   []string{"production"},
				Label:    []string{"enhancement"},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.rules.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_Rules_UnmarshalYAML(t *testing.T) {
	// setup types
	var (
		b   []byte
		err error
	)

	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *Rules
	}{
		{
			failure: false,
			file:    "testdata/ruleset_simple.yml",
			want: &Rules{
				Branch:   []string{"main"},
				Comment:  []string{"test comment"},
				Event:    []string{"push"},
				Instance: []string{"vela-server"},
				Label:    []string{"bug"},
				Path:     []string{"foo.txt"},
				Repo:     []string{"github/octocat"},
				Sender:   []string{"octocat"},
				Status:   []string{"success"},
				Tag:      []string{"v0.1.0"},
				Target:   []string{"production"},
			},
		},
		{
			failure: true,
			file:    "",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(Rules)

		if len(test.file) > 0 {
			b, err = os.ReadFile(test.file)
			if err != nil {
				t.Errorf("unable to read file: %v", err)
			}
		} else {
			b = []byte("``")
		}

		err = yaml.Unmarshal(b, got)

		if test.failure {
			if err == nil {
				t.Errorf("UnmarshalYAML should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalYAML is %v, want %v", got, test.want)
		}
	}
}
