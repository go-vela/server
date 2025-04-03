// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestPipeline_Build_Purge(t *testing.T) {
	// setup types
	stages := testBuildStages()
	stages.Stages = stages.Stages[:len(stages.Stages)-1]

	steps := testBuildSteps()
	steps.Steps = steps.Steps[:len(steps.Steps)-1]

	// setup tests
	tests := []struct {
		pipeline *Build
		want     *Build
		wantErr  bool
	}{
		{
			pipeline: testBuildStages(),
			want:     stages,
		},
		{
			pipeline: testBuildSteps(),
			want:     steps,
		},
		{
			pipeline: new(Build),
			want:     new(Build),
		},
		{
			pipeline: &Build{
				Stages: StageSlice{
					{
						Name: "init",
						Steps: ContainerSlice{
							{
								ID:          "github octocat._1_init_init",
								Directory:   "/home/github/octocat",
								Environment: map[string]string{"FOO": "bar"},
								Image:       "#init",
								Name:        "init",
								Number:      1,
								Pull:        "always",
							},
						},
					},
				},
				Steps: ContainerSlice{
					{
						ID:          "step_github octocat._1_init",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			pipeline: &Build{
				Steps: ContainerSlice{
					{
						ID:          "step_github octocat._1_init",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
					{
						ID:          "step_github octocat._1_bad_regexp",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine",
						Name:        "bad_regexp",
						Number:      2,
						Pull:        "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push"}, Branch: []string{"*-dev"}, Matcher: "regexp"},
							Operator: "and",
							Matcher:  "regexp",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	// run tests
	for _, test := range tests {
		r := &RuleData{
			Branch: "main",
			Event:  "pull_request",
			Path:   []string{},
			Repo:   "foo/bar",
			Tag:    "refs/heads/main",
		}

		got, err := test.pipeline.Purge(r)

		if test.wantErr && err == nil {
			t.Errorf("Purge should have returned an error, got: %v", got)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Purge is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Build_Sanitize(t *testing.T) {
	// setup types
	stages := testBuildStages()
	stages.ID = "github-Octocat._1"
	stages.Services[0].ID = "service_github-octocat._1_postgres"
	stages.Stages[0].Steps[0].ID = "github-octocat._1_init_init"
	stages.Stages[1].Steps[0].ID = "github-octocat._1_clone_clone"
	stages.Stages[2].Steps[0].ID = "github-octocat._1_echo_echo"
	stages.Secrets[0].Origin.ID = "secret_github-octocat._1_vault"

	kubeStages := testBuildStages()
	kubeStages.ID = "github-octocat--1"
	kubeStages.Services[0].ID = "service-github-octocat--1-postgres"
	kubeStages.Stages[0].Steps[0].ID = "github-octocat--1-init-init"
	kubeStages.Stages[1].Steps[0].ID = "github-octocat--1-clone-clone"
	kubeStages.Stages[2].Steps[0].ID = "github-octocat--1-echo-echo"
	kubeStages.Secrets[0].Origin.ID = "secret-github-octocat--1-vault"

	steps := testBuildSteps()
	steps.ID = "github-octocat._1"
	steps.Services[0].ID = "service_github-octocat._1_postgres"
	steps.Steps[0].ID = "step_github-octocat._1_init"
	steps.Steps[1].ID = "step_github-octocat._1_clone"
	steps.Steps[2].ID = "step_github-octocat._1_echo"
	steps.Secrets[0].Origin.ID = "secret_github-octocat._1_vault"

	kubeSteps := testBuildSteps()
	kubeSteps.ID = "github-octocat--1"
	kubeSteps.Services[0].ID = "service-github-octocat--1-postgres"
	kubeSteps.Steps[0].ID = "step-github-octocat--1-init"
	kubeSteps.Steps[1].ID = "step-github-octocat--1-clone"
	kubeSteps.Steps[2].ID = "step-github-octocat--1-echo"
	kubeSteps.Secrets[0].Origin.ID = "secret-github-octocat--1-vault"

	// setup tests
	tests := []struct {
		driver   string
		pipeline *Build
		want     *Build
	}{
		{
			driver:   constants.DriverDocker,
			pipeline: testBuildStages(),
			want:     stages,
		},
		{
			driver:   constants.DriverKubernetes,
			pipeline: testBuildStages(),
			want:     kubeStages,
		},
		{
			driver:   constants.DriverDocker,
			pipeline: testBuildSteps(),
			want:     steps,
		},
		{
			driver:   constants.DriverKubernetes,
			pipeline: testBuildSteps(),
			want:     kubeSteps,
		},
		{
			driver:   constants.DriverDocker,
			pipeline: new(Build),
			want:     new(Build),
		},
		{
			driver:   constants.DriverKubernetes,
			pipeline: new(Build),
			want:     new(Build),
		},
		{
			driver: constants.DriverDocker,
			pipeline: &Build{
				Stages: StageSlice{
					{
						Name: "init",
						Steps: ContainerSlice{
							{
								ID:          "github octocat._1_init_init",
								Directory:   "/home/github/octocat",
								Environment: map[string]string{"FOO": "bar"},
								Image:       "#init",
								Name:        "init",
								Number:      1,
								Pull:        "always",
							},
						},
					},
				},
				Steps: ContainerSlice{
					{
						ID:          "step_github octocat._1_init",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.pipeline.Sanitize(test.driver)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Sanitize is %v, want %v", got, test.want)
		}
	}
}

func testBuildStages() *Build {
	return &Build{
		Version:     "1",
		ID:          "github Octocat._1",
		Environment: map[string]string{"HELLO": "Hello, Global Message"},
		Services: ContainerSlice{
			{
				ID:          "service_github octocat._1_postgres",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
			},
		},
		Stages: StageSlice{
			{
				Name: "init",
				Steps: ContainerSlice{
					{
						ID:          "github octocat._1_init_init",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "clone",
				Needs: []string{"init"},
				Steps: ContainerSlice{
					{
						ID:          "github octocat._1_clone_clone",
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-git:v0.3.0",
						Name:        "clone",
						Number:      2,
						Pull:        "always",
					},
				},
			},
			{
				Name:  "echo",
				Needs: []string{"clone"},
				Steps: ContainerSlice{
					{
						ID:          "github octocat._1_echo_echo",
						Commands:    []string{"echo hello"},
						Directory:   "/home/github/octocat",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "alpine:latest",
						Name:        "echo",
						Number:      3,
						Pull:        "always",
						Ruleset: Ruleset{
							If:       Rules{Event: []string{"push"}},
							Operator: "and",
						},
					},
				},
			},
		},
		Secrets: SecretSlice{
			{
				Name: "foobar",
				Origin: &Container{
					ID:          "secret_github octocat._1_vault",
					Directory:   "/home/github/octocat",
					Environment: map[string]string{"FOO": "bar"},
					Image:       "vault:latest",
					Name:        "vault",
					Number:      1,
				},
			},
		},
	}
}

func testBuildSteps() *Build {
	return &Build{
		Version:     "1",
		ID:          "github octocat._1",
		Environment: map[string]string{"HELLO": "Hello, Global Message"},
		Services: ContainerSlice{
			{
				ID:          "service_github octocat._1_postgres",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "postgres:12-alpine",
				Name:        "postgres",
				Number:      1,
				Ports:       []string{"5432:5432"},
			},
		},
		Steps: ContainerSlice{
			{
				ID:          "step_github octocat._1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			{
				ID:          "step_github octocat._1_clone",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			{
				ID:          "step_github octocat._1_echo",
				Commands:    []string{"echo hello"},
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      3,
				Pull:        "always",
				Ruleset: Ruleset{
					If:       Rules{Event: []string{"push"}},
					Operator: "and",
				},
			},
		},
		Secrets: SecretSlice{
			{
				Name: "foobar",
				Origin: &Container{
					ID:          "secret_github octocat._1_vault",
					Directory:   "/home/github/octocat",
					Environment: map[string]string{"FOO": "bar"},
					Image:       "vault:latest",
					Name:        "vault",
					Number:      1,
				},
			},
		},
	}
}
