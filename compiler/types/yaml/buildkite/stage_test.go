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

func TestYaml_StageSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		stages *StageSlice
		want   *pipeline.StageSlice
	}{
		{
			stages: &StageSlice{
				{
					Name:  "echo",
					Needs: []string{"clone"},
					Steps: StepSlice{
						{
							Commands:    []string{"echo hello"},
							Detach:      false,
							Entrypoint:  []string{"/bin/sh"},
							Environment: map[string]string{"FOO": "bar"},
							Image:       "alpine:latest",
							Name:        "echo",
							Privileged:  false,
							Pull:        "not_present",
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
				},
			},
			want: &pipeline.StageSlice{
				{
					Name:  "echo",
					Needs: []string{"clone"},
					Steps: pipeline.ContainerSlice{
						{
							Commands:    []string{"echo hello"},
							Detach:      false,
							Entrypoint:  []string{"/bin/sh"},
							Environment: map[string]string{"FOO": "bar"},
							Image:       "alpine:latest",
							Name:        "echo",
							Privileged:  false,
							Pull:        "not_present",
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
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.stages.ToPipeline()

		// WARNING: hack to compare stages
		//
		// Channel values can only be compared for equality.
		// Two channel values are considered equal if they
		// originated from the same make call meaning they
		// refer to the same channel value in memory.
		for i, stage := range *got {
			tmp := *test.want

			tmp[i].Done = stage.Done
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_StageSlice_UnmarshalYAML(t *testing.T) {
	// setup types
	var (
		b   []byte
		err error
	)

	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StageSlice
	}{
		{
			failure: false,
			file:    "testdata/stage.yml",
			want: &StageSlice{
				{
					Name:  "dependencies",
					Needs: []string{"clone"},
					Environment: map[string]string{
						"STAGE_ENV_VAR": "stage",
					},
					Independent: true,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew downloadDependencies"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Image: "openjdk:latest",
							Name:  "install",
							Pull:  "always",
						},
					},
				},
				{
					Name:  "test",
					Needs: []string{"dependencies", "clone"},
					Environment: map[string]string{
						"STAGE_ENV_VAR":    "stage",
						"SECOND_STAGE_ENV": "stage2",
					},
					Independent: false,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew check"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Name:  "test",
							Image: "openjdk:latest",
							Pull:  "always",
						},
					},
				},
				{
					Name:  "build",
					Needs: []string{"dependencies", "clone"},
					Environment: map[string]string{
						"STAGE_ENV_VAR": "stage",
					},
					Independent: false,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew build"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Name:  "build",
							Image: "openjdk:latest",
							Pull:  "always",
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
			file:    "",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StageSlice)

		if len(test.file) > 0 {
			b, err = os.ReadFile(test.file)
			if err != nil {
				t.Errorf("unable to read file: %v", err)
			}
		} else {
			b = []byte("- foo")
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

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("(Unmarshal mismatch: -want +got):\n%s", diff)
		}
	}
}

func TestYaml_StageSlice_MarshalYAML(t *testing.T) {
	// setup types
	var (
		b   []byte
		err error
	)

	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StageSlice
	}{
		{
			failure: false,
			file:    "testdata/stage.yml",
			want: &StageSlice{
				{
					Name:        "dependencies",
					Needs:       []string{"clone"},
					Independent: true,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew downloadDependencies"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Image: "openjdk:latest",
							Name:  "install",
							Pull:  "always",
						},
					},
				},
				{
					Name:        "test",
					Needs:       []string{"dependencies", "clone"},
					Independent: false,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew check"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Name:  "test",
							Image: "openjdk:latest",
							Pull:  "always",
						},
					},
				},
				{
					Name:        "build",
					Needs:       []string{"dependencies", "clone"},
					Independent: false,
					Steps: StepSlice{
						{
							Commands: []string{"./gradlew build"},
							Environment: map[string]string{
								"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
								"GRADLE_USER_HOME": ".gradle",
							},
							Name:  "build",
							Image: "openjdk:latest",
							Pull:  "always",
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
			file:    "",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StageSlice)
		got2 := new(StageSlice)

		if len(test.file) > 0 {
			b, err = os.ReadFile(test.file)
			if err != nil {
				t.Errorf("unable to read file: %v", err)
			}
		} else {
			b = []byte("- foo")
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

		out, err := yaml.Marshal(got)
		if err != nil {
			t.Errorf("MarshalYAML returned err: %v", err)
		}

		err = yaml.Unmarshal(out, got2)
		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if diff := cmp.Diff(got2, test.want); diff != "" {
			t.Errorf("(Marshal mismatch: -got +want):\n%s", diff)
		}
	}
}

func TestYaml_Stage_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		stage       *Stage
		environment map[string]string
		failure     bool
	}{
		{
			stage: &Stage{
				Environment: map[string]string{"FOO": "bar"},
				Name:        "testStage",
			},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			stage:       &Stage{},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			stage:       nil,
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			stage: &Stage{
				Environment: map[string]string{"FOO": "bar"},
				Name:        "testStage",
			},
			environment: nil,
			failure:     true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.stage.MergeEnv(test.environment)

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
