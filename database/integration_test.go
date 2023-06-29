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

			t.Run("test_build", func(t *testing.T) {
				testBuilds(t, db, []*library.Build{buildOne, buildTwo}, []*library.Repo{repoOne, repoTwo})
			})

			t.Run("test_services", func(t *testing.T) {
				testServices(t, db, []*library.Build{buildOne, buildTwo})
			})

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
		t.Errorf("CountBuilds() is %v, want 2", count)
	}
	counter++

	// list the builds
	list, err := db.ListBuilds()
	if err != nil {
		t.Errorf("unable to list builds: %v", err)
	}
	if !reflect.DeepEqual(list, builds) {
		t.Errorf("ListBuilds() is %v, want %v", list, builds)
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

func testServices(t *testing.T, db Interface, builds []*library.Build) {
	// used to track the number of methods we call for services
	//
	// we start at 2 for creating the table and indexes for services
	// since those are already called when the database engine starts
	counter := 2

	one := new(library.Service)
	one.SetID(1)
	one.SetBuildID(1)
	one.SetRepoID(1)
	one.SetNumber(1)
	one.SetName("init")
	one.SetImage("#init")
	one.SetStatus("running")
	one.SetError("")
	one.SetExitCode(0)
	one.SetCreated(1563474076)
	one.SetStarted(1563474078)
	one.SetFinished(1563474079)
	one.SetHost("example.company.com")
	one.SetRuntime("docker")
	one.SetDistribution("linux")

	two := new(library.Service)
	two.SetID(2)
	two.SetBuildID(1)
	two.SetRepoID(1)
	two.SetNumber(2)
	two.SetName("clone")
	two.SetImage("target/vela-git:v0.3.0")
	two.SetStatus("pending")
	two.SetError("")
	two.SetExitCode(0)
	two.SetCreated(1563474086)
	two.SetStarted(1563474088)
	two.SetFinished(1563474089)
	two.SetHost("example.company.com")
	two.SetRuntime("docker")
	two.SetDistribution("linux")

	services := []*library.Service{one, two}

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
		t.Errorf("CountServices() is %v, want 2", count)
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
	if !reflect.DeepEqual(list, []*library.Service{two, one}) {
		t.Errorf("ListServicesForBuild() is %v, want %v", list, []*library.Service{two, one})
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
		t.Errorf("ListSteps() is %v, want %v", list, steps)
	}
	counter++

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", builds[0].GetID(), err)
	}
	if !reflect.DeepEqual(list, []*library.Step{two, one}) {
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*library.Step{two, one})
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
	liteOne.SetFavorites(nil)
	liteOne.SetActive(false)
	liteOne.SetAdmin(false)

	liteTwo := new(library.User)
	liteTwo.SetID(2)
	liteTwo.SetName("octokitty")
	liteTwo.SetToken("")
	liteTwo.SetRefreshToken("")
	liteTwo.SetHash("")
	liteTwo.SetFavorites(nil)
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
