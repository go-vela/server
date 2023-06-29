// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-vela/server/database/schedule"

	"github.com/go-vela/server/database/repo"

	"github.com/go-vela/server/database/pipeline"

	"github.com/go-vela/server/database/log"

	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/hook"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/database/service"

	"github.com/kr/pretty"

	"github.com/go-vela/types/raw"

	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/types/library"
)

func TestDatabase_Integration(t *testing.T) {
	// check if we should skip the integration test
	//
	// https://konradreiche.com/blog/how-to-separate-integration-tests-in-go
	if os.Getenv("INTEGRATION") == "" {
		t.Skipf("skipping %s integration test due to environment variable constraint", t.Name())
	}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		config  *config
	}{
		{
			name:    "success with postgres",
			failure: false,
			config: &config{
				Driver:           "postgres",
				Address:          "postgres://vela:notARealPassword12345@localhost:5432/vela",
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create resources for testing
			buildOne := new(library.Build)
			buildOne.SetID(1)
			buildOne.SetRepoID(1)
			buildOne.SetPipelineID(1)
			buildOne.SetNumber(1)
			buildOne.SetParent(1)
			buildOne.SetEvent("push")
			buildOne.SetStatus("running")
			buildOne.SetError("")
			buildOne.SetEnqueued(1563474077)
			buildOne.SetCreated(1563474076)
			buildOne.SetStarted(1563474078)
			buildOne.SetFinished(1563474079)
			buildOne.SetDeploy("")
			buildOne.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
			buildOne.SetClone("https://github.com/github/octocat.git")
			buildOne.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135163")
			buildOne.SetTitle("push received from https://github.com/github/octocat")
			buildOne.SetMessage("First commit...")
			buildOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
			buildOne.SetSender("OctoKitty")
			buildOne.SetAuthor("OctoKitty")
			buildOne.SetEmail("OctoKitty@github.com")
			buildOne.SetLink("https://example.company.com/github/octocat/1")
			buildOne.SetBranch("master")
			buildOne.SetRef("refs/heads/master")
			buildOne.SetBaseRef("")
			buildOne.SetHeadRef("changes")
			buildOne.SetHost("example.company.com")
			buildOne.SetRuntime("docker")
			buildOne.SetDistribution("linux")

			buildTwo := new(library.Build)
			buildTwo.SetID(2)
			buildTwo.SetRepoID(1)
			buildTwo.SetPipelineID(1)
			buildTwo.SetNumber(2)
			buildTwo.SetParent(1)
			buildTwo.SetEvent("pull_request")
			buildTwo.SetStatus("running")
			buildTwo.SetError("")
			buildTwo.SetEnqueued(1563474077)
			buildTwo.SetCreated(1563474076)
			buildTwo.SetStarted(1563474078)
			buildTwo.SetFinished(1563474079)
			buildTwo.SetDeploy("")
			buildTwo.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
			buildTwo.SetClone("https://github.com/github/octocat.git")
			buildTwo.SetSource("https://github.com/github/octocat/48afb5bdc41ad69bf22588491333f7cf71135164")
			buildTwo.SetTitle("pull_request received from https://github.com/github/octocat")
			buildTwo.SetMessage("Second commit...")
			buildTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
			buildTwo.SetSender("OctoKitty")
			buildTwo.SetAuthor("OctoKitty")
			buildTwo.SetEmail("OctoKitty@github.com")
			buildTwo.SetLink("https://example.company.com/github/octocat/2")
			buildTwo.SetBranch("master")
			buildTwo.SetRef("refs/heads/master")
			buildTwo.SetBaseRef("")
			buildTwo.SetHeadRef("changes")
			buildTwo.SetHost("example.company.com")
			buildTwo.SetRuntime("docker")
			buildTwo.SetDistribution("linux")

			repoOne := new(library.Repo)
			repoOne.SetID(1)
			repoOne.SetOrg("github")
			repoOne.SetName("octocat")
			repoOne.SetFullName("github/octocat")
			repoOne.SetLink("https://github.com/github/octocat")
			repoOne.SetClone("https://github.com/github/octocat.git")
			repoOne.SetBranch("main")
			repoOne.SetTopics([]string{"cloud", "security"})
			repoOne.SetBuildLimit(10)
			repoOne.SetTimeout(30)
			repoOne.SetCounter(0)
			repoOne.SetVisibility("public")
			repoOne.SetPrivate(false)
			repoOne.SetTrusted(false)
			repoOne.SetActive(true)
			repoOne.SetAllowPull(false)
			repoOne.SetAllowPush(true)
			repoOne.SetAllowDeploy(false)
			repoOne.SetAllowTag(false)
			repoOne.SetAllowComment(false)
			repoOne.SetPipelineType("")
			repoOne.SetPreviousName("")

			repoTwo := new(library.Repo)
			repoTwo.SetID(2)
			repoTwo.SetOrg("github")
			repoTwo.SetName("octokitty")
			repoTwo.SetFullName("github/octokitty")
			repoTwo.SetLink("https://github.com/github/octokitty")
			repoTwo.SetClone("https://github.com/github/octokitty.git")
			repoTwo.SetBranch("master")
			repoTwo.SetTopics([]string{"cloud", "security"})
			repoTwo.SetBuildLimit(10)
			repoTwo.SetTimeout(30)
			repoTwo.SetCounter(0)
			repoTwo.SetVisibility("public")
			repoTwo.SetPrivate(false)
			repoTwo.SetTrusted(false)
			repoTwo.SetActive(true)
			repoTwo.SetAllowPull(false)
			repoTwo.SetAllowPush(true)
			repoTwo.SetAllowDeploy(false)
			repoTwo.SetAllowTag(false)
			repoTwo.SetAllowComment(false)
			repoTwo.SetPipelineType("")
			repoTwo.SetPreviousName("")

			serviceOne := new(library.Service)
			serviceOne.SetID(1)
			serviceOne.SetBuildID(1)
			serviceOne.SetRepoID(1)
			serviceOne.SetNumber(1)
			serviceOne.SetName("init")
			serviceOne.SetImage("#init")
			serviceOne.SetStatus("running")
			serviceOne.SetError("")
			serviceOne.SetExitCode(0)
			serviceOne.SetCreated(1563474076)
			serviceOne.SetStarted(1563474078)
			serviceOne.SetFinished(1563474079)
			serviceOne.SetHost("example.company.com")
			serviceOne.SetRuntime("docker")
			serviceOne.SetDistribution("linux")

			serviceTwo := new(library.Service)
			serviceTwo.SetID(2)
			serviceTwo.SetBuildID(1)
			serviceTwo.SetRepoID(1)
			serviceTwo.SetNumber(2)
			serviceTwo.SetName("clone")
			serviceTwo.SetImage("target/vela-git:v0.3.0")
			serviceTwo.SetStatus("pending")
			serviceTwo.SetError("")
			serviceTwo.SetExitCode(0)
			serviceTwo.SetCreated(1563474086)
			serviceTwo.SetStarted(1563474088)
			serviceTwo.SetFinished(1563474089)
			serviceTwo.SetHost("example.company.com")
			serviceTwo.SetRuntime("docker")
			serviceTwo.SetDistribution("linux")

			stepOne := new(library.Step)
			stepOne.SetID(1)
			stepOne.SetBuildID(1)
			stepOne.SetRepoID(1)
			stepOne.SetNumber(1)
			stepOne.SetName("init")
			stepOne.SetImage("#init")
			stepOne.SetStage("init")
			stepOne.SetStatus("running")
			stepOne.SetError("")
			stepOne.SetExitCode(0)
			stepOne.SetCreated(1563474076)
			stepOne.SetStarted(1563474078)
			stepOne.SetFinished(1563474079)
			stepOne.SetHost("example.company.com")
			stepOne.SetRuntime("docker")
			stepOne.SetDistribution("linux")

			stepTwo := new(library.Step)
			stepTwo.SetID(2)
			stepTwo.SetBuildID(1)
			stepTwo.SetRepoID(1)
			stepTwo.SetNumber(2)
			stepTwo.SetName("clone")
			stepTwo.SetImage("target/vela-git:v0.3.0")
			stepTwo.SetStage("init")
			stepTwo.SetStatus("pending")
			stepTwo.SetError("")
			stepTwo.SetExitCode(0)
			stepTwo.SetCreated(1563474086)
			stepTwo.SetStarted(1563474088)
			stepTwo.SetFinished(1563474089)
			stepTwo.SetHost("example.company.com")
			stepTwo.SetRuntime("docker")
			stepTwo.SetDistribution("linux")

			db, err := New(
				WithAddress(test.config.Address),
				WithCompressionLevel(test.config.CompressionLevel),
				WithConnectionLife(test.config.ConnectionLife),
				WithConnectionIdle(test.config.ConnectionIdle),
				WithConnectionOpen(test.config.ConnectionOpen),
				WithDriver(test.config.Driver),
				WithEncryptionKey(test.config.EncryptionKey),
				WithSkipCreation(test.config.SkipCreation),
			)
			if err != nil {
				t.Errorf("unable to create new database engine for %s: %v", test.name, err)
			}

			driver := db.Driver()
			if !strings.EqualFold(driver, test.config.Driver) {
				t.Errorf("Driver() is %v, want %v", driver, test.config.Driver)
			}

			err = db.Ping()
			if err != nil {
				t.Errorf("unable to ping database engine for %s: %v", test.name, err)
			}

			t.Run("test_builds", func(t *testing.T) {
				testBuilds(t, db, []*library.Build{buildOne, buildTwo}, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_hooks", func(t *testing.T) {
				testHooks(t, db, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_logs", func(t *testing.T) {
				testLogs(t, db, []*library.Service{serviceOne, serviceTwo}, []*library.Step{stepOne, stepTwo})
			})

			t.Run("test_pipelines", func(t *testing.T) {
				testPipelines(t, db, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_repos", func(t *testing.T) {
				testRepos(t, db, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_schedules", func(t *testing.T) {
				testSchedules(t, db, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_services", func(t *testing.T) {
				testServices(t, db, []*library.Build{buildOne, buildTwo}, []*library.Service{serviceOne, serviceTwo})
			})

			t.Run("test_steps", func(t *testing.T) {
				testSteps(t, db, []*library.Build{buildOne, buildTwo}, []*library.Step{stepOne, stepTwo})
			})

			t.Run("test_users", func(t *testing.T) {
				testUsers(t, db)
			})

			t.Run("test_workers", func(t *testing.T) {
				testWorkers(t, db)
			})

			err = db.Close()
			if err != nil {
				t.Errorf("unable to close database engine for %s: %v", test.name, err)
			}
		})
	}
}

func testBuilds(t *testing.T, db Interface, builds []*library.Build, repos []*library.Repo) {
	// used to track the number of methods we call for builds
	//
	// we start at 2 for creating the table and indexes for builds
	// since those are already called when the database engine starts
	counter := 2

	// create the builds
	for _, build := range builds {
		_, err := db.CreateBuild(build)
		if err != nil {
			t.Errorf("unable to create build %d: %v", build.GetID(), err)
		}
	}
	counter++

	// count the builds
	count, err := db.CountBuilds()
	if err != nil {
		t.Errorf("unable to count builds: %v", err)
	}
	if int(count) != len(builds) {
		t.Errorf("CountBuilds() is %v, want %v", count, len(builds))
	}
	counter++

	// count the builds for a repo
	count, err = db.CountBuildsForRepo(repos[0], nil)
	if err != nil {
		t.Errorf("unable to count builds for repo %s: %v", repos[0].GetFullName(), err)
	}
	if int(count) != len(builds) {
		t.Errorf("CountBuildsForRepo() is %v, want %v", count, len(builds))
	}
	counter++

	// list the builds
	list, err := db.ListBuilds()
	if err != nil {
		t.Errorf("unable to list builds: %v", err)
	}
	if !reflect.DeepEqual(list, builds) {
		pretty.Ldiff(t, list, builds)
		t.Errorf("ListBuilds() is %v, want %v", list, builds)
	}
	counter++

	// list the builds for a repo
	list, count, err = db.ListBuildsForRepo(repos[0], nil, time.Now().UTC().Unix(), 0, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for repo %s: %v", repos[0].GetFullName(), err)
	}
	if !reflect.DeepEqual(list, builds) {
		pretty.Ldiff(t, list, builds)
		t.Errorf("ListBuildsForRepo() is %v, want %v", list, builds)
	}
	counter++

	// lookup the last build by repo
	got, err := db.LastBuildForRepo(repos[0], "main")
	if err != nil {
		t.Errorf("unable to get last build for repo %s: %v", repos[0].GetFullName(), err)
	}
	if !reflect.DeepEqual(got, builds[1]) {
		t.Errorf("GetBuildForRepo() is %v, want %v", got, builds[1])
	}
	counter++

	// lookup the builds by repo and number
	for _, build := range builds {
		got, err = db.GetBuildForRepo(repos[0], build.GetNumber())
		if err != nil {
			t.Errorf("unable to get build %d for repo %s: %v", build.GetID(), repos[0].GetFullName(), err)
		}
		if !reflect.DeepEqual(got, build) {
			t.Errorf("GetBuildForRepo() is %v, want %v", got, build)
		}
	}
	counter++

	// update the builds
	for _, build := range builds {
		build.SetStatus("success")
		_, err = db.UpdateBuild(build)
		if err != nil {
			t.Errorf("unable to update build %d: %v", build.GetID(), err)
		}

		// lookup the build by ID
		got, err = db.GetBuild(build.GetID())
		if err != nil {
			t.Errorf("unable to get build %d by ID: %v", build.GetID(), err)
		}
		if !reflect.DeepEqual(got, build) {
			t.Errorf("GetBuild() is %v, want %v", got, build)
		}
	}
	counter++
	counter++

	// delete the builds
	for _, build := range builds {
		err = db.DeleteBuild(build)
		if err != nil {
			t.Errorf("unable to delete build %d: %v", build.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(build.BuildInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testHooks(t *testing.T, db Interface, repos []*library.Repo) {
	// used to track the number of methods we call for hooks
	//
	// we start at 2 for creating the table and indexes for hooks
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Hook)
	one.SetID(1)

	two := new(library.Hook)
	two.SetID(2)

	hooks := []*library.Hook{one, two}

	// create the hooks
	for _, hook := range hooks {
		_, err := db.CreateHook(hook)
		if err != nil {
			t.Errorf("unable to create hook %s: %v", hook.GetSourceID(), err)
		}
	}
	counter++

	// count the hooks
	count, err := db.CountHooks()
	if err != nil {
		t.Errorf("unable to count hooks: %v", err)
	}
	if int(count) != len(hooks) {
		t.Errorf("CountHooks() is %v, want %v", count, len(hooks))
	}
	counter++

	// list the hooks
	list, err := db.ListHooks()
	if err != nil {
		t.Errorf("unable to list hooks: %v", err)
	}
	if !reflect.DeepEqual(list, hooks) {
		t.Errorf("ListHooks() is %v, want %v", list, hooks)
	}
	counter++

	// lookup the hooks by name
	for _, hook := range hooks {
		got, err := db.GetHookForRepo(repos[0], hook.GetNumber())
		if err != nil {
			t.Errorf("unable to get hook %s for repo %s: %v", hook.GetSourceID(), repos[0].GetFullName(), err)
		}
		if !reflect.DeepEqual(got, hook) {
			t.Errorf("GetHookForRepo() is %v, want %v", got, hook)
		}
	}
	counter++

	// update the hooks
	for _, hook := range hooks {
		hook.SetStatus("success")
		_, err = db.UpdateHook(hook)
		if err != nil {
			t.Errorf("unable to update hook %s: %v", hook.GetSourceID(), err)
		}

		// lookup the hook by ID
		got, err := db.GetHook(hook.GetID())
		if err != nil {
			t.Errorf("unable to get hook %s by ID: %v", hook.GetSourceID(), err)
		}
		if !reflect.DeepEqual(got, hook) {
			t.Errorf("GetHook() is %v, want %v", got, hook)
		}
	}
	counter++
	counter++

	// delete the hooks
	for _, hook := range hooks {
		err = db.DeleteHook(hook)
		if err != nil {
			t.Errorf("unable to delete hook %s: %v", hook.GetSourceID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(hook.HookInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testLogs(t *testing.T, db Interface, services []*library.Service, steps []*library.Step) {
	// used to track the number of methods we call for logs
	//
	// we start at 2 for creating the table and indexes for logs
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Log)
	one.SetID(1)

	two := new(library.Log)
	two.SetID(2)

	logs := []*library.Log{one, two}

	// create the logs
	for _, log := range logs {
		err := db.CreateLog(log)
		if err != nil {
			t.Errorf("unable to create log %d: %v", log.GetID(), err)
		}
	}
	counter++

	// count the logs
	count, err := db.CountLogs()
	if err != nil {
		t.Errorf("unable to count logs: %v", err)
	}
	if int(count) != len(logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(logs))
	}
	counter++

	// list the logs
	list, err := db.ListLogs()
	if err != nil {
		t.Errorf("unable to list logs: %v", err)
	}
	if !reflect.DeepEqual(list, logs) {
		t.Errorf("ListLogs() is %v, want %v", list, logs)
	}
	counter++

	// lookup the logs by service
	for _, log := range logs {
		got, err := db.GetLogForService(services[0])
		if err != nil {
			t.Errorf("unable to get log %d for service %d: %v", log.GetID(), services[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, log) {
			t.Errorf("GetLogForService() is %v, want %v", got, log)
		}
	}
	counter++

	// lookup the logs by service
	for _, log := range logs {
		got, err := db.GetLogForStep(steps[0])
		if err != nil {
			t.Errorf("unable to get log %d for step %d: %v", log.GetID(), steps[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, log) {
			t.Errorf("GetLogForStep() is %v, want %v", got, log)
		}
	}
	counter++

	// update the logs
	for _, log := range logs {
		log.SetData([]byte("bar"))
		err = db.UpdateLog(log)
		if err != nil {
			t.Errorf("unable to update log %d: %v", log.GetID(), err)
		}

		// lookup the log by ID
		got, err := db.GetLog(log.GetID())
		if err != nil {
			t.Errorf("unable to get log %d by ID: %v", log.GetID(), err)
		}
		if !reflect.DeepEqual(got, log) {
			t.Errorf("GetLog() is %v, want %v", got, log)
		}
	}
	counter++
	counter++

	// delete the logs
	for _, log := range logs {
		err = db.DeleteLog(log)
		if err != nil {
			t.Errorf("unable to delete log %d: %v", log.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(log.LogInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testPipelines(t *testing.T, db Interface, repos []*library.Repo) {
	// used to track the number of methods we call for pipelines
	//
	// we start at 2 for creating the table and indexes for pipelines
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Pipeline)
	one.SetID(1)

	two := new(library.Pipeline)
	two.SetID(2)

	pipelines := []*library.Pipeline{one, two}

	// create the pipelines
	for _, pipeline := range pipelines {
		_, err := db.CreatePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to create pipeline %s: %v", pipeline.GetCommit(), err)
		}
	}
	counter++

	// count the pipelines
	count, err := db.CountPipelines()
	if err != nil {
		t.Errorf("unable to count pipelines: %v", err)
	}
	if int(count) != len(pipelines) {
		t.Errorf("CountPipelines() is %v, want %v", count, len(pipelines))
	}
	counter++

	// list the pipelines
	list, err := db.ListPipelines()
	if err != nil {
		t.Errorf("unable to list pipelines: %v", err)
	}
	if !reflect.DeepEqual(list, pipelines) {
		t.Errorf("ListPipelines() is %v, want %v", list, pipelines)
	}
	counter++

	// lookup the pipelines by name
	for _, pipeline := range pipelines {
		got, err := db.GetPipelineForRepo(pipeline.GetCommit(), repos[0])
		if err != nil {
			t.Errorf("unable to get pipeline %s for repo %s: %v", pipeline.GetCommit(), repos[0].GetFullName(), err)
		}
		if !reflect.DeepEqual(got, pipeline) {
			t.Errorf("GetPipelineForRepo() is %v, want %v", got, pipeline)
		}
	}
	counter++

	// update the pipelines
	for _, pipeline := range pipelines {
		pipeline.SetVersion("2")
		_, err = db.UpdatePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to update pipeline %s: %v", pipeline.GetCommit(), err)
		}

		// lookup the pipeline by ID
		got, err := db.GetPipeline(pipeline.GetID())
		if err != nil {
			t.Errorf("unable to get pipeline %s by ID: %v", pipeline.GetCommit(), err)
		}
		if !reflect.DeepEqual(got, pipeline) {
			t.Errorf("GetPipeline() is %v, want %v", got, pipeline)
		}
	}
	counter++
	counter++

	// delete the pipelines
	for _, pipeline := range pipelines {
		err = db.DeletePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to delete pipeline %s: %v", pipeline.GetCommit(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(pipeline.PipelineInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testRepos(t *testing.T, db Interface, repos []*library.Repo) {
	// used to track the number of methods we call for repos
	//
	// we start at 2 for creating the table and indexes for repos
	// since those are already called when the database engine starts
	counter := 2

	// create the repos
	for _, repo := range repos {
		err := db.CreateRepo(repo)
		if err != nil {
			t.Errorf("unable to create repo %s: %v", repo.GetFullName(), err)
		}
	}
	counter++

	// count the repos
	count, err := db.CountRepos()
	if err != nil {
		t.Errorf("unable to count repos: %v", err)
	}
	if int(count) != len(repos) {
		t.Errorf("CountRepos() is %v, want %v", count, len(repos))
	}
	counter++

	// list the repos
	list, err := db.ListRepos()
	if err != nil {
		t.Errorf("unable to list repos: %v", err)
	}
	if !reflect.DeepEqual(list, repos) {
		t.Errorf("ListRepos() is %v, want %v", list, repos)
	}
	counter++

	// lookup the repos by name
	for _, repo := range repos {
		got, err := db.GetRepoForOrg(repo.GetOrg(), repo.GetName())
		if err != nil {
			t.Errorf("unable to get repo %s by org: %v", repo.GetFullName(), err)
		}
		if !reflect.DeepEqual(got, repo) {
			t.Errorf("GetRepoForOrg() is %v, want %v", got, repo)
		}
	}
	counter++

	// update the repos
	for _, repo := range repos {
		repo.SetActive(false)
		err = db.UpdateRepo(repo)
		if err != nil {
			t.Errorf("unable to update repo %s: %v", repo.GetFullName(), err)
		}

		// lookup the repo by ID
		got, err := db.GetRepo(repo.GetID())
		if err != nil {
			t.Errorf("unable to get repo %s by ID: %v", repo.GetFullName(), err)
		}
		if !reflect.DeepEqual(got, repo) {
			t.Errorf("GetRepo() is %v, want %v", got, repo)
		}
	}
	counter++
	counter++

	// delete the repos
	for _, repo := range repos {
		err = db.DeleteRepo(repo)
		if err != nil {
			t.Errorf("unable to delete repo %s: %v", repo.GetFullName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(repo.RepoInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testSchedules(t *testing.T, db Interface, repos []*library.Repo) {
	// used to track the number of methods we call for schedules
	//
	// we start at 2 for creating the table and indexes for schedules
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Schedule)
	one.SetID(1)

	two := new(library.Schedule)
	two.SetID(2)

	schedules := []*library.Schedule{one, two}

	// create the schedules
	for _, schedule := range schedules {
		err := db.CreateSchedule(schedule)
		if err != nil {
			t.Errorf("unable to create schedule %s: %v", schedule.GetName(), err)
		}
	}
	counter++

	// count the schedules
	count, err := db.CountSchedules()
	if err != nil {
		t.Errorf("unable to count schedules: %v", err)
	}
	if int(count) != len(schedules) {
		t.Errorf("CountSchedules() is %v, want %v", count, len(schedules))
	}
	counter++

	// list the schedules
	list, err := db.ListSchedules()
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !reflect.DeepEqual(list, schedules) {
		t.Errorf("ListSchedules() is %v, want %v", list, schedules)
	}
	counter++

	// lookup the schedules by name
	for _, schedule := range schedules {
		got, err := db.GetScheduleForRepo(repos[0], schedule.GetName())
		if err != nil {
			t.Errorf("unable to get schedule %s for repo %s: %v", schedule.GetName(), repos[0].GetFullName(), err)
		}
		if !reflect.DeepEqual(got, schedule) {
			t.Errorf("GetScheduleForRepo() is %v, want %v", got, schedule)
		}
	}
	counter++

	// update the schedules
	for _, schedule := range schedules {
		schedule.SetActive(false)
		err = db.UpdateSchedule(schedule, false)
		if err != nil {
			t.Errorf("unable to update schedule %s: %v", schedule.GetName(), err)
		}

		// lookup the schedule by ID
		got, err := db.GetSchedule(schedule.GetID())
		if err != nil {
			t.Errorf("unable to get schedule %s by ID: %v", schedule.GetName(), err)
		}
		if !reflect.DeepEqual(got, schedule) {
			t.Errorf("GetSchedule() is %v, want %v", got, schedule)
		}
	}
	counter++
	counter++

	// delete the schedules
	for _, schedule := range schedules {
		err = db.DeleteSchedule(schedule)
		if err != nil {
			t.Errorf("unable to delete schedule %s: %v", schedule.GetName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(schedule.ScheduleInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testServices(t *testing.T, db Interface, builds []*library.Build, services []*library.Service) {
	// used to track the number of methods we call for services
	//
	// we start at 2 for creating the table and indexes for services
	// since those are already called when the database engine starts
	counter := 2

	// create the services
	for _, service := range services {
		err := db.CreateService(service)
		if err != nil {
			t.Errorf("unable to create service %s: %v", service.GetName(), err)
		}
	}
	counter++

	// count the services
	count, err := db.CountServices()
	if err != nil {
		t.Errorf("unable to count services: %v", err)
	}
	if int(count) != len(services) {
		t.Errorf("CountServices() is %v, want %v", count, len(services))
	}
	counter++

	// count the services for a build
	count, err = db.CountServicesForBuild(builds[0], nil)
	if err != nil {
		t.Errorf("unable to count services for build %d: %v", builds[0].GetID(), err)
	}
	if int(count) != len(services) {
		t.Errorf("CountServicesForBuild() is %v, want %v", count, len(services))
	}
	counter++

	// list the services
	list, err := db.ListServices()
	if err != nil {
		t.Errorf("unable to list services: %v", err)
	}
	if !reflect.DeepEqual(list, services) {
		t.Errorf("ListServices() is %v, want %v", list, services)
	}
	counter++

	// list the services for a build
	list, count, err = db.ListServicesForBuild(builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list services for build %d: %v", builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Service{services[1], services[0]}) {
		t.Errorf("ListServicesForBuild() is %v, want %v", list, []*library.Service{services[1], services[0]})
	}
	if int(count) != len(services) {
		t.Errorf("ListServicesForBuild() is %v, want %v", count, len(services))
	}
	counter++

	expected := map[string]float64{
		"#init":                  1,
		"target/vela-git:v0.3.0": 1,
	}
	images, err := db.ListServiceImageCount()
	if err != nil {
		t.Errorf("unable to list service image count: %v", err)
	}
	if !reflect.DeepEqual(images, expected) {
		t.Errorf("ListServiceImageCount() is %v, want %v", images, expected)
	}
	counter++

	expected = map[string]float64{
		"pending": 1,
		"failure": 0,
		"killed":  0,
		"running": 1,
		"success": 0,
	}
	statuses, err := db.ListServiceStatusCount()
	if err != nil {
		t.Errorf("unable to list service status count: %v", err)
	}
	if !reflect.DeepEqual(statuses, expected) {
		t.Errorf("ListServiceStatusCount() is %v, want %v", statuses, expected)
	}
	counter++

	// lookup the services by name
	for _, service := range services {
		got, err := db.GetServiceForBuild(builds[0], service.GetNumber())
		if err != nil {
			t.Errorf("unable to get service %s for build %d: %v", service.GetName(), builds[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, service) {
			t.Errorf("GetServiceForBuild() is %v, want %v", got, service)
		}
	}
	counter++

	// update the services
	for _, service := range services {
		service.SetStatus("success")
		err = db.UpdateService(service)
		if err != nil {
			t.Errorf("unable to update service %s: %v", service.GetName(), err)
		}

		// lookup the service by ID
		got, err := db.GetService(service.GetID())
		if err != nil {
			t.Errorf("unable to get service %s by ID: %v", service.GetName(), err)
		}
		if !reflect.DeepEqual(got, service) {
			t.Errorf("GetService() is %v, want %v", got, service)
		}
	}
	counter++
	counter++

	// delete the services
	for _, service := range services {
		err = db.DeleteService(service)
		if err != nil {
			t.Errorf("unable to delete service %s: %v", service.GetName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(service.ServiceInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testSteps(t *testing.T, db Interface, builds []*library.Build, steps []*library.Step) {
	// used to track the number of methods we call for steps
	//
	// we start at 2 for creating the table and indexes for steps
	// since those are already called when the database engine starts
	counter := 2

	// create the steps
	for _, step := range steps {
		err := db.CreateStep(step)
		if err != nil {
			t.Errorf("unable to create step %s: %v", step.GetName(), err)
		}
	}
	counter++

	// count the steps
	count, err := db.CountSteps()
	if err != nil {
		t.Errorf("unable to count steps: %v", err)
	}
	if int(count) != len(steps) {
		t.Errorf("CountSteps() is %v, want %v", count, len(steps))
	}
	counter++

	// count the steps for a build
	count, err = db.CountStepsForBuild(builds[0], nil)
	if err != nil {
		t.Errorf("unable to count steps for build %d: %v", builds[0].GetID(), err)
	}
	if int(count) != len(steps) {
		t.Errorf("CountStepsForBuild() is %v, want %v", count, len(steps))
	}
	counter++

	// list the steps
	list, err := db.ListSteps()
	if err != nil {
		t.Errorf("unable to list steps: %v", err)
	}
	if !reflect.DeepEqual(list, steps) {
		t.Errorf("ListSteps() is %v, want %v", list, steps)
	}
	counter++

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Step{steps[1], steps[0]}) {
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*library.Step{steps[1], steps[0]})
	}
	if int(count) != len(steps) {
		t.Errorf("ListStepsForBuild() is %v, want %v", count, len(steps))
	}
	counter++

	expected := map[string]float64{
		"#init":                  1,
		"target/vela-git:v0.3.0": 1,
	}
	images, err := db.ListStepImageCount()
	if err != nil {
		t.Errorf("unable to list step image count: %v", err)
	}
	if !reflect.DeepEqual(images, expected) {
		t.Errorf("ListStepImageCount() is %v, want %v", images, expected)
	}
	counter++

	expected = map[string]float64{
		"pending": 1,
		"failure": 0,
		"killed":  0,
		"running": 1,
		"success": 0,
	}
	statuses, err := db.ListStepStatusCount()
	if err != nil {
		t.Errorf("unable to list step status count: %v", err)
	}
	if !reflect.DeepEqual(statuses, expected) {
		t.Errorf("ListStepStatusCount() is %v, want %v", statuses, expected)
	}
	counter++

	// lookup the steps by name
	for _, step := range steps {
		got, err := db.GetStepForBuild(builds[0], step.GetNumber())
		if err != nil {
			t.Errorf("unable to get step %s for build %d: %v", step.GetName(), builds[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, step) {
			t.Errorf("GetStepForBuild() is %v, want %v", got, step)
		}
	}
	counter++

	// update the steps
	for _, step := range steps {
		step.SetStatus("success")
		err = db.UpdateStep(step)
		if err != nil {
			t.Errorf("unable to update step %s: %v", step.GetName(), err)
		}

		// lookup the step by ID
		got, err := db.GetStep(step.GetID())
		if err != nil {
			t.Errorf("unable to get step %s by ID: %v", step.GetName(), err)
		}
		if !reflect.DeepEqual(got, step) {
			t.Errorf("GetStep() is %v, want %v", got, step)
		}
	}
	counter++
	counter++

	// delete the steps
	for _, step := range steps {
		err = db.DeleteStep(step)
		if err != nil {
			t.Errorf("unable to delete step %s: %v", step.GetName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(step.StepInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testUsers(t *testing.T, db Interface) {
	// used to track the number of methods we call for users
	//
	// we start at 2 for creating the table and indexes for users
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.User)
	one.SetID(1)
	one.SetName("octocat")
	one.SetToken("superSecretToken")
	one.SetRefreshToken("superSecretRefreshToken")
	one.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	one.SetFavorites([]string{"github/octocat"})
	one.SetActive(true)
	one.SetAdmin(false)

	two := new(library.User)
	two.SetID(2)
	two.SetName("octokitty")
	two.SetToken("superSecretToken")
	two.SetRefreshToken("superSecretRefreshToken")
	two.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	two.SetFavorites([]string{"github/octocat"})
	two.SetActive(true)
	two.SetAdmin(false)

	users := []*library.User{one, two}

	liteOne := new(library.User)
	liteOne.SetID(1)
	liteOne.SetName("octocat")
	liteOne.SetToken("")
	liteOne.SetRefreshToken("")
	liteOne.SetHash("")
	liteOne.SetFavorites([]string{})
	liteOne.SetActive(false)
	liteOne.SetAdmin(false)

	liteTwo := new(library.User)
	liteTwo.SetID(2)
	liteTwo.SetName("octokitty")
	liteTwo.SetToken("")
	liteTwo.SetRefreshToken("")
	liteTwo.SetHash("")
	liteTwo.SetFavorites([]string{})
	liteTwo.SetActive(false)
	liteTwo.SetAdmin(false)

	liteUsers := []*library.User{liteOne, liteTwo}

	// create the users
	for _, user := range users {
		err := db.CreateUser(user)
		if err != nil {
			t.Errorf("unable to create user %s: %v", user.GetName(), err)
		}
	}
	counter++

	// count the users
	count, err := db.CountUsers()
	if err != nil {
		t.Errorf("unable to count users: %v", err)
	}
	if int(count) != len(users) {
		t.Errorf("CountUsers() is %v, want %v", count, len(users))
	}
	counter++

	// list the users
	list, err := db.ListUsers()
	if err != nil {
		t.Errorf("unable to list users: %v", err)
	}
	if !reflect.DeepEqual(list, users) {
		t.Errorf("ListUsers() is %v, want %v", list, users)
	}
	counter++

	// lite list the users
	list, count, err = db.ListLiteUsers(1, 10)
	if err != nil {
		t.Errorf("unable to list lite users: %v", err)
	}
	if !reflect.DeepEqual(list, liteUsers) {
		pretty.Ldiff(t, list, liteUsers)
		if diff := cmp.Diff(list, liteUsers); diff != "" {
			t.Errorf("ListLiteUsers() mismatch (-want +got):\n%s", diff)
		}
		t.Errorf("ListLiteUsers() is %v, want %v", list, liteUsers)
	}
	if int(count) != len(users) {
		t.Errorf("ListLiteUsers() is %v, want %v", count, len(users))
	}
	counter++

	// lookup the users by name
	for _, user := range users {
		got, err := db.GetUserForName(user.GetName())
		if err != nil {
			t.Errorf("unable to get user %s by name: %v", user.GetName(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUserForName() is %v, want %v", got, user)
		}
	}
	counter++

	// update the users
	for _, user := range users {
		user.SetActive(false)
		err = db.UpdateUser(user)
		if err != nil {
			t.Errorf("unable to update user %s: %v", user.GetName(), err)
		}

		// lookup the user by ID
		got, err := db.GetUser(user.GetID())
		if err != nil {
			t.Errorf("unable to get user %s by ID: %v", user.GetName(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUser() is %v, want %v", got, user)
		}
	}
	counter++
	counter++

	// delete the users
	for _, user := range users {
		err = db.DeleteUser(user)
		if err != nil {
			t.Errorf("unable to delete user %s: %v", user.GetName(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(user.UserInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testWorkers(t *testing.T, db Interface) {
	// used to track the number of methods we call for workers
	//
	// we start at 2 for creating the table and indexes for users
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Worker)
	one.SetID(1)
	one.SetHostname("worker-1.example.com")
	one.SetAddress("https://worker-1.example.com")
	one.SetRoutes([]string{"vela"})
	one.SetActive(true)
	one.SetStatus("available")
	one.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	one.SetRunningBuildIDs([]string{"12345"})
	one.SetLastBuildStartedAt(time.Now().UTC().Unix())
	one.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	one.SetLastCheckedIn(time.Now().UTC().Unix())
	one.SetBuildLimit(1)

	two := new(library.Worker)
	two.SetID(2)
	two.SetHostname("worker-2.example.com")
	two.SetAddress("https://worker-2.example.com")
	two.SetRoutes([]string{"vela"})
	two.SetActive(true)
	two.SetStatus("available")
	two.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	two.SetRunningBuildIDs([]string{"12345"})
	two.SetLastBuildStartedAt(time.Now().UTC().Unix())
	two.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	two.SetLastCheckedIn(time.Now().UTC().Unix())
	two.SetBuildLimit(1)

	workers := []*library.Worker{one, two}

	// create the workers
	for _, worker := range workers {
		err := db.CreateWorker(worker)
		if err != nil {
			t.Errorf("unable to create worker %s: %v", worker.GetHostname(), err)
		}
	}
	counter++

	// count the workers
	count, err := db.CountWorkers()
	if err != nil {
		t.Errorf("unable to count workers: %v", err)
	}
	if int(count) != len(workers) {
		t.Errorf("CountWorkers() is %v, want %v", count, len(workers))
	}
	counter++

	// list the workers
	list, err := db.ListWorkers()
	if err != nil {
		t.Errorf("unable to list workers: %v", err)
	}
	if !reflect.DeepEqual(list, workers) {
		t.Errorf("ListWorkers() is %v, want %v", list, workers)
	}
	counter++

	// lookup the workers by hostname
	for _, worker := range workers {
		got, err := db.GetWorkerForHostname(worker.GetHostname())
		if err != nil {
			t.Errorf("unable to get worker %s by hostname: %v", worker.GetHostname(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorkerForHostname() is %v, want %v", got, worker)
		}
	}
	counter++

	// update the workers
	for _, worker := range workers {
		worker.SetActive(false)
		err = db.UpdateWorker(worker)
		if err != nil {
			t.Errorf("unable to update worker %s: %v", worker.GetHostname(), err)
		}

		// lookup the worker by ID
		got, err := db.GetWorker(worker.GetID())
		if err != nil {
			t.Errorf("unable to get worker %s by ID: %v", worker.GetHostname(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorker() is %v, want %v", got, worker)
		}
	}
	counter++
	counter++

	// delete the workers
	for _, worker := range workers {
		err = db.DeleteWorker(worker)
		if err != nil {
			t.Errorf("unable to delete worker %s: %v", worker.GetHostname(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(worker.WorkerInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}
