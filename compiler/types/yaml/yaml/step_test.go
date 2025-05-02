// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
)

func TestYaml_StepSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		steps *StepSlice
		want  *pipeline.ContainerSlice
	}{
		{
			steps: &StepSlice{
				{
					Commands:    []string{"echo hello"},
					Detach:      false,
					Entrypoint:  []string{"/bin/sh"},
					Environment: map[string]string{"FOO": "bar"},
					Image:       "alpine:latest",
					Name:        "echo",
					Privileged:  false,
					Pull:        "not_present",
					ReportAs:    "my-step",
					IDRequest:   "yes",
					Ruleset: Ruleset{
						If: Rules{
							Branch:   []string{"main"},
							Comment:  []string{"test comment"},
							Event:    []string{"push"},
							Path:     []string{"foo.txt"},
							Repo:     []string{"github/octocat"},
							Status:   []string{"success"},
							Tag:      []string{"v0.1.0"},
							Target:   []string{"production"},
							Operator: "and",
						},
						Unless: Rules{
							Branch:   []string{"main"},
							Comment:  []string{"real comment"},
							Event:    []string{"pull_request"},
							Path:     []string{"bar.txt"},
							Repo:     []string{"github/octocat"},
							Status:   []string{"failure"},
							Tag:      []string{"v0.2.0"},
							Target:   []string{"production"},
							Operator: "and",
						},
						Continue: false,
					},
					Secrets: StepSecretSlice{
						{
							Source: "docker_username",
							Target: "plugin_username",
						},
					},
					Ulimits: UlimitSlice{
						{
							Name: "foo",
							Soft: 1024,
							Hard: 2048,
						},
					},
					Volumes: VolumeSlice{
						{
							Source:      "/foo",
							Destination: "/bar",
							AccessMode:  "ro",
						},
					},
				},
			},
			want: &pipeline.ContainerSlice{
				{
					Commands:    []string{"echo hello"},
					Detach:      false,
					Entrypoint:  []string{"/bin/sh"},
					Environment: map[string]string{"FOO": "bar"},
					Image:       "alpine:latest",
					Name:        "echo",
					Privileged:  false,
					Pull:        "not_present",
					ReportAs:    "my-step",
					IDRequest:   "yes",
					Ruleset: pipeline.Ruleset{
						If: pipeline.Rules{
							Branch:   []string{"main"},
							Comment:  []string{"test comment"},
							Event:    []string{"push"},
							Path:     []string{"foo.txt"},
							Repo:     []string{"github/octocat"},
							Status:   []string{"success"},
							Tag:      []string{"v0.1.0"},
							Target:   []string{"production"},
							Operator: "and",
						},
						Unless: pipeline.Rules{
							Branch:   []string{"main"},
							Comment:  []string{"real comment"},
							Event:    []string{"pull_request"},
							Path:     []string{"bar.txt"},
							Repo:     []string{"github/octocat"},
							Status:   []string{"failure"},
							Tag:      []string{"v0.2.0"},
							Target:   []string{"production"},
							Operator: "and",
						},
						Continue: false,
					},
					Secrets: pipeline.StepSecretSlice{
						{
							Source: "docker_username",
							Target: "plugin_username",
						},
					},
					Ulimits: pipeline.UlimitSlice{
						{
							Name: "foo",
							Soft: 1024,
							Hard: 2048,
						},
					},
					Volumes: pipeline.VolumeSlice{
						{
							Source:      "/foo",
							Destination: "/bar",
							AccessMode:  "ro",
						},
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.steps.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_StepSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StepSlice
	}{
		{
			failure: false,
			file:    "testdata/step.yml",
			want: &StepSlice{
				{
					Commands: raw.StringSlice{"./gradlew downloadDependencies"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Name:  "install",
					Image: "openjdk:latest",
					Pull:  "always",
				},
				{
					Commands: raw.StringSlice{"./gradlew check"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Name:  "test",
					Image: "openjdk:latest",
					Pull:  "always",
				},
				{
					Commands: raw.StringSlice{"./gradlew build"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Name:  "build",
					Image: "openjdk:latest",
					Pull:  "always",
				},
				{
					Name:     "docker_build",
					Image:    "plugins/docker:18.09",
					Pull:     "always",
					ReportAs: "docker",
					Parameters: map[string]interface{}{
						"registry": "index.docker.io",
						"repo":     "github/octocat",
						"tags":     []interface{}{"latest", "dev"},
					},
				},
				{
					Name: "templated_publish",
					Pull: "not_present",
					Template: StepTemplate{
						Name: "docker_publish",
						Variables: map[string]interface{}{
							"registry": "index.docker.io",
							"repo":     "github/octocat",
							"tags":     []interface{}{"latest", "dev"},
						},
					},
				},
			},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/step_malformed.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/step_nil.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StepSlice)

		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
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

func TestYaml_Step_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		step        *Step
		environment map[string]string
		failure     bool
	}{
		{
			step: &Step{
				Commands:    []string{"echo hello"},
				Detach:      false,
				Entrypoint:  []string{"/bin/sh"},
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Privileged:  false,
				Pull:        "not_present",
			},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			step:        &Step{},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			step:        nil,
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			step: &Step{
				Commands:    []string{"echo hello"},
				Detach:      false,
				Entrypoint:  []string{"/bin/sh"},
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Privileged:  false,
				Pull:        "not_present",
			},
			environment: nil,
			failure:     true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.step.MergeEnv(test.environment)

		if test.failure {
			if err == nil {
				t.Errorf("MergeEnv should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("MergeEnv returned err: %v", err)
		}
	}
}
