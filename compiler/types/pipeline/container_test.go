// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestPipeline_ContainerSlice_Purge(t *testing.T) {
	// setup types
	containers := testContainers()
	*containers = (*containers)[:len(*containers)-1]

	// setup tests
	tests := []struct {
		containers *ContainerSlice
		want       *ContainerSlice
	}{
		{
			containers: testContainers(),
			want:       containers,
		},
		{
			containers: new(ContainerSlice),
			want:       new(ContainerSlice),
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

		got, _ := test.containers.Purge(r)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Purge is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_ContainerSlice_Sanitize(t *testing.T) {
	// setup types
	containers := testContainers()
	(*containers)[0].ID = "step_github-octocat._1_init"
	(*containers)[1].ID = "step_github-octocat._1_clone"
	(*containers)[2].ID = "step_github-octocat._1_echo"

	kubeContainers := testContainers()
	(*kubeContainers)[0].ID = "step-github-octocat--1-init"
	(*kubeContainers)[1].ID = "step-github-octocat--1-clone"
	(*kubeContainers)[2].ID = "step-github-octocat--1-echo"

	// setup tests
	tests := []struct {
		driver     string
		containers *ContainerSlice
		want       *ContainerSlice
	}{
		{
			driver:     constants.DriverDocker,
			containers: testContainers(),
			want:       containers,
		},
		{
			driver:     constants.DriverKubernetes,
			containers: testContainers(),
			want:       kubeContainers,
		},
		{
			driver:     constants.DriverDocker,
			containers: new(ContainerSlice),
			want:       new(ContainerSlice),
		},
		{
			driver:     constants.DriverKubernetes,
			containers: new(ContainerSlice),
			want:       new(ContainerSlice),
		},
		{
			driver:     "foo",
			containers: new(ContainerSlice),
			want:       new(ContainerSlice),
		},
	}

	// run tests
	for _, test := range tests {
		got := test.containers.Sanitize(test.driver)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Sanitize is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Container_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		container *Container
		want      bool
	}{
		{
			container: &Container{},
			want:      true,
		},
		{
			container: nil,
			want:      true,
		},
		{
			container: &Container{ID: "foo"},
			want:      false,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.container.Empty()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Empty is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Container_Execute(t *testing.T) {
	// setup types
	containers := testContainers()
	*containers = (*containers)[:len(*containers)-1]

	// setup tests
	tests := []struct {
		container *Container
		ruleData  *RuleData
		want      bool
	}{
		{ // empty/nil container
			container: nil,
			ruleData:  nil,
			want:      false,
		},
		{ // empty container ruleset with build running
			container: &Container{
				Name:     "empty-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
				Path:   []string{"README.md"},
			},
			want: true,
		},
		{ // empty container ruleset with build success
			container: &Container{
				Name:     "empty-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
				Path:   []string{"src/*, *.gradle"},
			},
			want: true,
		},
		{ // empty container ruleset with build failure
			container: &Container{
				Name:     "empty-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
			},
			ruleData: &RuleData{
				Branch:  "main",
				Event:   "push",
				Repo:    "foo/bar",
				Status:  "failure",
				Comment: "LGTM",
			},
			want: false,
		},
		{ // status success container with build running
			container: &Container{
				Name:     "status-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch:  "main",
				Event:   "push",
				Repo:    "foo/bar",
				Status:  "running",
				Comment: "CI Run",
			},
			want: true,
		},
		{ // status success container with build success
			container: &Container{
				Name:     "status-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // status success container with build failure
			container: &Container{
				Name:     "status-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: false,
		},
		{ // status/failure success container with build running
			container: &Container{
				Name:     "status/failure-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess, constants.StatusFailure},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: true,
		},
		{ // status/failure success container with build success
			container: &Container{
				Name:     "status/failure-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess, constants.StatusFailure},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // status/failure success container with build failure
			container: &Container{
				Name:     "status/failure-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Status: []string{constants.StatusSuccess, constants.StatusFailure},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: true,
		},
		{ // no status container with build running
			container: &Container{
				Name:     "branch/event/no-status-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: true,
		},
		{ // no status container with build failure
			container: &Container{
				Name:     "branch/event/no-status-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: false,
		},
		{ // branch/event/path container with build running
			container: &Container{
				Name:     "branch/event/path-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Path:   []string{"README.md"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: true,
		},
		{ // branch/event/path container with build success
			container: &Container{
				Name:     "branch/event/path-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Path:   []string{"README.md"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // branch/event/path container with build failure
			container: &Container{
				Name:     "branch/event/path-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Path:   []string{"README.md"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: false,
		},
		{ // branch/event/comment container with build running
			container: &Container{
				Name:     "branch/event/comment-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch:  []string{"develop"},
						Event:   []string{constants.EventComment},
						Comment: []string{"run vela"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "develop",
				Event:  "comment",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: true,
		},
		{ // branch/event/comment container with build success
			container: &Container{
				Name:     "branch/event/comment-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch:  []string{"develop"},
						Event:   []string{constants.EventComment},
						Comment: []string{"run vela"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "develop",
				Event:  "comment",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // branch/event/comment container with build failure
			container: &Container{
				Name:     "branch/event/comment-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch:  []string{"develop"},
						Event:   []string{constants.EventComment},
						Comment: []string{"run vela"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "develop",
				Event:  "comment",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: false,
		},
		{ // branch/event/status container with build running
			container: &Container{
				Name:     "branch/event/status-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: true,
		},
		{ // branch/event/status container with build success
			container: &Container{
				Name:     "branch/event/status-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // branch/event/status container with build failure
			container: &Container{
				Name:     "branch/event/status-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: false,
		},
		{ // branch/event/status container with or operator with build failure
			container: &Container{
				Name:     "branch/event/status-failure-or",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "or",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: true,
		},
		{ // tag/event/status container with build running
			container: &Container{
				Name:     "tag/event/status-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Tag:    []string{"v*"},
						Event:  []string{constants.EventTag},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "tag",
				Repo:   "foo/bar",
				Status: "running",
				Tag:    "v0.1.0",
			},
			want: true,
		},
		{ // tag/event/status container with build success
			container: &Container{
				Name:     "tag/event/status-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Tag:    []string{"v*"},
						Event:  []string{constants.EventTag},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "tag",
				Repo:   "foo/bar",
				Status: "success",
				Tag:    "v0.1.0",
			},
			want: true,
		},
		{ // tag/event/status container with build failure
			container: &Container{
				Name:     "tag/event/status-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Tag:    []string{"v*"},
						Event:  []string{constants.EventTag},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "tag",
				Repo:   "foo/bar",
				Status: "failure",
				Tag:    "v0.1.0",
			},
			want: false,
		},
		{ // status unless success container with build running
			container: &Container{
				Name:     "unless/status-running",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "running",
			},
			want: false,
		},
		{ // status unless success container with build success
			container: &Container{
				Name:     "unless/status-success",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: false,
		},
		{ // status unless success container with build failure
			container: &Container{
				Name:     "unless/status-failure",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: true,
		},
		{ // status unless success container with build success
			container: &Container{
				Name:     "status unless",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Branch: []string{"main"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "push",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: false,
		},
		{ // status unless success container with build failure
			container: &Container{
				Name:     "status unless",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Branch: []string{"dev"},
						Event:  []string{constants.EventPush},
						Status: []string{constants.StatusSuccess},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "main",
				Event:  "pull_request",
				Repo:   "foo/bar",
				Status: "failure",
			},
			want: true,
		},
		{ // pull request labeled success container with build success
			container: &Container{
				Name:     "pull-request-labeled",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Branch: []string{"fix/1234"},
						Event:  []string{constants.EventPull + constants.ActionLabeled},
						Label:  []string{"enhancement", "documentation"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "fix/1234",
				Event:  "pull_request:labeled",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // pull request unlabeled success container with build success
			container: &Container{
				Name:     "pull-request-unlabeled",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					If: Rules{
						Event: []string{constants.EventPull + constants.ActionUnlabeled},
						Label: []string{"enhancement"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "fix/1234",
				Event:  "pull_request:unlabeled",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // pull request labeled unless ruleset, success container with build success
			container: &Container{
				Name:     "pull-request-labeled",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Branch: []string{"fix/1234"},
						Event:  []string{constants.EventPull + constants.ActionLabeled},
						Label:  []string{"enhancement", "documentation"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "fix/1234",
				Event:  "pull_request:labeled",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
		{ // pull request unlabeled unless ruleset, success container with build success
			container: &Container{
				Name:     "pull-request-unlabeled",
				Image:    "alpine:latest",
				Commands: []string{"echo \"Hey Vela\""},
				Ruleset: Ruleset{
					Unless: Rules{
						Event: []string{constants.EventPull + constants.ActionUnlabeled},
						Label: []string{"enhancement"},
					},
					Operator: "and",
				},
			},
			ruleData: &RuleData{
				Branch: "fix/1234",
				Event:  "pull_request:unlabeled",
				Repo:   "foo/bar",
				Status: "success",
			},
			want: true,
		},
	}

	// run tests
	for _, test := range tests {
		got, _ := test.container.Execute(test.ruleData)

		if got != test.want {
			t.Errorf("Container Execute %s is %v, want %v", test.container.Name, got, test.want)
		}
	}
}

func TestPipeline_Container_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		container   *Container
		environment map[string]string
		failure     bool
	}{
		{
			container: &Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			container:   &Container{},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			container:   nil,
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			container: &Container{
				ID:          "step_github_octocat_1_init",
				Directory:   "/home/github/octocat",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "always",
			},
			environment: nil,
			failure:     true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.container.MergeEnv(test.environment)

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

func TestPipeline_Container_Substitute(t *testing.T) {
	// setup tests
	tests := []struct {
		container *Container
		want      *Container
		failure   bool
	}{
		{
			container: &Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo ${FOO}", "echo $${BAR}"},
				Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "always",
			},
			want: &Container{
				ID:          "step_github_octocat_1_echo",
				Commands:    []string{"echo baz", "echo ${BAR}"},
				Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
				Image:       "alpine:latest",
				Name:        "echo",
				Number:      1,
				Pull:        "always",
			},
			failure: false,
		},
		{
			container: &Container{
				ID:       "step_github_octocat_1_echo",
				Commands: []string{"echo ${FOO}", "echo ${BAR}"},
				Environment: map[string]string{
					"FOO": "1\n2\n",
					"BAR": "`~!@#$%^&*()-_=+[{]}\\|;:',<.>/?",
				},
				Image:  "alpine:latest",
				Name:   "echo",
				Number: 1,
				Pull:   "always",
			},
			want: &Container{
				ID:       "step_github_octocat_1_echo",
				Commands: []string{"echo ${FOO}", "echo ${BAR}"},
				Environment: map[string]string{
					"FOO": "1\n2\n",
					"BAR": "`~!@#$%^&*()-_=+[{]}\\|;:',<.>/?",
				},
				Image:  "alpine:latest",
				Name:   "echo",
				Number: 1,
				Pull:   "always",
			},
			failure: false,
		},
		{
			container: nil,
			want:      nil,
			failure:   true,
		},
		{
			container: new(Container),
			want:      new(Container),
			failure:   true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.container.Substitute()

		if test.failure {
			if err == nil {
				t.Errorf("Substitute should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Substitute returned err: %v", err)
		}

		if !reflect.DeepEqual(test.container, test.want) {
			t.Errorf("Substitute is %v, want %v", test.container, test.want)
		}
	}
}

func testContainers() *ContainerSlice {
	return &ContainerSlice{
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
			IDRequest:   "yes",
		},
		{
			ID:          "step_github/octocat._1_echo",
			Commands:    []string{"echo hello"},
			Directory:   "/home/github/octocat",
			Environment: map[string]string{"FOO": "bar"},
			Image:       "alpine:latest",
			Name:        "echo",
			Number:      3,
			Pull:        "always",
			ReportAs:    "echo-step",
			Ruleset: Ruleset{
				If:       Rules{Event: []string{"push"}},
				Operator: "and",
			},
		},
	}
}
