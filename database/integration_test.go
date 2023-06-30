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

	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
	"github.com/google/go-cmp/cmp"
	"github.com/kr/pretty"
)

type Resources struct {
	Builds      []*library.Build
	Deployments []*library.Deployment
	Hooks       []*library.Hook
	Logs        []*library.Log
	Pipelines   []*library.Pipeline
	Repos       []*library.Repo
	Schedules   []*library.Schedule
	Secrets     []*library.Secret
	Services    []*library.Service
	Steps       []*library.Step
	Users       []*library.User
	Workers     []*library.Worker
}

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
			resources := newResources()

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

			t.Run("test_builds", func(t *testing.T) { testBuilds(t, db, resources) })

			t.Run("test_hooks", func(t *testing.T) { testHooks(t, db, resources) })

			t.Run("test_logs", func(t *testing.T) { testLogs(t, db, resources) })

			t.Run("test_pipelines", func(t *testing.T) { testPipelines(t, db, resources) })

			t.Run("test_repos", func(t *testing.T) { testRepos(t, db, resources) })

			t.Run("test_schedules", func(t *testing.T) { testSchedules(t, db, resources) })

			t.Run("test_secrets", func(t *testing.T) { testSecrets(t, db, resources) })

			t.Run("test_services", func(t *testing.T) { testServices(t, db, resources) })

			t.Run("test_steps", func(t *testing.T) { testSteps(t, db, resources) })

			t.Run("test_users", func(t *testing.T) { testUsers(t, db, resources) })

			t.Run("test_workers", func(t *testing.T) { testWorkers(t, db, resources) })

			err = db.Close()
			if err != nil {
				t.Errorf("unable to close database engine for %s: %v", test.name, err)
			}
		})
	}
}

