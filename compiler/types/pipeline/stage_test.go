// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestPipeline_StageSlice_Purge(t *testing.T) {
	// setup types
	stages := testStages()
	*stages = (*stages)[:len(*stages)-1]

	// setup tests
	tests := []struct {
		stages *StageSlice
		want   *StageSlice
	}{
		{
			stages: testStages(),
			want:   stages,
		},
		{
			stages: new(StageSlice),
			want:   new(StageSlice),
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

		got, _ := test.stages.Purge(r)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Purge is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_StageSlice_Sanitize(t *testing.T) {
	// setup types
	stages := testStages()
	(*stages)[0].Steps[0].ID = "github-octocat._1_init_init"
	(*stages)[1].Steps[0].ID = "github-octocat._1_clone_clone"
	(*stages)[2].Steps[0].ID = "github-octocat._1_echo_echo"

	kubeStages := testStages()
	(*kubeStages)[0].Steps[0].ID = "github-octocat--1-init-init"
	(*kubeStages)[1].Steps[0].ID = "github-octocat--1-clone-clone"
	(*kubeStages)[2].Steps[0].ID = "github-octocat--1-echo-echo"

	// setup tests
	tests := []struct {
		driver string
		stages *StageSlice
		want   *StageSlice
	}{
		{
			driver: constants.DriverDocker,
			stages: testStages(),
			want:   stages,
		},
		{
			driver: constants.DriverKubernetes,
			stages: testStages(),
			want:   kubeStages,
		},
		{
			driver: constants.DriverDocker,
			stages: new(StageSlice),
			want:   new(StageSlice),
		},
		{
			driver: constants.DriverKubernetes,
			stages: new(StageSlice),
			want:   new(StageSlice),
		},
		{
			driver: "foo",
			stages: new(StageSlice),
			want:   nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.stages.Sanitize(test.driver)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Sanitize is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Stage_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		stage       *Stage
		environment map[string]string
		failure     bool
	}{
		{
			stage: &Stage{
				Name:        "testStage",
				Environment: map[string]string{"FOO": "bar"},
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

func testStages() *StageSlice {
	return &StageSlice{
		{
			Name:        "init",
			Environment: map[string]string{"FOO": "bar"},
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
			Name:        "clone",
			Needs:       []string{"init"},
			Environment: map[string]string{"FOO": "bar"},
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
			Name:        "echo",
			Needs:       []string{"clone"},
			Environment: map[string]string{"FOO": "bar"},
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
	}
}
