// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/raw"
)

func TestLibrary_Executor_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		executor *Executor
		want     *Executor
	}{
		{
			executor: testExecutor(),
			want:     testExecutor(),
		},
		{
			executor: new(Executor),
			want:     new(Executor),
		},
	}

	// run tests
	for _, test := range tests {
		if test.executor.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.executor.GetID(), test.want.GetID())
		}

		if test.executor.GetHost() != test.want.GetHost() {
			t.Errorf("GetHost is %v, want %v", test.executor.GetHost(), test.want.GetHost())
		}

		if test.executor.GetRuntime() != test.want.GetRuntime() {
			t.Errorf("GetRuntime is %v, want %v", test.executor.GetRuntime(), test.want.GetRuntime())
		}

		if test.executor.GetDistribution() != test.want.GetDistribution() {
			t.Errorf("GetDistribution is %v, want %v", test.executor.GetDistribution(), test.want.GetDistribution())
		}

		if !reflect.DeepEqual(test.executor.GetBuild(), test.want.GetBuild()) {
			t.Errorf("GetBuild is %v, want %v", test.executor.GetBuild(), test.want.GetBuild())
		}

		if !reflect.DeepEqual(test.executor.GetRepo(), test.want.GetRepo()) {
			t.Errorf("GetRepo is %v, want %v", test.executor.GetRepo(), test.want.GetRepo())
		}

		if !reflect.DeepEqual(test.executor.GetPipeline(), test.want.GetPipeline()) {
			t.Errorf("GetPipeline is %v, want %v", test.executor.GetPipeline(), test.want.GetPipeline())
		}
	}
}

func TestLibrary_Executor_Setters(t *testing.T) {
	// setup types
	var e *Executor

	// setup tests
	tests := []struct {
		executor *Executor
		want     *Executor
	}{
		{
			executor: testExecutor(),
			want:     testExecutor(),
		},
		{
			executor: e,
			want:     new(Executor),
		},
	}

	// run tests
	for _, test := range tests {
		test.executor.SetID(test.want.GetID())
		test.executor.SetHost(test.want.GetHost())
		test.executor.SetRuntime(test.want.GetRuntime())
		test.executor.SetDistribution(test.want.GetDistribution())
		test.executor.SetBuild(test.want.GetBuild())
		test.executor.SetRepo(test.want.GetRepo())
		test.executor.SetPipeline(test.want.GetPipeline())

		if test.executor.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.executor.GetID(), test.want.GetID())
		}

		if test.executor.GetHost() != test.want.GetHost() {
			t.Errorf("SetHost is %v, want %v", test.executor.GetHost(), test.want.GetHost())
		}

		if test.executor.GetRuntime() != test.want.GetRuntime() {
			t.Errorf("SetRuntime is %v, want %v", test.executor.GetRuntime(), test.want.GetRuntime())
		}

		if test.executor.GetDistribution() != test.want.GetDistribution() {
			t.Errorf("SetDistribution is %v, want %v", test.executor.GetDistribution(), test.want.GetDistribution())
		}

		if !reflect.DeepEqual(test.executor.GetBuild(), test.want.GetBuild()) {
			t.Errorf("SetBuild is %v, want %v", test.executor.GetBuild(), test.want.GetBuild())
		}

		if !reflect.DeepEqual(test.executor.GetRepo(), test.want.GetRepo()) {
			t.Errorf("SetRepo is %v, want %v", test.executor.GetRepo(), test.want.GetRepo())
		}

		if !reflect.DeepEqual(test.executor.GetPipeline(), test.want.GetPipeline()) {
			t.Errorf("SetPipeline is %v, want %v", test.executor.GetPipeline(), test.want.GetPipeline())
		}
	}
}

func TestLibrary_Executor_String(t *testing.T) {
	// setup types
	e := testExecutor()

	want := fmt.Sprintf(`{
  Build: %s,
  Distribution: %s,
  Host: %s,
  ID: %d,
  Repo: %v,
  Runtime: %s,
  Pipeline: %v,
}`,
		strings.ReplaceAll(e.Build.String(), " ", "  "),
		e.GetDistribution(),
		e.GetHost(),
		e.GetID(),
		strings.ReplaceAll(e.Repo.String(), " ", "  "),
		e.GetRuntime(),
		e.GetPipeline(),
	)

	// run test
	got := e.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testExecutor is a test helper function to create a Executor
// type with all fields set to a fake value.
func testExecutor() *Executor {
	e := new(Executor)

	e.SetID(1)
	e.SetHost("company.example.com")
	e.SetRuntime("docker")
	e.SetDistribution("linux")
	e.SetBuild(*testBuild())
	e.SetRepo(*testRepo())
	e.SetPipeline(pipeline.Build{
		Version: "1",
		ID:      "github_octocat_1",
		Steps: pipeline.ContainerSlice{
			{
				ID:        "step_github_octocat_1_init",
				Directory: "/home/github/octocat",
				Image:     "#init",
				Name:      "init",
				Number:    1,
				Pull:      "always",
			},
			{
				ID:        "step_github_octocat_1_clone",
				Directory: "/home/github/octocat",
				Image:     "target/vela-git:v0.3.0",
				Name:      "clone",
				Number:    2,
				Pull:      "always",
			},
			{
				ID:        "step_github_octocat_1_echo",
				Commands:  []string{"echo hello"},
				Directory: "/home/github/octocat",
				Image:     "alpine:latest",
				Name:      "echo",
				Number:    3,
				Pull:      "always",
			},
		},
	})

	return e
}

// testBuild is a test helper function to create a Build
// type with all fields set to a fake value.
//
// TODO: remove this function once the Build type is moved to server.
func testBuild() *library.Build {
	b := new(library.Build)

	b.SetID(1)
	b.SetRepoID(1)
	b.SetPipelineID(1)
	b.SetNumber(1)
	b.SetParent(1)
	b.SetEvent("push")
	b.SetStatus("running")
	b.SetError("")
	b.SetEnqueued(1563474077)
	b.SetCreated(1563474076)
	b.SetStarted(1563474078)
	b.SetFinished(1563474079)
	b.SetDeploy("")
	b.SetDeployNumber(0)
	b.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	b.SetClone("https://github.com/github/octocat.git")
	b.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetTitle("push received from https://github.com/github/octocat")
	b.SetMessage("First commit...")
	b.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	b.SetSender("OctoKitty")
	b.SetAuthor("OctoKitty")
	b.SetEmail("OctoKitty@github.com")
	b.SetLink("https://example.company.com/github/octocat/1")
	b.SetBranch("main")
	b.SetRef("refs/heads/main")
	b.SetBaseRef("")
	b.SetHeadRef("changes")
	b.SetHost("example.company.com")
	b.SetRuntime("docker")
	b.SetDistribution("linux")
	b.SetApprovedAt(1563474076)
	b.SetApprovedBy("OctoCat")

	return b
}
