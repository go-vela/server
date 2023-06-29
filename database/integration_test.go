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

			t.Run("test_steps", func(t *testing.T) {
				testSteps(t, db, []*library.Build{buildOne, buildTwo})
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

func testSteps(t *testing.T, db Interface, builds []*library.Build) {
	// used to track the number of methods we call for steps
	//
	// we start at 2 for creating the table and indexes for steps
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Step)
	one.SetID(1)
	one.SetBuildID(1)
	one.SetRepoID(1)
	one.SetNumber(1)
	one.SetName("init")
	one.SetImage("#init")
	one.SetStage("init")
	one.SetStatus("running")
	one.SetError("")
	one.SetExitCode(0)
	one.SetCreated(1563474076)
	one.SetStarted(1563474078)
	one.SetFinished(1563474079)
	one.SetHost("example.company.com")
	one.SetRuntime("docker")
	one.SetDistribution("linux")

	two := new(library.Step)
	two.SetID(2)
	two.SetBuildID(1)
	two.SetRepoID(1)
	two.SetNumber(2)
	two.SetName("clone")
	two.SetImage("target/vela-git:v0.3.0")
	two.SetStage("init")
	two.SetStatus("pending")
	two.SetError("")
	two.SetExitCode(0)
	two.SetCreated(1563474086)
	two.SetStarted(1563474088)
	two.SetFinished(1563474089)
	two.SetHost("example.company.com")
	two.SetRuntime("docker")
	two.SetDistribution("linux")

	steps := []*library.Step{one, two}

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
		t.Errorf("CountSteps() is %v, want 2", count)
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
		pretty.Ldiff(t, list, steps)
		t.Errorf("ListSteps() is %v, want %v", list, steps)
	}
	counter++

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Step{two, one}) {
		pretty.Ldiff(t, list, steps)
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*library.Step{two, one})
	}
	if int(count) != len(steps) {
		t.Errorf("ListStepsForBuild() is %v, want %v", count, len(steps))
	}
	counter++

	images, err := db.ListStepImageCount()
	if err != nil {
		t.Errorf("unable to list step image count: %v",, err)
	}
	if len(images) != len(steps) {
		t.Errorf("ListStepImageCount() is %v, want %v", len(images), len(steps))
	}
	counter++

	statuses, err := db.ListStepStatusCount()
	if err != nil {
		t.Errorf("unable to list step status count: %v", err)
	}
	if len(statuses) != len(steps) {
		t.Errorf("ListStepStatusCount() is %v, want %v", len(images), len(steps))
	}
	counter++

	// lookup the steps by name
	for _, step := range steps {
		got, err := db.GetStepForBuild(builds[0], step.GetNumber())
		if err != nil {
			t.Errorf("unable to get step %s for build %d: %v", step.GetName(), builds[0].GetID(), err)
		}
		if !reflect.DeepEqual(got, step) {
			pretty.Ldiff(t, got, step)
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
			pretty.Ldiff(t, got, step)
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