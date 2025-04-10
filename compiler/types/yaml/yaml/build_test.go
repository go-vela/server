// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
)

func TestYaml_Build_ToAPI(t *testing.T) {
	build := new(api.Pipeline)
	build.SetFlavor("16cpu8gb")
	build.SetPlatform("gcp")
	build.SetVersion("1")
	build.SetExternalSecrets(true)
	build.SetInternalSecrets(true)
	build.SetServices(true)
	build.SetStages(false)
	build.SetSteps(true)
	build.SetTemplates(true)

	stages := new(api.Pipeline)
	stages.SetFlavor("")
	stages.SetPlatform("")
	stages.SetVersion("1")
	stages.SetExternalSecrets(false)
	stages.SetInternalSecrets(false)
	stages.SetServices(false)
	stages.SetStages(true)
	stages.SetSteps(false)
	stages.SetTemplates(false)

	steps := new(api.Pipeline)
	steps.SetFlavor("")
	steps.SetPlatform("")
	steps.SetVersion("1")
	steps.SetExternalSecrets(false)
	steps.SetInternalSecrets(false)
	steps.SetServices(false)
	steps.SetStages(false)
	steps.SetSteps(true)
	steps.SetTemplates(false)

	// setup tests
	tests := []struct {
		name string
		file string
		want *api.Pipeline
	}{
		{
			name: "build",
			file: "testdata/build.yml",
			want: build,
		},
		{
			name: "stages",
			file: "testdata/build_anchor_stage.yml",
			want: stages,
		},
		{
			name: "steps",
			file: "testdata/build_anchor_step.yml",
			want: steps,
		},
	}

	// run tests
	for _, test := range tests {
		b := new(Build)

		data, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file %s for %s: %v", test.file, test.name, err)
		}

		err = yaml.Unmarshal(data, b)
		if err != nil {
			t.Errorf("unable to unmarshal YAML for %s: %v", test.name, err)
		}

		got := b.ToPipelineAPI()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipelineAPI for %s is %v, want %v", test.name, got, test.want)
		}
	}
}