func testBuilds(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for builds
	//
	// we start at 2 for creating the table and indexes for builds
	// since those are already called when the database engine starts
	counter := 2

	// create the repos for build related functions
	for _, repo := range resources.Repos {
		err := db.CreateRepo(repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	buildOne := new(library.BuildQueue)
	buildOne.SetCreated(1563474076)
	buildOne.SetFullName("github/octokitty")
	buildOne.SetNumber(1)
	buildOne.SetStatus("running")

	buildTwo := new(library.BuildQueue)
	buildTwo.SetCreated(1563474076)
	buildTwo.SetFullName("github/octokitty")
	buildTwo.SetNumber(2)
	buildTwo.SetStatus("running")

	queueBuilds := []*library.BuildQueue{buildOne, buildTwo}

	// create the builds
	for _, build := range resources.Builds {
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
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuilds() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// count the builds for a deployment
	count, err = db.CountBuildsForDeployment(resources.Deployments[0], nil)
	if err != nil {
		t.Errorf("unable to count builds for deployment %d: %v", resources.Deployments[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForDeployment() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// count the builds for an org
	count, err = db.CountBuildsForOrg(resources.Repos[0].GetOrg(), nil)
	if err != nil {
		t.Errorf("unable to count builds for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForOrg() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// count the builds for a repo
	count, err = db.CountBuildsForRepo(resources.Repos[0], nil)
	if err != nil {
		t.Errorf("unable to count builds for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// count the builds for a status
	count, err = db.CountBuildsForStatus("running", nil)
	if err != nil {
		t.Errorf("unable to count builds for status %s: %v", "running", err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForStatus() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// list the builds
	list, err := db.ListBuilds()
	if err != nil {
		t.Errorf("unable to list builds: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Builds) {
		t.Errorf("ListBuilds() is %v, want %v", list, resources.Builds)
	}
	counter++

	// list the builds for a deployment
	list, count, err = db.ListBuildsForDeployment(resources.Deployments[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for deployment %d: %v", resources.Deployments[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForDeployment() is %v, want %v", count, len(resources.Builds))
	}
	if !reflect.DeepEqual(list, []*library.Build{resources.Builds[1], resources.Builds[0]}) {
		t.Errorf("ListBuildsForDeployment() is %v, want %v", list, []*library.Build{resources.Builds[1], resources.Builds[0]})
	}
	counter++

	// list the builds for an org
	list, count, err = db.ListBuildsForOrg(resources.Repos[0].GetOrg(), nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForOrg() is %v, want %v", count, len(resources.Builds))
	}
	if !reflect.DeepEqual(list, resources.Builds) {
		t.Errorf("ListBuildsForOrg() is %v, want %v", list, resources.Builds)
	}
	counter++

	// list the builds for a repo
	list, count, err = db.ListBuildsForRepo(resources.Repos[0], nil, time.Now().UTC().Unix(), 0, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	if !reflect.DeepEqual(list, []*library.Build{resources.Builds[1], resources.Builds[0]}) {
		t.Errorf("ListBuildsForRepo() is %v, want %v", list, []*library.Build{resources.Builds[1], resources.Builds[0]})
	}
	counter++

	// list the pending and running builds
	queueList, err := db.ListPendingAndRunningBuilds("0")
	if err != nil {
		t.Errorf("unable to list pending and running builds: %v", err)
	}
	if !reflect.DeepEqual(queueList, queueBuilds) {
		t.Errorf("ListPendingAndRunningBuilds() is %v, want %v", queueList, queueBuilds)
	}
	counter++

	// lookup the last build by repo
	got, err := db.LastBuildForRepo(resources.Repos[0], "main")
	if err != nil {
		t.Errorf("unable to get last build for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if !reflect.DeepEqual(got, resources.Builds[1]) {
		t.Errorf("GetBuildForRepo() is %v, want %v", got, resources.Builds[1])
	}
	counter++

	// lookup the builds by repo and number
	for _, build := range resources.Builds {
		repo := resources.Repos[build.GetRepoID()-1]
		got, err = db.GetBuildForRepo(repo, build.GetNumber())
		if err != nil {
			t.Errorf("unable to get build %d for repo %d: %v", build.GetID(), repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, build) {
			t.Errorf("GetBuildForRepo() is %v, want %v", got, build)
		}
	}
	counter++

	// update the builds
	for _, build := range resources.Builds {
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

	// clean the builds
	count, err = db.CleanBuilds("msg", time.Now().UTC().Unix())
	if err != nil {
		t.Errorf("unable to clean builds: %v", err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CleanBuilds() is %v, want %v", count, len(resources.Builds))
	}
	counter++

	// delete the builds
	for _, build := range resources.Builds {
		err = db.DeleteBuild(build)
		if err != nil {
			t.Errorf("unable to delete build %d: %v", build.GetID(), err)
		}
	}
	counter++

	// delete the repos for build related functions
	for _, repo := range resources.Repos {
		err := db.DeleteRepo(repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(build.BuildInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testHooks(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for hooks
	//
	// we start at 2 for creating the table and indexes for hooks
	// since those are already called when the database engine starts
	counter := 2

	// create the hooks
	for _, hook := range resources.Hooks {
		_, err := db.CreateHook(hook)
		if err != nil {
			t.Errorf("unable to create hook %d: %v", hook.GetID(), err)
		}
	}
	counter++

	// count the hooks
	count, err := db.CountHooks()
	if err != nil {
		t.Errorf("unable to count hooks: %v", err)
	}
	if int(count) != len(resources.Hooks) {
		t.Errorf("CountHooks() is %v, want %v", count, len(resources.Hooks))
	}
	counter++

	// list the hooks
	list, err := db.ListHooks()
	if err != nil {
		t.Errorf("unable to list hooks: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Hooks) {
		t.Errorf("ListHooks() is %v, want %v", list, resources.Hooks)
	}
	counter++

	// lookup the hooks by name
	for _, hook := range resources.Hooks {
		repo := resources.Repos[hook.GetRepoID()-1]
		got, err := db.GetHookForRepo(repo, hook.GetNumber())
		if err != nil {
			t.Errorf("unable to get hook %d for repo %d: %v", hook.GetID(), repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, hook) {
			t.Errorf("GetHookForRepo() is %v, want %v", got, hook)
		}
	}
	counter++

	// update the hooks
	for _, hook := range resources.Hooks {
		hook.SetStatus("success")
		_, err = db.UpdateHook(hook)
		if err != nil {
			t.Errorf("unable to update hook %d: %v", hook.GetID(), err)
		}

		// lookup the hook by ID
		got, err := db.GetHook(hook.GetID())
		if err != nil {
			t.Errorf("unable to get hook %d by ID: %v", hook.GetID(), err)
		}
		if !reflect.DeepEqual(got, hook) {
			t.Errorf("GetHook() is %v, want %v", got, hook)
		}
	}
	counter++
	counter++

	// delete the hooks
	for _, hook := range resources.Hooks {
		err = db.DeleteHook(hook)
		if err != nil {
			t.Errorf("unable to delete hook %d: %v", hook.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(hook.HookInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testLogs(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for logs
	//
	// we start at 2 for creating the table and indexes for logs
	// since those are already called when the database engine starts
	counter := 2

	// create the logs
	for _, log := range resources.Logs {
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
	if int(count) != len(resources.Logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(resources.Logs))
	}
	counter++

	// list the logs
	list, err := db.ListLogs()
	if err != nil {
		t.Errorf("unable to list logs: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Logs) {
		t.Errorf("ListLogs() is %v, want %v", list, resources.Logs)
	}
	counter++

	// lookup the logs by service
	for _, log := range []*library.Log{resources.Logs[0], resources.Logs[1]} {
		service := resources.Services[log.GetServiceID()-1]
		got, err := db.GetLogForService(service)
		if err != nil {
			t.Errorf("unable to get log %d for service %d: %v", log.GetID(), service.GetID(), err)
		}
		if !reflect.DeepEqual(got, log) {
			t.Errorf("GetLogForService() is %v, want %v", got, log)
		}
	}
	counter++

	// lookup the logs by service
	for _, log := range []*library.Log{resources.Logs[2], resources.Logs[3]} {
		step := resources.Steps[log.GetStepID()-1]
		got, err := db.GetLogForStep(step)
		if err != nil {
			t.Errorf("unable to get log %d for step %d: %v", log.GetID(), step.GetID(), err)
		}
		if !reflect.DeepEqual(got, log) {
			t.Errorf("GetLogForStep() is %v, want %v", got, log)
		}
	}
	counter++

	// update the logs
	for _, log := range resources.Logs {
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
	for _, log := range resources.Logs {
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

func testPipelines(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for pipelines
	//
	// we start at 2 for creating the table and indexes for pipelines
	// since those are already called when the database engine starts
	counter := 2

	// create the pipelines
	for _, pipeline := range resources.Pipelines {
		_, err := db.CreatePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to create pipeline %d: %v", pipeline.GetID(), err)
		}
	}
	counter++

	// count the pipelines
	count, err := db.CountPipelines()
	if err != nil {
		t.Errorf("unable to count pipelines: %v", err)
	}
	if int(count) != len(resources.Pipelines) {
		t.Errorf("CountPipelines() is %v, want %v", count, len(resources.Pipelines))
	}
	counter++

	// list the pipelines
	list, err := db.ListPipelines()
	if err != nil {
		t.Errorf("unable to list pipelines: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Pipelines) {
		t.Errorf("ListPipelines() is %v, want %v", list, resources.Pipelines)
	}
	counter++

	// lookup the pipelines by name
	for _, pipeline := range resources.Pipelines {
		repo := resources.Repos[pipeline.GetRepoID()-1]
		got, err := db.GetPipelineForRepo(pipeline.GetCommit(), repo)
		if err != nil {
			t.Errorf("unable to get pipeline %d for repo %d: %v", pipeline.GetID(), repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, pipeline) {
			t.Errorf("GetPipelineForRepo() is %v, want %v", got, pipeline)
		}
	}
	counter++

	// update the pipelines
	for _, pipeline := range resources.Pipelines {
		pipeline.SetVersion("2")
		_, err = db.UpdatePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to update pipeline %d: %v", pipeline.GetID(), err)
		}

		// lookup the pipeline by ID
		got, err := db.GetPipeline(pipeline.GetID())
		if err != nil {
			t.Errorf("unable to get pipeline %d by ID: %v", pipeline.GetID(), err)
		}
		if !reflect.DeepEqual(got, pipeline) {
			t.Errorf("GetPipeline() is %v, want %v", got, pipeline)
		}
	}
	counter++
	counter++

	// delete the pipelines
	for _, pipeline := range resources.Pipelines {
		err = db.DeletePipeline(pipeline)
		if err != nil {
			t.Errorf("unable to delete pipeline %d: %v", pipeline.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(pipeline.PipelineInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testRepos(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for repos
	//
	// we start at 2 for creating the table and indexes for repos
	// since those are already called when the database engine starts
	counter := 2

	// create the repos
	for _, repo := range resources.Repos {
		err := db.CreateRepo(repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}
	counter++

	// count the repos
	count, err := db.CountRepos()
	if err != nil {
		t.Errorf("unable to count repos: %v", err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("CountRepos() is %v, want %v", count, len(resources.Repos))
	}
	counter++

	// list the repos
	list, err := db.ListRepos()
	if err != nil {
		t.Errorf("unable to list repos: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Repos) {
		t.Errorf("ListRepos() is %v, want %v", list, resources.Repos)
	}
	counter++

	// lookup the repos by name
	for _, repo := range resources.Repos {
		got, err := db.GetRepoForOrg(repo.GetOrg(), repo.GetName())
		if err != nil {
			t.Errorf("unable to get repo %d by org: %v", repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, repo) {
			t.Errorf("GetRepoForOrg() is %v, want %v", got, repo)
		}
	}
	counter++

	// update the repos
	for _, repo := range resources.Repos {
		repo.SetActive(false)
		err = db.UpdateRepo(repo)
		if err != nil {
			t.Errorf("unable to update repo %d: %v", repo.GetID(), err)
		}

		// lookup the repo by ID
		got, err := db.GetRepo(repo.GetID())
		if err != nil {
			t.Errorf("unable to get repo %d by ID: %v", repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, repo) {
			t.Errorf("GetRepo() is %v, want %v", got, repo)
		}
	}
	counter++
	counter++

	// delete the repos
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(repo.RepoInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testSchedules(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for schedules
	//
	// we start at 2 for creating the table and indexes for schedules
	// since those are already called when the database engine starts
	counter := 2

	// create the schedules
	for _, schedule := range resources.Schedules {
		err := db.CreateSchedule(schedule)
		if err != nil {
			t.Errorf("unable to create schedule %d: %v", schedule.GetID(), err)
		}
	}
	counter++

	// count the schedules
	count, err := db.CountSchedules()
	if err != nil {
		t.Errorf("unable to count schedules: %v", err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("CountSchedules() is %v, want %v", count, len(resources.Schedules))
	}
	counter++

	// list the schedules
	list, err := db.ListSchedules()
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Schedules) {
		t.Errorf("ListSchedules() is %v, want %v", list, resources.Schedules)
	}
	counter++

	// lookup the schedules by name
	for _, schedule := range resources.Schedules {
		repo := resources.Repos[schedule.GetRepoID()-1]
		got, err := db.GetScheduleForRepo(repo, schedule.GetName())
		if err != nil {
			t.Errorf("unable to get schedule %d for repo %d: %v", schedule.GetID(), repo.GetID(), err)
		}
		if !reflect.DeepEqual(got, schedule) {
			t.Errorf("GetScheduleForRepo() is %v, want %v", got, schedule)
		}
	}
	counter++

	// update the schedules
	for _, schedule := range resources.Schedules {
		schedule.SetActive(false)
		err = db.UpdateSchedule(schedule, true)
		if err != nil {
			t.Errorf("unable to update schedule %d: %v", schedule.GetID(), err)
		}

		// lookup the schedule by ID
		got, err := db.GetSchedule(schedule.GetID())
		if err != nil {
			t.Errorf("unable to get schedule %d by ID: %v", schedule.GetID(), err)
		}
		if !reflect.DeepEqual(got, schedule) {
			t.Errorf("GetSchedule() is %v, want %v", got, schedule)
		}
	}
	counter++
	counter++

	// delete the schedules
	for _, schedule := range resources.Schedules {
		err = db.DeleteSchedule(schedule)
		if err != nil {
			t.Errorf("unable to delete schedule %d: %v", schedule.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(schedule.ScheduleInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testSecrets(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for secrets
	//
	// we start at 2 for creating the table and indexes for secrets
	// since those are already called when the database engine starts
	counter := 2

	// create the secrets
	for _, secret := range resources.Secrets {
		err := db.CreateSecret(secret)
		if err != nil {
			t.Errorf("unable to create secret %d: %v", secret.GetID(), err)
		}
	}
	counter++

	// count the secrets
	count, err := db.CountSecrets()
	if err != nil {
		t.Errorf("unable to count secrets: %v", err)
	}
	if int(count) != len(resources.Secrets) {
		t.Errorf("CountSecrets() is %v, want %v", count, len(resources.Secrets))
	}
	counter++

	// list the secrets
	list, err := db.ListSecrets()
	if err != nil {
		t.Errorf("unable to list secrets: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Secrets) {
		t.Errorf("ListSecrets() is %v, want %v", list, resources.Secrets)
	}
	counter++

	// lookup the secrets by name
	for _, secret := range resources.Secrets {
		got, err := db.GetSecretForRepo(secret.GetName(), resources.Repos[0])
		if err != nil {
			t.Errorf("unable to get secret %d for repo %d: %v", secret.GetID(), resources.Repos[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, secret) {
			t.Errorf("GetSecretForRepo() is %v, want %v", got, secret)
		}
	}
	counter++

	// update the secrets
	for _, secret := range resources.Secrets {
		secret.SetAllowCommand(false)
		err = db.UpdateSecret(secret)
		if err != nil {
			t.Errorf("unable to update secret %d: %v", secret.GetID(), err)
		}

		// lookup the secret by ID
		got, err := db.GetSecret(secret.GetID())
		if err != nil {
			t.Errorf("unable to get secret %d by ID: %v", secret.GetID(), err)
		}
		if !reflect.DeepEqual(got, secret) {
			t.Errorf("GetSecret() is %v, want %v", got, secret)
		}
	}
	counter++
	counter++

	// delete the secrets
	for _, secret := range resources.Secrets {
		err = db.DeleteSecret(secret)
		if err != nil {
			t.Errorf("unable to delete secret %d: %v", secret.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(secret.SecretInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testServices(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for services
	//
	// we start at 2 for creating the table and indexes for services
	// since those are already called when the database engine starts
	counter := 2

	// create the services
	for _, service := range resources.Services {
		err := db.CreateService(service)
		if err != nil {
			t.Errorf("unable to create service %d: %v", service.GetID(), err)
		}
	}
	counter++

	// count the services
	count, err := db.CountServices()
	if err != nil {
		t.Errorf("unable to count services: %v", err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServices() is %v, want %v", count, len(resources.Services))
	}
	counter++

	// count the services for a build
	count, err = db.CountServicesForBuild(resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServicesForBuild() is %v, want %v", count, len(resources.Services))
	}
	counter++

	// list the services
	list, err := db.ListServices()
	if err != nil {
		t.Errorf("unable to list services: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Services) {
		t.Errorf("ListServices() is %v, want %v", list, resources.Services)
	}
	counter++

	// list the services for a build
	list, count, err = db.ListServicesForBuild(resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Service{resources.Services[1], resources.Services[0]}) {
		t.Errorf("ListServicesForBuild() is %v, want %v", list, []*library.Service{resources.Services[1], resources.Services[0]})
	}
	if int(count) != len(resources.Services) {
		t.Errorf("ListServicesForBuild() is %v, want %v", count, len(resources.Services))
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
	for _, service := range resources.Services {
		build := resources.Builds[service.GetBuildID()-1]
		got, err := db.GetServiceForBuild(build, service.GetNumber())
		if err != nil {
			t.Errorf("unable to get service %d for build %d: %v", service.GetID(), build.GetID(), err)
		}
		if !reflect.DeepEqual(got, service) {
			t.Errorf("GetServiceForBuild() is %v, want %v", got, service)
		}
	}
	counter++

	// update the services
	for _, service := range resources.Services {
		service.SetStatus("success")
		err = db.UpdateService(service)
		if err != nil {
			t.Errorf("unable to update service %d: %v", service.GetID(), err)
		}

		// lookup the service by ID
		got, err := db.GetService(service.GetID())
		if err != nil {
			t.Errorf("unable to get service %d by ID: %v", service.GetID(), err)
		}
		if !reflect.DeepEqual(got, service) {
			t.Errorf("GetService() is %v, want %v", got, service)
		}
	}
	counter++
	counter++

	// delete the services
	for _, service := range resources.Services {
		err = db.DeleteService(service)
		if err != nil {
			t.Errorf("unable to delete service %d: %v", service.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(service.ServiceInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testSteps(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for steps
	//
	// we start at 2 for creating the table and indexes for steps
	// since those are already called when the database engine starts
	counter := 2

	// create the steps
	for _, step := range resources.Steps {
		err := db.CreateStep(step)
		if err != nil {
			t.Errorf("unable to create step %d: %v", step.GetID(), err)
		}
	}
	counter++

	// count the steps
	count, err := db.CountSteps()
	if err != nil {
		t.Errorf("unable to count steps: %v", err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountSteps() is %v, want %v", count, len(resources.Steps))
	}
	counter++

	// count the steps for a build
	count, err = db.CountStepsForBuild(resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountStepsForBuild() is %v, want %v", count, len(resources.Steps))
	}
	counter++

	// list the steps
	list, err := db.ListSteps()
	if err != nil {
		t.Errorf("unable to list steps: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Steps) {
		t.Errorf("ListSteps() is %v, want %v", list, resources.Steps)
	}
	counter++

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Step{resources.Steps[1], resources.Steps[0]}) {
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*library.Step{resources.Steps[1], resources.Steps[0]})
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("ListStepsForBuild() is %v, want %v", count, len(resources.Steps))
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
	for _, step := range resources.Steps {
		build := resources.Builds[step.GetBuildID()-1]
		got, err := db.GetStepForBuild(build, step.GetNumber())
		if err != nil {
			t.Errorf("unable to get step %d for build %d: %v", step.GetID(), build.GetID(), err)
		}
		if !reflect.DeepEqual(got, step) {
			t.Errorf("GetStepForBuild() is %v, want %v", got, step)
		}
	}
	counter++

	// update the steps
	for _, step := range resources.Steps {
		step.SetStatus("success")
		err = db.UpdateStep(step)
		if err != nil {
			t.Errorf("unable to update step %d: %v", step.GetID(), err)
		}

		// lookup the step by ID
		got, err := db.GetStep(step.GetID())
		if err != nil {
			t.Errorf("unable to get step %d by ID: %v", step.GetID(), err)
		}
		if !reflect.DeepEqual(got, step) {
			t.Errorf("GetStep() is %v, want %v", got, step)
		}
	}
	counter++
	counter++

	// delete the steps
	for _, step := range resources.Steps {
		err = db.DeleteStep(step)
		if err != nil {
			t.Errorf("unable to delete step %d: %v", step.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(step.StepInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testUsers(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for users
	//
	// we start at 2 for creating the table and indexes for users
	// since those are already called when the database engine starts
	counter := 2

	userOne := new(library.User)
	userOne.SetID(1)
	userOne.SetName("octocat")
	userOne.SetToken("")
	userOne.SetRefreshToken("")
	userOne.SetHash("")
	userOne.SetActive(false)
	userOne.SetAdmin(false)

	userTwo := new(library.User)
	userTwo.SetID(2)
	userTwo.SetName("octokitty")
	userTwo.SetToken("")
	userTwo.SetRefreshToken("")
	userTwo.SetHash("")
	userTwo.SetActive(false)
	userTwo.SetAdmin(false)

	liteUsers := []*library.User{userOne, userTwo}

	// create the users
	for _, user := range resources.Users {
		err := db.CreateUser(user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}
	counter++

	// count the users
	count, err := db.CountUsers()
	if err != nil {
		t.Errorf("unable to count users: %v", err)
	}
	if int(count) != len(resources.Users) {
		t.Errorf("CountUsers() is %v, want %v", count, len(resources.Users))
	}
	counter++

	// list the users
	list, err := db.ListUsers()
	if err != nil {
		t.Errorf("unable to list users: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Users) {
		t.Errorf("ListUsers() is %v, want %v", list, resources.Users)
	}
	counter++

	// lite list the users
	list, count, err = db.ListLiteUsers(1, 10)
	if err != nil {
		t.Errorf("unable to list lite users: %v", err)
	}
	if !reflect.DeepEqual(list, liteUsers) {
		pretty.Ldiff(t, list, liteUsers)
		if diff := cmp.Diff(liteUsers, list); diff != "" {
			t.Errorf("ListLiteUsers() mismatch (-want +got):\n%s", diff)
		}
		t.Errorf("ListLiteUsers() is %v, want %v", list, liteUsers)
	}
	if int(count) != len(liteUsers) {
		t.Errorf("ListLiteUsers() is %v, want %v", count, len(liteUsers))
	}
	counter++

	// lookup the users by name
	for _, user := range resources.Users {
		got, err := db.GetUserForName(user.GetName())
		if err != nil {
			t.Errorf("unable to get user %d by name: %v", user.GetID(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUserForName() is %v, want %v", got, user)
		}
	}
	counter++

	// update the users
	for _, user := range resources.Users {
		user.SetActive(false)
		err = db.UpdateUser(user)
		if err != nil {
			t.Errorf("unable to update user %d: %v", user.GetID(), err)
		}

		// lookup the user by ID
		got, err := db.GetUser(user.GetID())
		if err != nil {
			t.Errorf("unable to get user %d by ID: %v", user.GetID(), err)
		}
		if !reflect.DeepEqual(got, user) {
			t.Errorf("GetUser() is %v, want %v", got, user)
		}
	}
	counter++
	counter++

	// delete the users
	for _, user := range resources.Users {
		err = db.DeleteUser(user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(user.UserInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func testWorkers(t *testing.T, db Interface, resources *Resources) {
	// used to track the number of methods we call for workers
	//
	// we start at 2 for creating the table and indexes for users
	// since those are already called when the database engine starts
	counter := 2

	// create the workers
	for _, worker := range resources.Workers {
		err := db.CreateWorker(worker)
		if err != nil {
			t.Errorf("unable to create worker %d: %v", worker.GetID(), err)
		}
	}
	counter++

	// count the workers
	count, err := db.CountWorkers()
	if err != nil {
		t.Errorf("unable to count workers: %v", err)
	}
	if int(count) != len(resources.Workers) {
		t.Errorf("CountWorkers() is %v, want %v", count, len(resources.Workers))
	}
	counter++

	// list the workers
	list, err := db.ListWorkers()
	if err != nil {
		t.Errorf("unable to list workers: %v", err)
	}
	if !reflect.DeepEqual(list, resources.Workers) {
		t.Errorf("ListWorkers() is %v, want %v", list, resources.Workers)
	}
	counter++

	// lookup the workers by hostname
	for _, worker := range resources.Workers {
		got, err := db.GetWorkerForHostname(worker.GetHostname())
		if err != nil {
			t.Errorf("unable to get worker %d by hostname: %v", worker.GetID(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorkerForHostname() is %v, want %v", got, worker)
		}
	}
	counter++

	// update the workers
	for _, worker := range resources.Workers {
		worker.SetActive(false)
		err = db.UpdateWorker(worker)
		if err != nil {
			t.Errorf("unable to update worker %d: %v", worker.GetID(), err)
		}

		// lookup the worker by ID
		got, err := db.GetWorker(worker.GetID())
		if err != nil {
			t.Errorf("unable to get worker %d by ID: %v", worker.GetID(), err)
		}
		if !reflect.DeepEqual(got, worker) {
			t.Errorf("GetWorker() is %v, want %v", got, worker)
		}
	}
	counter++
	counter++

	// delete the workers
	for _, worker := range resources.Workers {
		err = db.DeleteWorker(worker)
		if err != nil {
			t.Errorf("unable to delete worker %d: %v", worker.GetID(), err)
		}
	}
	counter++

	// ensure we called all the functions we should have
	methods := reflect.TypeOf(new(worker.WorkerInterface)).Elem().NumMethod()
	if counter != methods {
		t.Errorf("total number of methods called is %v, want %v", counter, methods)
	}
}

func newResources() *Resources {
	buildOne := new(library.Build)
	buildOne.SetID(1)
	buildOne.SetRepoID(1)
	buildOne.SetPipelineID(1)
	buildOne.SetNumber(1)
	buildOne.SetParent(1)
	buildOne.SetEvent("push")
	buildOne.SetEventAction("")
	buildOne.SetStatus("running")
	buildOne.SetError("")
	buildOne.SetEnqueued(1563474077)
	buildOne.SetCreated(1563474076)
	buildOne.SetStarted(1563474078)
	buildOne.SetFinished(1563474079)
	buildOne.SetDeploy("")
	buildOne.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	buildOne.SetClone("https://github.com/github/octocat.git")
	buildOne.SetSource("https://github.com/github/octocat/deployments/1")
	buildOne.SetTitle("push received from https://github.com/github/octocat")
	buildOne.SetMessage("First commit...")
	buildOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	buildOne.SetSender("OctoKitty")
	buildOne.SetAuthor("OctoKitty")
	buildOne.SetEmail("OctoKitty@github.com")
	buildOne.SetLink("https://example.company.com/github/octocat/1")
	buildOne.SetBranch("main")
	buildOne.SetRef("refs/heads/main")
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
	buildTwo.SetEventAction("")
	buildTwo.SetStatus("running")
	buildTwo.SetError("")
	buildTwo.SetEnqueued(1563474077)
	buildTwo.SetCreated(1563474076)
	buildTwo.SetStarted(1563474078)
	buildTwo.SetFinished(1563474079)
	buildTwo.SetDeploy("")
	buildTwo.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	buildTwo.SetClone("https://github.com/github/octocat.git")
	buildTwo.SetSource("https://github.com/github/octocat/deployments/1")
	buildTwo.SetTitle("pull_request received from https://github.com/github/octocat")
	buildTwo.SetMessage("Second commit...")
	buildTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	buildTwo.SetSender("OctoKitty")
	buildTwo.SetAuthor("OctoKitty")
	buildTwo.SetEmail("OctoKitty@github.com")
	buildTwo.SetLink("https://example.company.com/github/octocat/2")
	buildTwo.SetBranch("main")
	buildTwo.SetRef("refs/heads/main")
	buildTwo.SetBaseRef("")
	buildTwo.SetHeadRef("changes")
	buildTwo.SetHost("example.company.com")
	buildTwo.SetRuntime("docker")
	buildTwo.SetDistribution("linux")

	deploymentOne := new(library.Deployment)
	deploymentOne.SetID(1)
	deploymentOne.SetRepoID(1)
	deploymentOne.SetURL("https://github.com/github/octocat/deployments/1")
	deploymentOne.SetUser("octocat")
	deploymentOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	deploymentOne.SetRef("refs/heads/master")
	deploymentOne.SetTask("vela-deploy")
	deploymentOne.SetTarget("production")
	deploymentOne.SetDescription("Deployment request from Vela")
	deploymentOne.SetPayload(map[string]string{"foo": "test1"})

	deploymentTwo := new(library.Deployment)
	deploymentTwo.SetID(1)
	deploymentTwo.SetRepoID(1)
	deploymentTwo.SetURL("https://github.com/github/octocat/deployments/2")
	deploymentTwo.SetUser("octocat")
	deploymentTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	deploymentTwo.SetRef("refs/heads/master")
	deploymentTwo.SetTask("vela-deploy")
	deploymentTwo.SetTarget("production")
	deploymentTwo.SetDescription("Deployment request from Vela")
	deploymentTwo.SetPayload(map[string]string{"foo": "test1"})

	hookOne := new(library.Hook)
	hookOne.SetID(1)
	hookOne.SetRepoID(1)
	hookOne.SetBuildID(1)
	hookOne.SetNumber(1)
	hookOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	hookOne.SetCreated(time.Now().UTC().Unix())
	hookOne.SetHost("github.com")
	hookOne.SetEvent("push")
	hookOne.SetEventAction("")
	hookOne.SetBranch("main")
	hookOne.SetError("")
	hookOne.SetStatus("success")
	hookOne.SetLink("https://github.com/github/octocat/settings/hooks/1")
	hookOne.SetWebhookID(123456)

	hookTwo := new(library.Hook)
	hookTwo.SetID(2)
	hookTwo.SetRepoID(1)
	hookTwo.SetBuildID(1)
	hookTwo.SetNumber(2)
	hookTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	hookTwo.SetCreated(time.Now().UTC().Unix())
	hookTwo.SetHost("github.com")
	hookTwo.SetEvent("push")
	hookTwo.SetEventAction("")
	hookTwo.SetBranch("main")
	hookTwo.SetError("")
	hookTwo.SetStatus("success")
	hookTwo.SetLink("https://github.com/github/octocat/settings/hooks/1")
	hookTwo.SetWebhookID(123456)

	logServiceOne := new(library.Log)
	logServiceOne.SetID(1)
	logServiceOne.SetBuildID(1)
	logServiceOne.SetRepoID(1)
	logServiceOne.SetServiceID(1)
	logServiceOne.SetStepID(0)
	logServiceOne.SetData([]byte("foo"))

	logServiceTwo := new(library.Log)
	logServiceTwo.SetID(2)
	logServiceTwo.SetBuildID(1)
	logServiceTwo.SetRepoID(1)
	logServiceTwo.SetServiceID(2)
	logServiceTwo.SetStepID(0)
	logServiceTwo.SetData([]byte("foo"))

	logStepOne := new(library.Log)
	logStepOne.SetID(3)
	logStepOne.SetBuildID(1)
	logStepOne.SetRepoID(1)
	logStepOne.SetServiceID(0)
	logStepOne.SetStepID(1)
	logStepOne.SetData([]byte("foo"))

	logStepTwo := new(library.Log)
	logStepTwo.SetID(4)
	logStepTwo.SetBuildID(1)
	logStepTwo.SetRepoID(1)
	logStepTwo.SetServiceID(0)
	logStepTwo.SetStepID(2)
	logStepTwo.SetData([]byte("foo"))

	pipelineOne := new(library.Pipeline)
	pipelineOne.SetID(1)
	pipelineOne.SetRepoID(1)
	pipelineOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	pipelineOne.SetFlavor("large")
	pipelineOne.SetPlatform("docker")
	pipelineOne.SetRef("refs/heads/main")
	pipelineOne.SetType("yaml")
	pipelineOne.SetVersion("1")
	pipelineOne.SetExternalSecrets(false)
	pipelineOne.SetInternalSecrets(false)
	pipelineOne.SetServices(true)
	pipelineOne.SetStages(false)
	pipelineOne.SetSteps(true)
	pipelineOne.SetTemplates(false)
	pipelineOne.SetData([]byte("version: 1"))

	pipelineTwo := new(library.Pipeline)
	pipelineTwo.SetID(2)
	pipelineTwo.SetRepoID(1)
	pipelineTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	pipelineTwo.SetFlavor("large")
	pipelineTwo.SetPlatform("docker")
	pipelineTwo.SetRef("refs/heads/main")
	pipelineTwo.SetType("yaml")
	pipelineTwo.SetVersion("1")
	pipelineTwo.SetExternalSecrets(false)
	pipelineTwo.SetInternalSecrets(false)
	pipelineTwo.SetServices(true)
	pipelineTwo.SetStages(false)
	pipelineTwo.SetSteps(true)
	pipelineTwo.SetTemplates(false)
	pipelineTwo.SetData([]byte("version: 1"))

	repoOne := new(library.Repo)
	repoOne.SetID(1)
	repoOne.SetUserID(1)
	repoOne.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
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
	repoTwo.SetUserID(1)
	repoTwo.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	repoTwo.SetOrg("github")
	repoTwo.SetName("octokitty")
	repoTwo.SetFullName("github/octokitty")
	repoTwo.SetLink("https://github.com/github/octokitty")
	repoTwo.SetClone("https://github.com/github/octokitty.git")
	repoTwo.SetBranch("main")
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

	scheduleOne := new(library.Schedule)
	scheduleOne.SetID(1)
	scheduleOne.SetRepoID(1)
	scheduleOne.SetActive(true)
	scheduleOne.SetName("nightly")
	scheduleOne.SetEntry("0 0 * * *")
	scheduleOne.SetCreatedAt(time.Now().UTC().Unix())
	scheduleOne.SetCreatedBy("octocat")
	scheduleOne.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	scheduleOne.SetUpdatedBy("octokitty")
	scheduleOne.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())

	scheduleTwo := new(library.Schedule)
	scheduleTwo.SetID(2)
	scheduleTwo.SetRepoID(1)
	scheduleTwo.SetActive(true)
	scheduleTwo.SetName("hourly")
	scheduleTwo.SetEntry("0 * * * *")
	scheduleTwo.SetCreatedAt(time.Now().UTC().Unix())
	scheduleTwo.SetCreatedBy("octocat")
	scheduleTwo.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	scheduleTwo.SetUpdatedBy("octokitty")
	scheduleTwo.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())

	secretOne := new(library.Secret)
	secretOne.SetID(1)
	secretOne.SetOrg("github")
	secretOne.SetRepo("octocat")
	secretOne.SetTeam("")
	secretOne.SetName("foo")
	secretOne.SetValue("bar")
	secretOne.SetType("repo")
	secretOne.SetImages([]string{"alpine"})
	secretOne.SetEvents([]string{"push", "tag", "deployment"})
	secretOne.SetAllowCommand(true)
	secretOne.SetCreatedAt(time.Now().UTC().Unix())
	secretOne.SetCreatedBy("octocat")
	secretOne.SetUpdatedAt(time.Now().UTC().Unix())
	secretOne.SetUpdatedBy("octokitty")

	secretTwo := new(library.Secret)
	secretTwo.SetID(2)
	secretTwo.SetOrg("github")
	secretTwo.SetRepo("octocat")
	secretTwo.SetTeam("")
	secretTwo.SetName("bar")
	secretTwo.SetValue("baz")
	secretTwo.SetType("repo")
	secretTwo.SetImages([]string{"alpine"})
	secretTwo.SetEvents([]string{"push", "tag", "deployment"})
	secretTwo.SetAllowCommand(true)
	secretTwo.SetCreatedAt(time.Now().UTC().Unix())
	secretTwo.SetCreatedBy("octocat")
	secretTwo.SetUpdatedAt(time.Now().UTC().Unix())
	secretTwo.SetUpdatedBy("octokitty")

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

	userOne := new(library.User)
	userOne.SetID(1)
	userOne.SetName("octocat")
	userOne.SetToken("superSecretToken")
	userOne.SetRefreshToken("superSecretRefreshToken")
	userOne.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	userOne.SetFavorites([]string{"github/octocat"})
	userOne.SetActive(true)
	userOne.SetAdmin(false)

	userTwo := new(library.User)
	userTwo.SetID(2)
	userTwo.SetName("octokitty")
	userTwo.SetToken("superSecretToken")
	userTwo.SetRefreshToken("superSecretRefreshToken")
	userTwo.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	userTwo.SetFavorites([]string{"github/octocat"})
	userTwo.SetActive(true)
	userTwo.SetAdmin(false)

	workerOne := new(library.Worker)
	workerOne.SetID(1)
	workerOne.SetHostname("worker-1.example.com")
	workerOne.SetAddress("https://worker-1.example.com")
	workerOne.SetRoutes([]string{"vela"})
	workerOne.SetActive(true)
	workerOne.SetStatus("available")
	workerOne.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	workerOne.SetRunningBuildIDs([]string{"12345"})
	workerOne.SetLastBuildStartedAt(time.Now().UTC().Unix())
	workerOne.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	workerOne.SetLastCheckedIn(time.Now().UTC().Unix())
	workerOne.SetBuildLimit(1)

	workerTwo := new(library.Worker)
	workerTwo.SetID(2)
	workerTwo.SetHostname("worker-2.example.com")
	workerTwo.SetAddress("https://worker-2.example.com")
	workerTwo.SetRoutes([]string{"vela"})
	workerTwo.SetActive(true)
	workerTwo.SetStatus("available")
	workerTwo.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	workerTwo.SetRunningBuildIDs([]string{"12345"})
	workerTwo.SetLastBuildStartedAt(time.Now().UTC().Unix())
	workerTwo.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	workerTwo.SetLastCheckedIn(time.Now().UTC().Unix())
	workerTwo.SetBuildLimit(1)

	return &Resources{
		Builds:      []*library.Build{buildOne, buildTwo},
		Deployments: []*library.Deployment{deploymentOne, deploymentTwo},
		Hooks:       []*library.Hook{hookOne, hookTwo},
		Logs:        []*library.Log{logServiceOne, logServiceTwo, logStepOne, logStepTwo},
		Pipelines:   []*library.Pipeline{pipelineOne, pipelineTwo},
		Repos:       []*library.Repo{repoOne, repoTwo},
		Schedules:   []*library.Schedule{scheduleOne, scheduleTwo},
		Secrets:     []*library.Secret{secretOne, secretTwo},
		Services:    []*library.Service{serviceOne, serviceTwo},
		Steps:       []*library.Step{stepOne, stepTwo},
		Users:       []*library.User{userOne, userTwo},
		Workers:     []*library.Worker{workerOne, workerTwo},
	}
}
