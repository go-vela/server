// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/executable"
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
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/raw"
	"github.com/google/go-cmp/cmp"
)

// Resources represents the object containing test resources.
type Resources struct {
	Builds      []*library.Build
	Deployments []*library.Deployment
	Executables []*library.BuildExecutable
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
		name   string
		config *config
	}{
		{
			name: "postgres",
			config: &config{
				Driver:           "postgres",
				Address:          os.Getenv("POSTGRES_ADDR"),
				CompressionLevel: 3,
				ConnectionLife:   10 * time.Second,
				ConnectionIdle:   5,
				ConnectionOpen:   20,
				EncryptionKey:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
				SkipCreation:     false,
			},
		},
		{
			name: "sqlite3",
			config: &config{
				Driver:           "sqlite3",
				Address:          os.Getenv("SQLITE_ADDR"),
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

			t.Run("test_executables", func(t *testing.T) { testExecutables(t, db, resources) })

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
	// create a variable to track the number of methods called for builds
	methods := make(map[string]bool)
	// capture the element type of the build interface
	element := reflect.TypeOf(new(build.BuildInterface)).Elem()
	// iterate through all methods found in the build interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for builds
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the repos for build related functions
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	buildOne := new(library.BuildQueue)
	buildOne.SetCreated(1563474076)
	buildOne.SetFullName("github/octocat")
	buildOne.SetNumber(1)
	buildOne.SetStatus("running")

	buildTwo := new(library.BuildQueue)
	buildTwo.SetCreated(1563474076)
	buildTwo.SetFullName("github/octocat")
	buildTwo.SetNumber(2)
	buildTwo.SetStatus("running")

	queueBuilds := []*library.BuildQueue{buildOne, buildTwo}

	// create the builds
	for _, build := range resources.Builds {
		_, err := db.CreateBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to create build %d: %v", build.GetID(), err)
		}
	}
	methods["CreateBuild"] = true

	// count the builds
	count, err := db.CountBuilds(context.TODO())
	if err != nil {
		t.Errorf("unable to count builds: %v", err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuilds() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountBuilds"] = true

	// count the builds for a deployment
	count, err = db.CountBuildsForDeployment(context.TODO(), resources.Deployments[0], nil)
	if err != nil {
		t.Errorf("unable to count builds for deployment %d: %v", resources.Deployments[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForDeployment() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountBuildsForDeployment"] = true

	// count the builds for an org
	count, err = db.CountBuildsForOrg(context.TODO(), resources.Repos[0].GetOrg(), nil)
	if err != nil {
		t.Errorf("unable to count builds for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForOrg() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountBuildsForOrg"] = true

	// count the builds for a repo
	count, err = db.CountBuildsForRepo(context.TODO(), resources.Repos[0], nil)
	if err != nil {
		t.Errorf("unable to count builds for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountBuildsForRepo"] = true

	// count the builds for a status
	count, err = db.CountBuildsForStatus(context.TODO(), "running", nil)
	if err != nil {
		t.Errorf("unable to count builds for status %s: %v", "running", err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountBuildsForStatus() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountBuildsForStatus"] = true

	// list the builds
	list, err := db.ListBuilds(context.TODO())
	if err != nil {
		t.Errorf("unable to list builds: %v", err)
	}
	if !cmp.Equal(list, resources.Builds) {
		t.Errorf("ListBuilds() is %v, want %v", list, resources.Builds)
	}
	methods["ListBuilds"] = true

	// list the builds for a deployment
	list, count, err = db.ListBuildsForDeployment(context.TODO(), resources.Deployments[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for deployment %d: %v", resources.Deployments[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForDeployment() is %v, want %v", count, len(resources.Builds))
	}
	if !cmp.Equal(list, []*library.Build{resources.Builds[1], resources.Builds[0]}) {
		t.Errorf("ListBuildsForDeployment() is %v, want %v", list, []*library.Build{resources.Builds[1], resources.Builds[0]})
	}
	methods["ListBuildsForDeployment"] = true

	// list the builds for an org
	list, count, err = db.ListBuildsForOrg(context.TODO(), resources.Repos[0].GetOrg(), nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForOrg() is %v, want %v", count, len(resources.Builds))
	}
	if !cmp.Equal(list, resources.Builds) {
		t.Errorf("ListBuildsForOrg() is %v, want %v", list, resources.Builds)
	}
	methods["ListBuildsForOrg"] = true

	// list the builds for a repo
	list, count, err = db.ListBuildsForRepo(context.TODO(), resources.Repos[0], nil, time.Now().UTC().Unix(), 0, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	if !cmp.Equal(list, []*library.Build{resources.Builds[1], resources.Builds[0]}) {
		t.Errorf("ListBuildsForRepo() is %v, want %v", list, []*library.Build{resources.Builds[1], resources.Builds[0]})
	}
	methods["ListBuildsForRepo"] = true

	// list the pending and running builds
	queueList, err := db.ListPendingAndRunningBuilds(context.TODO(), "0")
	if err != nil {
		t.Errorf("unable to list pending and running builds: %v", err)
	}
	if !cmp.Equal(queueList, queueBuilds) {
		t.Errorf("ListPendingAndRunningBuilds() is %v, want %v", queueList, queueBuilds)
	}
	methods["ListPendingAndRunningBuilds"] = true

	// lookup the last build by repo
	got, err := db.LastBuildForRepo(context.TODO(), resources.Repos[0], "main")
	if err != nil {
		t.Errorf("unable to get last build for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if !cmp.Equal(got, resources.Builds[1]) {
		t.Errorf("LastBuildForRepo() is %v, want %v", got, resources.Builds[1])
	}
	methods["LastBuildForRepo"] = true

	// lookup the builds by repo and number
	for _, build := range resources.Builds {
		repo := resources.Repos[build.GetRepoID()-1]
		got, err = db.GetBuildForRepo(context.TODO(), repo, build.GetNumber())
		if err != nil {
			t.Errorf("unable to get build %d for repo %d: %v", build.GetID(), repo.GetID(), err)
		}
		if !cmp.Equal(got, build) {
			t.Errorf("GetBuildForRepo() is %v, want %v", got, build)
		}
	}
	methods["GetBuildForRepo"] = true

	// clean the builds
	count, err = db.CleanBuilds(context.TODO(), "integration testing", 1563474090)
	if err != nil {
		t.Errorf("unable to clean builds: %v", err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CleanBuilds() is %v, want %v", count, len(resources.Builds))
	}
	methods["CleanBuilds"] = true

	// update the builds
	for _, build := range resources.Builds {
		build.SetStatus("success")
		_, err = db.UpdateBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to update build %d: %v", build.GetID(), err)
		}

		// lookup the build by ID
		got, err = db.GetBuild(context.TODO(), build.GetID())
		if err != nil {
			t.Errorf("unable to get build %d by ID: %v", build.GetID(), err)
		}
		if !cmp.Equal(got, build) {
			t.Errorf("GetBuild() is %v, want %v", got, build)
		}
	}
	methods["UpdateBuild"] = true
	methods["GetBuild"] = true

	// delete the builds
	for _, build := range resources.Builds {
		err = db.DeleteBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to delete build %d: %v", build.GetID(), err)
		}
	}
	methods["DeleteBuild"] = true

	// delete the repos for build related functions
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for builds", method)
		}
	}
}

func testExecutables(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for pipelines
	methods := make(map[string]bool)
	// capture the element type of the pipeline interface
	element := reflect.TypeOf(new(executable.BuildExecutableInterface)).Elem()
	// iterate through all methods found in the pipeline interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for pipelines
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the pipelines
	for _, executable := range resources.Executables {
		err := db.CreateBuildExecutable(context.TODO(), executable)
		if err != nil {
			t.Errorf("unable to create executable %d: %v", executable.GetID(), err)
		}
	}
	methods["CreateBuildExecutable"] = true

	// pop executables for builds
	for _, executable := range resources.Executables {
		got, err := db.PopBuildExecutable(context.TODO(), executable.GetBuildID())
		if err != nil {
			t.Errorf("unable to get executable %d for build %d: %v", executable.GetID(), executable.GetBuildID(), err)
		}
		if !cmp.Equal(got, executable) {
			t.Errorf("PopBuildExecutable() is %v, want %v", got, executable)
		}
	}
	methods["PopBuildExecutable"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for pipelines", method)
		}
	}
}

func testHooks(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for hooks
	methods := make(map[string]bool)
	// capture the element type of the hook interface
	element := reflect.TypeOf(new(hook.HookInterface)).Elem()
	// iterate through all methods found in the hook interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for hooks
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the hooks
	for _, hook := range resources.Hooks {
		_, err := db.CreateHook(context.TODO(), hook)
		if err != nil {
			t.Errorf("unable to create hook %d: %v", hook.GetID(), err)
		}
	}
	methods["CreateHook"] = true

	// count the hooks
	count, err := db.CountHooks(context.TODO())
	if err != nil {
		t.Errorf("unable to count hooks: %v", err)
	}
	if int(count) != len(resources.Hooks) {
		t.Errorf("CountHooks() is %v, want %v", count, len(resources.Hooks))
	}
	methods["CountHooks"] = true

	// count the hooks for a repo
	count, err = db.CountHooksForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to count hooks for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountHooksForRepo() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountHooksForRepo"] = true

	// list the hooks
	list, err := db.ListHooks(context.TODO())
	if err != nil {
		t.Errorf("unable to list hooks: %v", err)
	}
	if !cmp.Equal(list, resources.Hooks) {
		t.Errorf("ListHooks() is %v, want %v", list, resources.Hooks)
	}
	methods["ListHooks"] = true

	// list the hooks for a repo
	list, count, err = db.ListHooksForRepo(context.TODO(), resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list hooks for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Hooks) {
		t.Errorf("ListHooksForRepo() is %v, want %v", count, len(resources.Hooks))
	}
	if !cmp.Equal(list, []*library.Hook{resources.Hooks[1], resources.Hooks[0]}) {
		t.Errorf("ListHooksForRepo() is %v, want %v", list, []*library.Hook{resources.Hooks[1], resources.Hooks[0]})
	}
	methods["ListHooksForRepo"] = true

	// lookup the last build by repo
	got, err := db.LastHookForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to get last hook for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if !cmp.Equal(got, resources.Hooks[1]) {
		t.Errorf("LastHookForRepo() is %v, want %v", got, resources.Hooks[1])
	}
	methods["LastHookForRepo"] = true

	// lookup the hooks by name
	for _, hook := range resources.Hooks {
		repo := resources.Repos[hook.GetRepoID()-1]
		got, err = db.GetHookForRepo(context.TODO(), repo, hook.GetNumber())
		if err != nil {
			t.Errorf("unable to get hook %d for repo %d: %v", hook.GetID(), repo.GetID(), err)
		}
		if !cmp.Equal(got, hook) {
			t.Errorf("GetHookForRepo() is %v, want %v", got, hook)
		}
	}
	methods["GetHookForRepo"] = true

	// update the hooks
	for _, hook := range resources.Hooks {
		hook.SetStatus("success")
		_, err = db.UpdateHook(context.TODO(), hook)
		if err != nil {
			t.Errorf("unable to update hook %d: %v", hook.GetID(), err)
		}

		// lookup the hook by ID
		got, err = db.GetHook(context.TODO(), hook.GetID())
		if err != nil {
			t.Errorf("unable to get hook %d by ID: %v", hook.GetID(), err)
		}
		if !cmp.Equal(got, hook) {
			t.Errorf("GetHook() is %v, want %v", got, hook)
		}
	}
	methods["UpdateHook"] = true
	methods["GetHook"] = true

	// delete the hooks
	for _, hook := range resources.Hooks {
		err = db.DeleteHook(context.TODO(), hook)
		if err != nil {
			t.Errorf("unable to delete hook %d: %v", hook.GetID(), err)
		}
	}
	methods["DeleteHook"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for hooks", method)
		}
	}
}

func testLogs(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for logs
	methods := make(map[string]bool)
	// capture the element type of the log interface
	element := reflect.TypeOf(new(log.LogInterface)).Elem()
	// iterate through all methods found in the log interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for logs
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the logs
	for _, log := range resources.Logs {
		err := db.CreateLog(log)
		if err != nil {
			t.Errorf("unable to create log %d: %v", log.GetID(), err)
		}
	}
	methods["CreateLog"] = true

	// count the logs
	count, err := db.CountLogs()
	if err != nil {
		t.Errorf("unable to count logs: %v", err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(resources.Logs))
	}
	methods["CountLogs"] = true

	// count the logs for a build
	count, err = db.CountLogsForBuild(resources.Builds[0])
	if err != nil {
		t.Errorf("unable to count logs for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(resources.Logs))
	}
	methods["CountLogsForBuild"] = true

	// list the logs
	list, err := db.ListLogs()
	if err != nil {
		t.Errorf("unable to list logs: %v", err)
	}
	if !cmp.Equal(list, resources.Logs) {
		t.Errorf("ListLogs() is %v, want %v", list, resources.Logs)
	}
	methods["ListLogs"] = true

	// list the logs for a build
	list, count, err = db.ListLogsForBuild(resources.Builds[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list logs for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("ListLogsForBuild() is %v, want %v", count, len(resources.Logs))
	}
	if !cmp.Equal(list, resources.Logs) {
		t.Errorf("ListLogsForBuild() is %v, want %v", list, resources.Logs)
	}
	methods["ListLogsForBuild"] = true

	// lookup the logs by service
	for _, log := range []*library.Log{resources.Logs[0], resources.Logs[1]} {
		service := resources.Services[log.GetServiceID()-1]
		got, err := db.GetLogForService(service)
		if err != nil {
			t.Errorf("unable to get log %d for service %d: %v", log.GetID(), service.GetID(), err)
		}
		if !cmp.Equal(got, log) {
			t.Errorf("GetLogForService() is %v, want %v", got, log)
		}
	}
	methods["GetLogForService"] = true

	// lookup the logs by service
	for _, log := range []*library.Log{resources.Logs[2], resources.Logs[3]} {
		step := resources.Steps[log.GetStepID()-1]
		got, err := db.GetLogForStep(step)
		if err != nil {
			t.Errorf("unable to get log %d for step %d: %v", log.GetID(), step.GetID(), err)
		}
		if !cmp.Equal(got, log) {
			t.Errorf("GetLogForStep() is %v, want %v", got, log)
		}
	}
	methods["GetLogForStep"] = true

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
		if !cmp.Equal(got, log) {
			t.Errorf("GetLog() is %v, want %v", got, log)
		}
	}
	methods["UpdateLog"] = true
	methods["GetLog"] = true

	// delete the logs
	for _, log := range resources.Logs {
		err = db.DeleteLog(log)
		if err != nil {
			t.Errorf("unable to delete log %d: %v", log.GetID(), err)
		}
	}
	methods["DeleteLog"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for logs", method)
		}
	}
}

func testPipelines(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for pipelines
	methods := make(map[string]bool)
	// capture the element type of the pipeline interface
	element := reflect.TypeOf(new(pipeline.PipelineInterface)).Elem()
	// iterate through all methods found in the pipeline interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for pipelines
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the pipelines
	for _, pipeline := range resources.Pipelines {
		_, err := db.CreatePipeline(context.TODO(), pipeline)
		if err != nil {
			t.Errorf("unable to create pipeline %d: %v", pipeline.GetID(), err)
		}
	}
	methods["CreatePipeline"] = true

	// count the pipelines
	count, err := db.CountPipelines(context.TODO())
	if err != nil {
		t.Errorf("unable to count pipelines: %v", err)
	}
	if int(count) != len(resources.Pipelines) {
		t.Errorf("CountPipelines() is %v, want %v", count, len(resources.Pipelines))
	}
	methods["CountPipelines"] = true

	// count the pipelines for a repo
	count, err = db.CountPipelinesForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to count pipelines for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Pipelines) {
		t.Errorf("CountPipelinesForRepo() is %v, want %v", count, len(resources.Pipelines))
	}
	methods["CountPipelinesForRepo"] = true

	// list the pipelines
	list, err := db.ListPipelines(context.TODO())
	if err != nil {
		t.Errorf("unable to list pipelines: %v", err)
	}
	if !cmp.Equal(list, resources.Pipelines) {
		t.Errorf("ListPipelines() is %v, want %v", list, resources.Pipelines)
	}
	methods["ListPipelines"] = true

	// list the pipelines for a repo
	list, count, err = db.ListPipelinesForRepo(context.TODO(), resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list pipelines for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Pipelines) {
		t.Errorf("ListPipelinesForRepo() is %v, want %v", count, len(resources.Pipelines))
	}
	if !cmp.Equal(list, resources.Pipelines) {
		t.Errorf("ListPipelines() is %v, want %v", list, resources.Pipelines)
	}
	methods["ListPipelinesForRepo"] = true

	// lookup the pipelines by name
	for _, pipeline := range resources.Pipelines {
		repo := resources.Repos[pipeline.GetRepoID()-1]
		got, err := db.GetPipelineForRepo(context.TODO(), pipeline.GetCommit(), repo)
		if err != nil {
			t.Errorf("unable to get pipeline %d for repo %d: %v", pipeline.GetID(), repo.GetID(), err)
		}
		if !cmp.Equal(got, pipeline) {
			t.Errorf("GetPipelineForRepo() is %v, want %v", got, pipeline)
		}
	}
	methods["GetPipelineForRepo"] = true

	// update the pipelines
	for _, pipeline := range resources.Pipelines {
		pipeline.SetVersion("2")
		_, err = db.UpdatePipeline(context.TODO(), pipeline)
		if err != nil {
			t.Errorf("unable to update pipeline %d: %v", pipeline.GetID(), err)
		}

		// lookup the pipeline by ID
		got, err := db.GetPipeline(context.TODO(), pipeline.GetID())
		if err != nil {
			t.Errorf("unable to get pipeline %d by ID: %v", pipeline.GetID(), err)
		}
		if !cmp.Equal(got, pipeline) {
			t.Errorf("GetPipeline() is %v, want %v", got, pipeline)
		}
	}
	methods["UpdatePipeline"] = true
	methods["GetPipeline"] = true

	// delete the pipelines
	for _, pipeline := range resources.Pipelines {
		err = db.DeletePipeline(context.TODO(), pipeline)
		if err != nil {
			t.Errorf("unable to delete pipeline %d: %v", pipeline.GetID(), err)
		}
	}
	methods["DeletePipeline"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for pipelines", method)
		}
	}
}

func testRepos(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for repos
	methods := make(map[string]bool)
	// capture the element type of the repo interface
	element := reflect.TypeOf(new(repo.RepoInterface)).Elem()
	// iterate through all methods found in the repo interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for repos
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the repos
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}
	methods["CreateRepo"] = true

	// count the repos
	count, err := db.CountRepos(context.TODO())
	if err != nil {
		t.Errorf("unable to count repos: %v", err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("CountRepos() is %v, want %v", count, len(resources.Repos))
	}
	methods["CountRepos"] = true

	// count the repos for an org
	count, err = db.CountReposForOrg(context.TODO(), resources.Repos[0].GetOrg(), nil)
	if err != nil {
		t.Errorf("unable to count repos for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("CountReposForOrg() is %v, want %v", count, len(resources.Repos))
	}
	methods["CountReposForOrg"] = true

	// count the repos for a user
	count, err = db.CountReposForUser(context.TODO(), resources.Users[0], nil)
	if err != nil {
		t.Errorf("unable to count repos for user %d: %v", resources.Users[0].GetID(), err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("CountReposForUser() is %v, want %v", count, len(resources.Repos))
	}
	methods["CountReposForUser"] = true

	// list the repos
	list, err := db.ListRepos(context.TODO())
	if err != nil {
		t.Errorf("unable to list repos: %v", err)
	}
	if !cmp.Equal(list, resources.Repos) {
		t.Errorf("ListRepos() is %v, want %v", list, resources.Repos)
	}
	methods["ListRepos"] = true

	// list the repos for an org
	list, count, err = db.ListReposForOrg(context.TODO(), resources.Repos[0].GetOrg(), "name", nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list repos for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("ListReposForOrg() is %v, want %v", count, len(resources.Repos))
	}
	if !cmp.Equal(list, resources.Repos) {
		t.Errorf("ListReposForOrg() is %v, want %v", list, resources.Repos)
	}
	methods["ListReposForOrg"] = true

	// list the repos for a user
	list, count, err = db.ListReposForUser(context.TODO(), resources.Users[0], "name", nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list repos for user %d: %v", resources.Users[0].GetID(), err)
	}
	if int(count) != len(resources.Repos) {
		t.Errorf("ListReposForUser() is %v, want %v", count, len(resources.Repos))
	}
	if !cmp.Equal(list, resources.Repos) {
		t.Errorf("ListReposForUser() is %v, want %v", list, resources.Repos)
	}
	methods["ListReposForUser"] = true

	// lookup the repos by name
	for _, repo := range resources.Repos {
		got, err := db.GetRepoForOrg(context.TODO(), repo.GetOrg(), repo.GetName())
		if err != nil {
			t.Errorf("unable to get repo %d by org: %v", repo.GetID(), err)
		}
		if !cmp.Equal(got, repo) {
			t.Errorf("GetRepoForOrg() is %v, want %v", got, repo)
		}
	}
	methods["GetRepoForOrg"] = true

	// update the repos
	for _, repo := range resources.Repos {
		repo.SetActive(false)
		_, err = db.UpdateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to update repo %d: %v", repo.GetID(), err)
		}

		// lookup the repo by ID
		got, err := db.GetRepo(context.TODO(), repo.GetID())
		if err != nil {
			t.Errorf("unable to get repo %d by ID: %v", repo.GetID(), err)
		}
		if !cmp.Equal(got, repo) {
			t.Errorf("GetRepo() is %v, want %v", got, repo)
		}
	}
	methods["UpdateRepo"] = true
	methods["GetRepo"] = true

	// delete the repos
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}
	methods["DeleteRepo"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for repos", method)
		}
	}
}

func testSchedules(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for schedules
	methods := make(map[string]bool)
	// capture the element type of the schedule interface
	element := reflect.TypeOf(new(schedule.ScheduleInterface)).Elem()
	// iterate through all methods found in the schedule interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for schedules
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	ctx := context.TODO()

	// create the schedules
	for _, schedule := range resources.Schedules {
		_, err := db.CreateSchedule(ctx, schedule)
		if err != nil {
			t.Errorf("unable to create schedule %d: %v", schedule.GetID(), err)
		}
	}
	methods["CreateSchedule"] = true

	// count the schedules
	count, err := db.CountSchedules(ctx)
	if err != nil {
		t.Errorf("unable to count schedules: %v", err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("CountSchedules() is %v, want %v", count, len(resources.Schedules))
	}
	methods["CountSchedules"] = true

	// count the schedules for a repo
	count, err = db.CountSchedulesForRepo(ctx, resources.Repos[0])
	if err != nil {
		t.Errorf("unable to count schedules for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("CountSchedulesForRepo() is %v, want %v", count, len(resources.Schedules))
	}
	methods["CountSchedulesForRepo"] = true

	// list the schedules
	list, err := db.ListSchedules(ctx)
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !cmp.Equal(list, resources.Schedules, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListSchedules() is %v, want %v", list, resources.Schedules)
	}
	methods["ListSchedules"] = true

	// list the active schedules
	list, err = db.ListActiveSchedules(ctx)
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !cmp.Equal(list, resources.Schedules, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListActiveSchedules() is %v, want %v", list, resources.Schedules)
	}
	methods["ListActiveSchedules"] = true

	// list the schedules for a repo
	list, count, err = db.ListSchedulesForRepo(ctx, resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to count schedules for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("ListSchedulesForRepo() is %v, want %v", count, len(resources.Schedules))
	}
	if !cmp.Equal(list, []*library.Schedule{resources.Schedules[1], resources.Schedules[0]}, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListSchedulesForRepo() is %v, want %v", list, []*library.Schedule{resources.Schedules[1], resources.Schedules[0]})
	}
	methods["ListSchedulesForRepo"] = true

	// lookup the schedules by name
	for _, schedule := range resources.Schedules {
		repo := resources.Repos[schedule.GetRepoID()-1]
		got, err := db.GetScheduleForRepo(ctx, repo, schedule.GetName())
		if err != nil {
			t.Errorf("unable to get schedule %d for repo %d: %v", schedule.GetID(), repo.GetID(), err)
		}
		if !cmp.Equal(got, schedule, CmpOptApproxUpdatedAt()) {
			t.Errorf("GetScheduleForRepo() is %v, want %v", got, schedule)
		}
	}
	methods["GetScheduleForRepo"] = true

	// update the schedules
	for _, schedule := range resources.Schedules {
		schedule.SetUpdatedAt(time.Now().UTC().Unix())
		got, err := db.UpdateSchedule(ctx, schedule, true)
		if err != nil {
			t.Errorf("unable to update schedule %d: %v", schedule.GetID(), err)
		}

		if !cmp.Equal(got, schedule, CmpOptApproxUpdatedAt()) {
			t.Errorf("GetSchedule() is %v, want %v", got, schedule)
		}
	}
	methods["UpdateSchedule"] = true
	methods["GetSchedule"] = true

	// delete the schedules
	for _, schedule := range resources.Schedules {
		err = db.DeleteSchedule(ctx, schedule)
		if err != nil {
			t.Errorf("unable to delete schedule %d: %v", schedule.GetID(), err)
		}
	}
	methods["DeleteSchedule"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for schedules", method)
		}
	}
}

func testSecrets(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for secrets
	methods := make(map[string]bool)
	// capture the element type of the secret interface
	element := reflect.TypeOf(new(secret.SecretInterface)).Elem()
	// iterate through all methods found in the secret interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for secrets
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the secrets
	for _, secret := range resources.Secrets {
		_, err := db.CreateSecret(secret)
		if err != nil {
			t.Errorf("unable to create secret %d: %v", secret.GetID(), err)
		}
	}
	methods["CreateSecret"] = true

	// count the secrets
	count, err := db.CountSecrets()
	if err != nil {
		t.Errorf("unable to count secrets: %v", err)
	}
	if int(count) != len(resources.Secrets) {
		t.Errorf("CountSecrets() is %v, want %v", count, len(resources.Secrets))
	}
	methods["CountSecrets"] = true

	for _, secret := range resources.Secrets {
		switch secret.GetType() {
		case constants.SecretOrg:
			// count the secrets for an org
			count, err = db.CountSecretsForOrg(secret.GetOrg(), nil)
			if err != nil {
				t.Errorf("unable to count secrets for org %s: %v", secret.GetOrg(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForOrg() is %v, want %v", count, 1)
			}
			methods["CountSecretsForOrg"] = true
		case constants.SecretRepo:
			// count the secrets for a repo
			count, err = db.CountSecretsForRepo(resources.Repos[0], nil)
			if err != nil {
				t.Errorf("unable to count secrets for repo %d: %v", resources.Repos[0].GetID(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForRepo() is %v, want %v", count, 1)
			}
			methods["CountSecretsForRepo"] = true
		case constants.SecretShared:
			// count the secrets for a team
			count, err = db.CountSecretsForTeam(secret.GetOrg(), secret.GetTeam(), nil)
			if err != nil {
				t.Errorf("unable to count secrets for team %s: %v", secret.GetTeam(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForTeam() is %v, want %v", count, 1)
			}
			methods["CountSecretsForTeam"] = true

			// count the secrets for a list of teams
			count, err = db.CountSecretsForTeams(secret.GetOrg(), []string{secret.GetTeam()}, nil)
			if err != nil {
				t.Errorf("unable to count secrets for teams %s: %v", []string{secret.GetTeam()}, err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForTeams() is %v, want %v", count, 1)
			}
			methods["CountSecretsForTeams"] = true
		default:
			t.Errorf("unsupported type %s for secret %d", secret.GetType(), secret.GetID())
		}
	}

	// list the secrets
	list, err := db.ListSecrets()
	if err != nil {
		t.Errorf("unable to list secrets: %v", err)
	}
	if !cmp.Equal(list, resources.Secrets, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListSecrets() is %v, want %v", list, resources.Secrets)
	}
	methods["ListSecrets"] = true

	for _, secret := range resources.Secrets {
		switch secret.GetType() {
		case constants.SecretOrg:
			// list the secrets for an org
			list, count, err = db.ListSecretsForOrg(secret.GetOrg(), nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for org %s: %v", secret.GetOrg(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForOrg() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*library.Secret{secret}) {
				t.Errorf("ListSecretsForOrg() is %v, want %v", list, []*library.Secret{secret})
			}
			methods["ListSecretsForOrg"] = true
		case constants.SecretRepo:
			// list the secrets for a repo
			list, count, err = db.ListSecretsForRepo(resources.Repos[0], nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for repo %d: %v", resources.Repos[0].GetID(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForRepo() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*library.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForRepo() is %v, want %v", list, []*library.Secret{secret})
			}
			methods["ListSecretsForRepo"] = true
		case constants.SecretShared:
			// list the secrets for a team
			list, count, err = db.ListSecretsForTeam(secret.GetOrg(), secret.GetTeam(), nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for team %s: %v", secret.GetTeam(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForTeam() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*library.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForTeam() is %v, want %v", list, []*library.Secret{secret})
			}
			methods["ListSecretsForTeam"] = true

			// list the secrets for a list of teams
			list, count, err = db.ListSecretsForTeams(secret.GetOrg(), []string{secret.GetTeam()}, nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for teams %s: %v", []string{secret.GetTeam()}, err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForTeams() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*library.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForTeams() is %v, want %v", list, []*library.Secret{secret})
			}
			methods["ListSecretsForTeams"] = true
		default:
			t.Errorf("unsupported type %s for secret %d", secret.GetType(), secret.GetID())
		}
	}

	for _, secret := range resources.Secrets {
		switch secret.GetType() {
		case constants.SecretOrg:
			// lookup the secret by org
			got, err := db.GetSecretForOrg(secret.GetOrg(), secret.GetName())
			if err != nil {
				t.Errorf("unable to get secret %d for org %s: %v", secret.GetID(), secret.GetOrg(), err)
			}
			if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
				t.Errorf("GetSecretForOrg() is %v, want %v", got, secret)
			}
			methods["GetSecretForOrg"] = true
		case constants.SecretRepo:
			// lookup the secret by repo
			got, err := db.GetSecretForRepo(secret.GetName(), resources.Repos[0])
			if err != nil {
				t.Errorf("unable to get secret %d for repo %d: %v", secret.GetID(), resources.Repos[0].GetID(), err)
			}
			if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
				t.Errorf("GetSecretForRepo() is %v, want %v", got, secret)
			}
			methods["GetSecretForRepo"] = true
		case constants.SecretShared:
			// lookup the secret by team
			got, err := db.GetSecretForTeam(secret.GetOrg(), secret.GetTeam(), secret.GetName())
			if err != nil {
				t.Errorf("unable to get secret %d for team %s: %v", secret.GetID(), secret.GetTeam(), err)
			}
			if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
				t.Errorf("GetSecretForTeam() is %v, want %v", got, secret)
			}
			methods["GetSecretForTeam"] = true
		default:
			t.Errorf("unsupported type %s for secret %d", secret.GetType(), secret.GetID())
		}
	}

	// update the secrets
	for _, secret := range resources.Secrets {
		secret.SetUpdatedAt(time.Now().UTC().Unix())
		got, err := db.UpdateSecret(secret)
		if err != nil {
			t.Errorf("unable to update secret %d: %v", secret.GetID(), err)
		}

		if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
			t.Errorf("GetSecret() is %v, want %v", got, secret)
		}
	}
	methods["UpdateSecret"] = true
	methods["GetSecret"] = true

	// delete the secrets
	for _, secret := range resources.Secrets {
		err = db.DeleteSecret(secret)
		if err != nil {
			t.Errorf("unable to delete secret %d: %v", secret.GetID(), err)
		}
	}
	methods["DeleteSecret"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for secrets", method)
		}
	}
}

func testServices(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for services
	methods := make(map[string]bool)
	// capture the element type of the service interface
	element := reflect.TypeOf(new(service.ServiceInterface)).Elem()
	// iterate through all methods found in the service interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for services
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the services
	for _, service := range resources.Services {
		_, err := db.CreateService(service)
		if err != nil {
			t.Errorf("unable to create service %d: %v", service.GetID(), err)
		}
	}
	methods["CreateService"] = true

	// count the services
	count, err := db.CountServices()
	if err != nil {
		t.Errorf("unable to count services: %v", err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServices() is %v, want %v", count, len(resources.Services))
	}
	methods["CountServices"] = true

	// count the services for a build
	count, err = db.CountServicesForBuild(resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServicesForBuild() is %v, want %v", count, len(resources.Services))
	}
	methods["CountServicesForBuild"] = true

	// list the services
	list, err := db.ListServices()
	if err != nil {
		t.Errorf("unable to list services: %v", err)
	}
	if !cmp.Equal(list, resources.Services) {
		t.Errorf("ListServices() is %v, want %v", list, resources.Services)
	}
	methods["ListServices"] = true

	// list the services for a build
	list, count, err = db.ListServicesForBuild(resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !cmp.Equal(list, []*library.Service{resources.Services[1], resources.Services[0]}) {
		t.Errorf("ListServicesForBuild() is %v, want %v", list, []*library.Service{resources.Services[1], resources.Services[0]})
	}
	if int(count) != len(resources.Services) {
		t.Errorf("ListServicesForBuild() is %v, want %v", count, len(resources.Services))
	}
	methods["ListServicesForBuild"] = true

	expected := map[string]float64{
		"#init":                  1,
		"target/vela-git:v0.3.0": 1,
	}
	images, err := db.ListServiceImageCount()
	if err != nil {
		t.Errorf("unable to list service image count: %v", err)
	}
	if !cmp.Equal(images, expected) {
		t.Errorf("ListServiceImageCount() is %v, want %v", images, expected)
	}
	methods["ListServiceImageCount"] = true

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
	if !cmp.Equal(statuses, expected) {
		t.Errorf("ListServiceStatusCount() is %v, want %v", statuses, expected)
	}
	methods["ListServiceStatusCount"] = true

	// lookup the services by name
	for _, service := range resources.Services {
		build := resources.Builds[service.GetBuildID()-1]
		got, err := db.GetServiceForBuild(build, service.GetNumber())
		if err != nil {
			t.Errorf("unable to get service %d for build %d: %v", service.GetID(), build.GetID(), err)
		}
		if !cmp.Equal(got, service) {
			t.Errorf("GetServiceForBuild() is %v, want %v", got, service)
		}
	}
	methods["GetServiceForBuild"] = true

	// clean the services
	count, err = db.CleanServices("integration testing", 1563474090)
	if err != nil {
		t.Errorf("unable to clean services: %v", err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CleanServices() is %v, want %v", count, len(resources.Services))
	}
	methods["CleanServices"] = true

	// update the services
	for _, service := range resources.Services {
		service.SetStatus("success")
		got, err := db.UpdateService(service)
		if err != nil {
			t.Errorf("unable to update service %d: %v", service.GetID(), err)
		}

		if !cmp.Equal(got, service) {
			t.Errorf("UpdateService() is %v, want %v", got, service)
		}
	}
	methods["UpdateService"] = true
	methods["GetService"] = true

	// delete the services
	for _, service := range resources.Services {
		err = db.DeleteService(service)
		if err != nil {
			t.Errorf("unable to delete service %d: %v", service.GetID(), err)
		}
	}
	methods["DeleteService"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for services", method)
		}
	}
}

func testSteps(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for steps
	methods := make(map[string]bool)
	// capture the element type of the step interface
	element := reflect.TypeOf(new(step.StepInterface)).Elem()
	// iterate through all methods found in the step interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for steps
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the steps
	for _, step := range resources.Steps {
		_, err := db.CreateStep(step)
		if err != nil {
			t.Errorf("unable to create step %d: %v", step.GetID(), err)
		}
	}
	methods["CreateStep"] = true

	// count the steps
	count, err := db.CountSteps()
	if err != nil {
		t.Errorf("unable to count steps: %v", err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountSteps() is %v, want %v", count, len(resources.Steps))
	}
	methods["CountSteps"] = true

	// count the steps for a build
	count, err = db.CountStepsForBuild(resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountStepsForBuild() is %v, want %v", count, len(resources.Steps))
	}
	methods["CountStepsForBuild"] = true

	// list the steps
	list, err := db.ListSteps()
	if err != nil {
		t.Errorf("unable to list steps: %v", err)
	}
	if !cmp.Equal(list, resources.Steps) {
		t.Errorf("ListSteps() is %v, want %v", list, resources.Steps)
	}
	methods["ListSteps"] = true

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !cmp.Equal(list, []*library.Step{resources.Steps[1], resources.Steps[0]}) {
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*library.Step{resources.Steps[1], resources.Steps[0]})
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("ListStepsForBuild() is %v, want %v", count, len(resources.Steps))
	}
	methods["ListStepsForBuild"] = true

	expected := map[string]float64{
		"#init":                  1,
		"target/vela-git:v0.3.0": 1,
	}
	images, err := db.ListStepImageCount()
	if err != nil {
		t.Errorf("unable to list step image count: %v", err)
	}
	if !cmp.Equal(images, expected) {
		t.Errorf("ListStepImageCount() is %v, want %v", images, expected)
	}
	methods["ListStepImageCount"] = true

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
	if !cmp.Equal(statuses, expected) {
		t.Errorf("ListStepStatusCount() is %v, want %v", statuses, expected)
	}
	methods["ListStepStatusCount"] = true

	// lookup the steps by name
	for _, step := range resources.Steps {
		build := resources.Builds[step.GetBuildID()-1]
		got, err := db.GetStepForBuild(build, step.GetNumber())
		if err != nil {
			t.Errorf("unable to get step %d for build %d: %v", step.GetID(), build.GetID(), err)
		}
		if !cmp.Equal(got, step) {
			t.Errorf("GetStepForBuild() is %v, want %v", got, step)
		}
	}
	methods["GetStepForBuild"] = true

	// clean the steps
	count, err = db.CleanSteps("integration testing", 1563474090)
	if err != nil {
		t.Errorf("unable to clean steps: %v", err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CleanSteps() is %v, want %v", count, len(resources.Steps))
	}
	methods["CleanSteps"] = true

	// update the steps
	for _, step := range resources.Steps {
		step.SetStatus("success")
		got, err := db.UpdateStep(step)
		if err != nil {
			t.Errorf("unable to update step %d: %v", step.GetID(), err)
		}

		if !cmp.Equal(got, step) {
			t.Errorf("GetStep() is %v, want %v", got, step)
		}
	}
	methods["UpdateStep"] = true
	methods["GetStep"] = true

	// delete the steps
	for _, step := range resources.Steps {
		err = db.DeleteStep(step)
		if err != nil {
			t.Errorf("unable to delete step %d: %v", step.GetID(), err)
		}
	}
	methods["DeleteStep"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for steps", method)
		}
	}
}

func testUsers(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for users
	methods := make(map[string]bool)
	// capture the element type of the user interface
	element := reflect.TypeOf(new(user.UserInterface)).Elem()
	// iterate through all methods found in the user interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for users
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	userOne := new(library.User)
	userOne.SetID(1)
	userOne.SetName("octocat")
	userOne.SetToken("")
	userOne.SetRefreshToken("")
	userOne.SetHash("")
	userOne.SetFavorites(nil)
	userOne.SetActive(false)
	userOne.SetAdmin(false)

	userTwo := new(library.User)
	userTwo.SetID(2)
	userTwo.SetName("octokitty")
	userTwo.SetToken("")
	userTwo.SetRefreshToken("")
	userTwo.SetHash("")
	userTwo.SetFavorites(nil)
	userTwo.SetActive(false)
	userTwo.SetAdmin(false)

	liteUsers := []*library.User{userOne, userTwo}

	// create the users
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}
	methods["CreateUser"] = true

	// count the users
	count, err := db.CountUsers(context.TODO())
	if err != nil {
		t.Errorf("unable to count users: %v", err)
	}
	if int(count) != len(resources.Users) {
		t.Errorf("CountUsers() is %v, want %v", count, len(resources.Users))
	}
	methods["CountUsers"] = true

	// list the users
	list, err := db.ListUsers(context.TODO())
	if err != nil {
		t.Errorf("unable to list users: %v", err)
	}
	if !cmp.Equal(list, resources.Users) {
		t.Errorf("ListUsers() is %v, want %v", list, resources.Users)
	}
	methods["ListUsers"] = true

	// lite list the users
	list, count, err = db.ListLiteUsers(context.TODO(), 1, 10)
	if err != nil {
		t.Errorf("unable to list lite users: %v", err)
	}
	if !cmp.Equal(list, liteUsers) {
		t.Errorf("ListLiteUsers() is %v, want %v", list, liteUsers)
	}
	if int(count) != len(liteUsers) {
		t.Errorf("ListLiteUsers() is %v, want %v", count, len(liteUsers))
	}
	methods["ListLiteUsers"] = true

	// lookup the users by name
	for _, user := range resources.Users {
		got, err := db.GetUserForName(context.TODO(), user.GetName())
		if err != nil {
			t.Errorf("unable to get user %d by name: %v", user.GetID(), err)
		}
		if !cmp.Equal(got, user) {
			t.Errorf("GetUserForName() is %v, want %v", got, user)
		}
	}
	methods["GetUserForName"] = true

	// update the users
	for _, user := range resources.Users {
		user.SetActive(false)
		got, err := db.UpdateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to update user %d: %v", user.GetID(), err)
		}

		if !cmp.Equal(got, user) {
			t.Errorf("GetUser() is %v, want %v", got, user)
		}
	}
	methods["UpdateUser"] = true
	methods["GetUser"] = true

	// delete the users
	for _, user := range resources.Users {
		err = db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}
	methods["DeleteUser"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for users", method)
		}
	}
}

func testWorkers(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for workers
	methods := make(map[string]bool)
	// capture the element type of the worker interface
	element := reflect.TypeOf(new(worker.WorkerInterface)).Elem()
	// iterate through all methods found in the worker interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for workers
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the workers
	for _, worker := range resources.Workers {
		_, err := db.CreateWorker(context.TODO(), worker)
		if err != nil {
			t.Errorf("unable to create worker %d: %v", worker.GetID(), err)
		}
	}
	methods["CreateWorker"] = true

	// count the workers
	count, err := db.CountWorkers(context.TODO())
	if err != nil {
		t.Errorf("unable to count workers: %v", err)
	}
	if int(count) != len(resources.Workers) {
		t.Errorf("CountWorkers() is %v, want %v", count, len(resources.Workers))
	}
	methods["CountWorkers"] = true

	// list the workers
	list, err := db.ListWorkers(context.TODO())
	if err != nil {
		t.Errorf("unable to list workers: %v", err)
	}
	if !cmp.Equal(list, resources.Workers) {
		t.Errorf("ListWorkers() is %v, want %v", list, resources.Workers)
	}
	methods["ListWorkers"] = true

	// lookup the workers by hostname
	for _, worker := range resources.Workers {
		got, err := db.GetWorkerForHostname(context.TODO(), worker.GetHostname())
		if err != nil {
			t.Errorf("unable to get worker %d by hostname: %v", worker.GetID(), err)
		}
		if !cmp.Equal(got, worker) {
			t.Errorf("GetWorkerForHostname() is %v, want %v", got, worker)
		}
	}
	methods["GetWorkerForHostname"] = true

	// update the workers
	for _, worker := range resources.Workers {
		worker.SetActive(false)
		got, err := db.UpdateWorker(context.TODO(), worker)
		if err != nil {
			t.Errorf("unable to update worker %d: %v", worker.GetID(), err)
		}

		if !cmp.Equal(got, worker) {
			t.Errorf("GetWorker() is %v, want %v", got, worker)
		}
	}
	methods["UpdateWorker"] = true
	methods["GetWorker"] = true

	// delete the workers
	for _, worker := range resources.Workers {
		err = db.DeleteWorker(context.TODO(), worker)
		if err != nil {
			t.Errorf("unable to delete worker %d: %v", worker.GetID(), err)
		}
	}
	methods["DeleteWorker"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for workers", method)
		}
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

	executableOne := new(library.BuildExecutable)
	executableOne.SetID(1)
	executableOne.SetBuildID(1)
	executableOne.SetData([]byte("foo"))

	executableTwo := new(library.BuildExecutable)
	executableTwo.SetID(2)
	executableTwo.SetBuildID(2)
	executableTwo.SetData([]byte("foo"))

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
	scheduleOne.SetBranch("main")

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
	scheduleTwo.SetBranch("main")

	secretOrg := new(library.Secret)
	secretOrg.SetID(1)
	secretOrg.SetOrg("github")
	secretOrg.SetRepo("*")
	secretOrg.SetTeam("")
	secretOrg.SetName("foo")
	secretOrg.SetValue("bar")
	secretOrg.SetType("org")
	secretOrg.SetImages([]string{"alpine"})
	secretOrg.SetEvents([]string{"push", "tag", "deployment"})
	secretOrg.SetAllowCommand(true)
	secretOrg.SetCreatedAt(time.Now().UTC().Unix())
	secretOrg.SetCreatedBy("octocat")
	secretOrg.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretOrg.SetUpdatedBy("octokitty")

	secretRepo := new(library.Secret)
	secretRepo.SetID(2)
	secretRepo.SetOrg("github")
	secretRepo.SetRepo("octocat")
	secretRepo.SetTeam("")
	secretRepo.SetName("foo")
	secretRepo.SetValue("bar")
	secretRepo.SetType("repo")
	secretRepo.SetImages([]string{"alpine"})
	secretRepo.SetEvents([]string{"push", "tag", "deployment"})
	secretRepo.SetAllowCommand(true)
	secretRepo.SetCreatedAt(time.Now().UTC().Unix())
	secretRepo.SetCreatedBy("octocat")
	secretRepo.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretRepo.SetUpdatedBy("octokitty")

	secretShared := new(library.Secret)
	secretShared.SetID(3)
	secretShared.SetOrg("github")
	secretShared.SetRepo("")
	secretShared.SetTeam("octocat")
	secretShared.SetName("foo")
	secretShared.SetValue("bar")
	secretShared.SetType("shared")
	secretShared.SetImages([]string{"alpine"})
	secretShared.SetEvents([]string{"push", "tag", "deployment"})
	secretShared.SetAllowCommand(true)
	secretShared.SetCreatedAt(time.Now().UTC().Unix())
	secretShared.SetCreatedBy("octocat")
	secretShared.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretShared.SetUpdatedBy("octokitty")

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
		Executables: []*library.BuildExecutable{executableOne, executableTwo},
		Hooks:       []*library.Hook{hookOne, hookTwo},
		Logs:        []*library.Log{logServiceOne, logServiceTwo, logStepOne, logStepTwo},
		Pipelines:   []*library.Pipeline{pipelineOne, pipelineTwo},
		Repos:       []*library.Repo{repoOne, repoTwo},
		Schedules:   []*library.Schedule{scheduleOne, scheduleTwo},
		Secrets:     []*library.Secret{secretOrg, secretRepo, secretShared},
		Services:    []*library.Service{serviceOne, serviceTwo},
		Steps:       []*library.Step{stepOne, stepTwo},
		Users:       []*library.User{userOne, userTwo},
		Workers:     []*library.Worker{workerOne, workerTwo},
	}
}

// CmpOptApproxUpdatedAt is a custom comparator for cmp.Equal
// to reduce flakiness in tests when comparing structs with UpdatedAt field.
func CmpOptApproxUpdatedAt() cmp.Option {
	// Custom Comparer
	//
	// https://pkg.go.dev/github.com/google/go-cmp/cmp#Comparer
	cmpApproximateUnixTime := cmp.Comparer(func(x, y *int64) bool {
		if x == nil && y == nil {
			return true
		}

		if x != nil && y == nil || y != nil && x == nil {
			return false
		}

		// make sure we subtract smaller value from larger one
		if *y < *x {
			*x, *y = *y, *x
		}

		// is it less than 5 seconds? consider the values equal
		return *y-*x < 5
	})

	// only apply to structs with UpdatedAt field
	//
	// https://pkg.go.dev/github.com/google/go-cmp/cmp#FilterPath
	return cmp.FilterPath(
		func(p cmp.Path) bool {
			return p.Last().String() == "UpdatedAt" || p.Last().String() == ".UpdatedAt"
		},
		cmpApproximateUnixTime)
}