func TestYaml_Build_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		file string
		want *Build
	}{
		{
			file: "testdata/build.yml",
			want: &Build{
				Version: "1",
				Metadata: Metadata{
					Template:    false,
					Clone:       nil,
					Environment: []string{"steps", "services", "secrets"},
				},
				Environment: raw.StringSliceMap{
					"HELLO": "Hello, Global Message",
				},
				Worker: Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Services: ServiceSlice{
					{
						Ports: []string{"5432:5432"},
						Environment: raw.StringSliceMap{
							"POSTGRES_DB": "foo",
						},
						Name:  "postgres",
						Image: "postgres:latest",
						Pull:  "not_present",
					},
				},
				Steps: StepSlice{
					{
						Commands: raw.StringSlice{"./gradlew downloadDependencies"},
						Environment: raw.StringSliceMap{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "install",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:edited"}},
							Matcher:  "filepath",
							Operator: "and",
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
					{
						Commands: raw.StringSlice{"./gradlew check"},
						Environment: raw.StringSliceMap{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Name:  "test",
						Image: "openjdk:latest",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Volumes: VolumeSlice{
							{
								Source:      "/foo",
								Destination: "/bar",
								AccessMode:  "ro",
							},
						},
						Ulimits: UlimitSlice{
							{
								Name: "foo",
								Soft: 1024,
								Hard: 2048,
							},
						},
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
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Volumes: VolumeSlice{
							{
								Source:      "/foo",
								Destination: "/bar",
								AccessMode:  "ro",
							},
						},
						Ulimits: UlimitSlice{
							{
								Name: "foo",
								Soft: 1024,
								Hard: 2048,
							},
						},
					},
					{
						Name: "docker_build",
						Parameters: map[string]any{
							"dry_run":  true,
							"registry": "index.docker.io",
							"repo":     "github/octocat",
							"tags":     []any{"latest", "dev"},
						},
						Image: "plugins/docker:18.09",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
					},
					{
						Name: "docker_publish",
						Parameters: map[string]any{
							"registry": "index.docker.io",
							"repo":     "github/octocat",
							"tags":     []any{"latest", "dev"},
						},
						Image: "plugins/docker:18.09",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Branch: []string{"main"}, Event: []string{"push"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Secrets: StepSecretSlice{
							{
								Source: "docker_username",
								Target: "PLUGIN_USERNAME",
							},
							{
								Source: "docker_password",
								Target: "PLUGIN_PASSWORD",
							},
						},
					},
				},
				Secrets: SecretSlice{
					{
						Name:   "docker_username",
						Key:    "org/repo/docker/username",
						Engine: "native",
						Type:   "repo",
						Pull:   "build_start",
					},
					{
						Name:   "docker_password",
						Key:    "org/repo/docker/password",
						Engine: "vault",
						Type:   "repo",
						Pull:   "build_start",
					},
					{
						Name:   "docker_username",
						Key:    "org/docker/username",
						Engine: "native",
						Type:   "org",
						Pull:   "build_start",
					},
					{
						Name:   "docker_password",
						Key:    "org/docker/password",
						Engine: "vault",
						Type:   "org",
						Pull:   "build_start",
					},
					{
						Name:   "docker_username",
						Key:    "org/team/docker/username",
						Engine: "native",
						Type:   "shared",
						Pull:   "build_start",
					},
					{
						Name:   "docker_password",
						Key:    "org/team/docker/password",
						Engine: "vault",
						Type:   "shared",
						Pull:   "build_start",
					},
					{
						Origin: Origin{
							Image: "target/vela-vault:latest",
							Parameters: map[string]any{
								"addr": "vault.example.com",
							},
							Pull: "always",
							Secrets: StepSecretSlice{
								{
									Source: "docker_username",
									Target: "DOCKER_USERNAME",
								},
								{
									Source: "docker_password",
									Target: "DOCKER_PASSWORD",
								},
							},
						},
					},
				},
				Templates: TemplateSlice{
					{
						Name:   "docker_publish",
						Source: "github.com/go-vela/atlas/stable/docker_publish",
						Type:   "github",
					},
				},
			},
		},
		{
			file: "testdata/build_anchor_stage.yml",
			want: &Build{
				Version: "1",
				Metadata: Metadata{
					Template:    false,
					Clone:       nil,
					Environment: []string{"steps", "services", "secrets"},
				},
				Stages: StageSlice{
					{
						Name:        "dependencies",
						Needs:       []string{"clone"},
						Independent: false,
						Steps: StepSlice{
							{
								Commands: raw.StringSlice{"./gradlew downloadDependencies"},
								Environment: raw.StringSliceMap{
									"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
									"GRADLE_USER_HOME": ".gradle",
								},
								Image: "openjdk:latest",
								Name:  "install",
								Pull:  "always",
								Ruleset: Ruleset{
									If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
									Matcher:  "filepath",
									Operator: "and",
								},
								Volumes: VolumeSlice{
									{
										Source:      "/foo",
										Destination: "/bar",
										AccessMode:  "ro",
									},
								},
								Ulimits: UlimitSlice{
									{
										Name: "foo",
										Soft: 1024,
										Hard: 2048,
									},
								},
							},
						},
					},
					{
						Name:        "test",
						Needs:       []string{"dependencies", "clone"},
						Independent: false,
						Steps: StepSlice{
							{
								Commands: raw.StringSlice{"./gradlew check"},
								Environment: raw.StringSliceMap{
									"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
									"GRADLE_USER_HOME": ".gradle",
								},
								Name:  "test",
								Image: "openjdk:latest",
								Pull:  "always",
								Ruleset: Ruleset{
									If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
									Matcher:  "filepath",
									Operator: "and",
								},
								Volumes: VolumeSlice{
									{
										Source:      "/foo",
										Destination: "/bar",
										AccessMode:  "ro",
									},
								},
								Ulimits: UlimitSlice{
									{
										Name: "foo",
										Soft: 1024,
										Hard: 2048,
									},
								},
							},
						},
					},
					{
						Name:        "build",
						Needs:       []string{"dependencies", "clone"},
						Independent: true,
						Steps: StepSlice{
							{
								Commands: raw.StringSlice{"./gradlew build"},
								Environment: raw.StringSliceMap{
									"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
									"GRADLE_USER_HOME": ".gradle",
								},
								Name:  "build",
								Image: "openjdk:latest",
								Pull:  "always",
								Ruleset: Ruleset{
									If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
									Matcher:  "filepath",
									Operator: "and",
								},
								Volumes: VolumeSlice{
									{
										Source:      "/foo",
										Destination: "/bar",
										AccessMode:  "ro",
									},
								},
								Ulimits: UlimitSlice{
									{
										Name: "foo",
										Soft: 1024,
										Hard: 2048,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			file: "testdata/build_anchor_step.yml",
			want: &Build{
				Version: "1",
				Metadata: Metadata{
					Template:    false,
					Clone:       nil,
					Environment: []string{"steps", "services", "secrets"},
				},
				Steps: StepSlice{
					{
						Commands: raw.StringSlice{"./gradlew downloadDependencies"},
						Environment: raw.StringSliceMap{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "install",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Volumes: VolumeSlice{
							{
								Source:      "/foo",
								Destination: "/bar",
								AccessMode:  "ro",
							},
						},
						Ulimits: UlimitSlice{
							{
								Name: "foo",
								Soft: 1024,
								Hard: 2048,
							},
						},
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
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Volumes: VolumeSlice{
							{
								Source:      "/foo",
								Destination: "/bar",
								AccessMode:  "ro",
							},
						},
						Ulimits: UlimitSlice{
							{
								Name: "foo",
								Soft: 1024,
								Hard: 2048,
							},
						},
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
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
						},
						Volumes: VolumeSlice{
							{
								Source:      "/foo",
								Destination: "/bar",
								AccessMode:  "ro",
							},
						},
						Ulimits: UlimitSlice{
							{
								Name: "foo",
								Soft: 1024,
								Hard: 2048,
							},
						},
					},
				},
			},
		},
		{
			file: "testdata/build_empty_env.yml",
			want: &Build{
				Version: "1",
				Metadata: Metadata{
					Template:    false,
					Clone:       nil,
					Environment: []string{},
				},
				Environment: raw.StringSliceMap{
					"HELLO": "Hello, Global Message",
				},
				Worker: Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Steps: StepSlice{
					{
						Commands: raw.StringSlice{"./gradlew downloadDependencies"},
						Environment: raw.StringSliceMap{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "install",
						Pull:  "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened"}},
							Matcher:  "filepath",
							Operator: "and",
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
		{
			file: "testdata/merge_anchor.yml",
			want: &Build{
				Version: "1",
				Metadata: Metadata{
					Template:    false,
					Clone:       nil,
					Environment: []string{"steps", "services", "secrets"},
				},
				Services: ServiceSlice{
					{
						Name:  "service-a",
						Ports: []string{"5432:5432"},
						Environment: raw.StringSliceMap{
							"REGION": "dev",
						},
						Image: "postgres",
						Pull:  "not_present",
					},
				},
				Steps: StepSlice{
					{
						Commands: raw.StringSlice{"echo alpha"},
						Name:     "alpha",
						Image:    "alpine:latest",
						Pull:     "not_present",
						Ruleset: Ruleset{
							If: Rules{
								Event: []string{"push"},
							},
							Matcher:  "filepath",
							Operator: "and",
						},
					},
					{
						Commands: raw.StringSlice{"echo beta"},
						Name:     "beta",
						Image:    "alpine:latest",
						Pull:     "not_present",
						Ruleset: Ruleset{
							If: Rules{
								Event: []string{"push"},
							},
							Matcher:  "filepath",
							Operator: "and",
						},
					},
					{
						Commands: raw.StringSlice{"echo gamma"},
						Name:     "gamma",
						Image:    "alpine:latest",
						Pull:     "not_present",
						Environment: raw.StringSliceMap{
							"REGION": "dev",
						},
						Ruleset: Ruleset{
							If: Rules{
								Event: []string{"push"},
							},
							Matcher:  "filepath",
							Operator: "and",
						},
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := new(Build)

		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("Reading file for UnmarshalYAML returned err: %v", err)
		}

		err = yaml.Unmarshal(b, got)
		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalYAML is %v, want %v", got, test.want)
		}
	}
}
