// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/adhocore/gronx"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/v2/jwk"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/build"
	"github.com/go-vela/server/database/dashboard"
	"github.com/go-vela/server/database/deployment"
	"github.com/go-vela/server/database/executable"
	"github.com/go-vela/server/database/hook"
	dbJWK "github.com/go-vela/server/database/jwk"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/schedule"
	"github.com/go-vela/server/database/secret"
	"github.com/go-vela/server/database/service"
	dbSettings "github.com/go-vela/server/database/settings"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/server/tracing"
)

// Resources represents the object containing test resources.
type Resources struct {
	Builds      []*api.Build
	Dashboards  []*api.Dashboard
	Deployments []*api.Deployment
	Executables []*api.BuildExecutable
	Hooks       []*api.Hook
	JWKs        jwk.Set
	Logs        []*api.Log
	Pipelines   []*api.Pipeline
	Repos       []*api.Repo
	Schedules   []*api.Schedule
	Secrets     []*api.Secret
	Services    []*api.Service
	Steps       []*api.Step
	Users       []*api.User
	Workers     []*api.Worker
	Platform    []*settings.Platform
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
				WithTracing(&tracing.Client{Config: tracing.Config{EnableTracing: false}}),
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

			t.Run("test_dashboards", func(t *testing.T) { testDashboards(t, db, resources) })

			t.Run("test_deployments", func(t *testing.T) { testDeployments(t, db, resources) })

			t.Run("test_executables", func(t *testing.T) { testExecutables(t, db, resources) })

			t.Run("test_hooks", func(t *testing.T) { testHooks(t, db, resources) })

			t.Run("test_jwks", func(t *testing.T) { testJWKs(t, db, resources) })

			t.Run("test_logs", func(t *testing.T) { testLogs(t, db, resources) })

			t.Run("test_pipelines", func(t *testing.T) { testPipelines(t, db, resources) })

			t.Run("test_repos", func(t *testing.T) { testRepos(t, db, resources) })

			t.Run("test_schedules", func(t *testing.T) { testSchedules(t, db, resources) })

			t.Run("test_secrets", func(t *testing.T) { testSecrets(t, db, resources) })

			t.Run("test_services", func(t *testing.T) { testServices(t, db, resources) })

			t.Run("test_steps", func(t *testing.T) { testSteps(t, db, resources) })

			t.Run("test_users", func(t *testing.T) { testUsers(t, db, resources) })

			t.Run("test_workers", func(t *testing.T) { testWorkers(t, db, resources) })

			t.Run("test_settings", func(t *testing.T) { testSettings(t, db, resources) })

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

	// create the users for build related functions (owners of repos)
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}

	// create the repos for build related functions
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	buildOne := new(api.QueueBuild)
	buildOne.SetCreated(1563474076)
	buildOne.SetFullName("github/octocat")
	buildOne.SetNumber(1)
	buildOne.SetStatus("running")

	buildTwo := new(api.QueueBuild)
	buildTwo.SetCreated(1563474076)
	buildTwo.SetFullName("github/octocat")
	buildTwo.SetNumber(2)
	buildTwo.SetStatus("running")

	queueBuilds := []*api.QueueBuild{buildOne, buildTwo}

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
	count, err = db.CountBuildsForRepo(context.TODO(), resources.Repos[0], nil, time.Now().Unix(), 0)
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
	if diff := cmp.Diff(resources.Builds, list); diff != "" {
		t.Errorf("ListBuilds() mismatch (-want +got):\n%s", diff)
	}
	methods["ListBuilds"] = true

	// list the builds for an org
	list, count, err = db.ListBuildsForOrg(context.TODO(), resources.Repos[0].GetOrg(), nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list builds for org %s: %v", resources.Repos[0].GetOrg(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListBuildsForOrg() is %v, want %v", count, len(resources.Builds))
	}
	if diff := cmp.Diff(resources.Builds, list); diff != "" {
		t.Errorf("ListBuildsForOrg() mismatch (-want +got):\n%s", diff)
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
	if diff := cmp.Diff([]*api.Build{resources.Builds[1], resources.Builds[0]}, list); diff != "" {
		t.Errorf("ListBuildsForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["ListBuildsForRepo"] = true

	list, err = db.ListBuildsForDashboardRepo(context.TODO(), resources.Repos[0], []string{"main"}, []string{"push"})
	if err != nil {
		t.Errorf("unable to list build for dashboard repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if len(list) != 1 {
		t.Errorf("Number of results for ListBuildsForDashboardRepo() is %v, want %v", len(list), 1)
	}

	if diff := cmp.Diff([]*api.Build{resources.Hooks[0].GetBuild()}, list); diff != "" {
		t.Errorf("ListBuildsForDashboardRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["ListBuildsForDashboardRepo"] = true

	// list the pending / running builds for a repo
	list, err = db.ListPendingAndRunningBuildsForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to list pending and running builds for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("ListPendingAndRunningBuildsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	if diff := cmp.Diff([]*api.Build{resources.Builds[0], resources.Builds[1]}, list); diff != "" {
		t.Errorf("ListPendingAndRunningBuildsForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["ListPendingAndRunningBuildsForRepo"] = true

	// list the pending and running builds
	queueList, err := db.ListPendingAndRunningBuilds(context.TODO(), "0")
	if err != nil {
		t.Errorf("unable to list pending and running builds: %v", err)
	}
	if diff := cmp.Diff(queueBuilds, queueList); diff != "" {
		t.Errorf("ListPendingAndRunningBuilds() mismatch (-want +got):\n%s", diff)
	}
	methods["ListPendingAndRunningBuilds"] = true

	// lookup the last build by repo
	got, err := db.LastBuildForRepo(context.TODO(), resources.Repos[0], "main")
	if err != nil {
		t.Errorf("unable to get last build for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if diff := cmp.Diff(resources.Builds[1], got); diff != "" {
		t.Errorf("LastBuildForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["LastBuildForRepo"] = true

	// lookup the builds by repo and number
	for _, build := range resources.Builds {
		repo := resources.Repos[build.GetRepo().GetID()-1]
		got, err = db.GetBuildForRepo(context.TODO(), repo, build.GetNumber())
		if err != nil {
			t.Errorf("unable to get build %d for repo %d: %v", build.GetID(), repo.GetID(), err)
		}
		if diff := cmp.Diff(build, got); diff != "" {
			t.Errorf("GetBuildForRepo() mismatch (-want +got):\n%s", diff)
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
	for i, build := range resources.Builds {
		prevStatus := build.GetStatus()

		build.SetStatus("pending approval")
		_, err = db.UpdateBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to update build %d: %v", build.GetID(), err)
		}

		// lookup the build by ID
		got, err = db.GetBuild(context.TODO(), build.GetID())
		if err != nil {
			t.Errorf("unable to get build %d by ID: %v", build.GetID(), err)
		}
		if diff := cmp.Diff(build, got); diff != "" {
			t.Errorf("GetBuild() mismatch (-want +got):\n%s", diff)
		}

		if i == 1 {
			pABuilds, err := db.ListPendingApprovalBuilds(context.TODO(), "1663474076")
			if err != nil {
				t.Errorf("unable to list pending approval builds: %v", err)
			}

			if len(pABuilds) != 2 {
				t.Errorf("ListPendingApprovalBuilds() is %v, want %v", len(pABuilds), 2)
			}
		}

		build.SetStatus(prevStatus)
	}
	methods["UpdateBuild"] = true
	methods["GetBuild"] = true
	methods["ListPendingApprovalBuilds"] = true

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

	// delete the users for the build related functions
	for _, user := range resources.Users {
		err = db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for builds", method)
		}
	}
}

func testDashboards(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for schedules
	methods := make(map[string]bool)
	// capture the element type of the schedule interface
	element := reflect.TypeOf(new(dashboard.DashboardInterface)).Elem()
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

	// create the dashboard
	for _, dashboard := range resources.Dashboards {
		_, err := db.CreateDashboard(ctx, dashboard)
		if err != nil {
			t.Errorf("unable to create dashboard %s: %v", dashboard.GetID(), err)
		}
	}
	methods["CreateDashboard"] = true

	// lookup the dashboards by ID
	for _, dashboard := range resources.Dashboards {
		got, err := db.GetDashboard(ctx, dashboard.GetID())
		if err != nil {
			t.Errorf("unable to get dashboard %s: %v", dashboard.GetID(), err)
		}

		// JSON tags of `-` prevent unmarshaling of tokens, but they are sanitized anyway
		cmpAdmins := []*api.User{}
		for _, admin := range got.GetAdmins() {
			cmpAdmins = append(cmpAdmins, admin.Crop())
		}
		got.SetAdmins(cmpAdmins)

		if !cmp.Equal(got, dashboard, CmpOptApproxUpdatedAt()) {
			t.Errorf("GetDashboard() is %v, want %v", got, dashboard)
		}
	}
	methods["GetDashboard"] = true

	// update the dashboards
	for _, dashboard := range resources.Dashboards {
		dashboard.SetUpdatedAt(time.Now().UTC().Unix())
		got, err := db.UpdateDashboard(ctx, dashboard)
		if err != nil {
			t.Errorf("unable to update dashboard %s: %v", dashboard.GetID(), err)
		}

		// JSON marshaling does not include comparing token due to `-` struct tag
		cmpAdmins := got.GetAdmins()
		for i, admin := range cmpAdmins {
			admin.SetToken(resources.Users[i].GetToken())
		}

		if diff := cmp.Diff(dashboard, got, CmpOptApproxUpdatedAt()); diff != "" {
			t.Errorf("UpdateDashboard() mismatch (-want +got):\n%s", diff)
		}
	}
	methods["UpdateDashboard"] = true

	// delete the schedules
	for _, dashboard := range resources.Dashboards {
		err := db.DeleteDashboard(ctx, dashboard)
		if err != nil {
			t.Errorf("unable to delete dashboard %s: %v", dashboard.GetID(), err)
		}
	}
	methods["DeleteDashboard"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for dashboards", method)
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
		if diff := cmp.Diff(executable, got); diff != "" {
			t.Errorf("PopBuildExecutable() mismatch (-want +got):\n%s", diff)
		}
	}
	methods["PopBuildExecutable"] = true

	prevBuildStatus := resources.Builds[0].GetStatus()

	resources.Builds[0].SetStatus(constants.StatusError)

	_, err := db.UpdateBuild(context.TODO(), resources.Builds[0])
	if err != nil {
		t.Errorf("unable to update build for clean executables test")
	}

	// reset build status for other tests
	resources.Builds[0].SetStatus(prevBuildStatus)

	err = db.CreateBuildExecutable(context.TODO(), resources.Executables[0])
	if err != nil {
		t.Errorf("unable to create executable %d: %v", resources.Executables[0].GetID(), err)
	}

	count, err := db.CleanBuildExecutables(context.TODO())
	if err != nil {
		t.Errorf("unable to clean executable %d: %v", resources.Executables[0].GetID(), err)
	}

	if count != 1 {
		t.Errorf("CleanBuildExecutables should have affected 1 row, affected %d", count)
	}

	_, err = db.PopBuildExecutable(context.TODO(), resources.Builds[0].GetID())
	if err == nil {
		t.Errorf("build executable not cleaned")
	}

	methods["CleanBuildExecutables"] = true

	// remove build used for clean executables test
	err = db.DeleteBuild(context.TODO(), resources.Builds[0])
	if err != nil {
		t.Errorf("unable to delete build %d: %v", resources.Builds[0].GetID(), err)
	}

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for pipelines", method)
		}
	}
}

func testDeployments(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for deployments
	methods := make(map[string]bool)
	// capture the element type of the deployment interface
	element := reflect.TypeOf(new(deployment.DeploymentInterface)).Elem()
	// iterate through all methods found in the deployment interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for deployments
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the users for deployment related functions (owners of repos)
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}

	// create the repos for deployment related functions
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	// create the builds for deployment related functions
	for _, build := range resources.Builds {
		_, err := db.CreateBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to create build %d: %v", build.GetID(), err)
		}
	}

	// create the deployments
	for _, deployment := range resources.Deployments {
		_, err := db.CreateDeployment(context.TODO(), deployment)
		if err != nil {
			t.Errorf("unable to create deployment %d: %v", deployment.GetID(), err)
		}
	}
	methods["CreateDeployment"] = true

	// count the deployments
	count, err := db.CountDeployments(context.TODO())
	if err != nil {
		t.Errorf("unable to count deployment: %v", err)
	}
	if int(count) != len(resources.Deployments) {
		t.Errorf("CountDeployments() is %v, want %v", count, len(resources.Deployments))
	}
	methods["CountDeployments"] = true

	// count the deployments for a repo
	count, err = db.CountDeploymentsForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to count deployments for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Builds) {
		t.Errorf("CountDeploymentsForRepo() is %v, want %v", count, len(resources.Builds))
	}
	methods["CountDeploymentsForRepo"] = true

	// list the deployments
	list, err := db.ListDeployments(context.TODO())
	if err != nil {
		t.Errorf("unable to list deployments: %v", err)
	}
	if diff := cmp.Diff(resources.Deployments, list); diff != "" {
		t.Errorf("ListDeployments() mismatch (-want +got):\n%s", diff)
	}
	methods["ListDeployments"] = true

	// list the deployments for a repo
	list, err = db.ListDeploymentsForRepo(context.TODO(), resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list deployments for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Deployments) {
		t.Errorf("ListDeploymentsForRepo() is %v, want %v", count, len(resources.Deployments))
	}
	if diff := cmp.Diff([]*api.Deployment{resources.Deployments[1], resources.Deployments[0]}, list); diff != "" {
		t.Errorf("ListDeploymentsForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["ListDeploymentsForRepo"] = true

	// lookup the deployments by name
	for _, deployment := range resources.Deployments {
		repo := resources.Repos[deployment.GetRepo().GetID()-1]
		got, err := db.GetDeploymentForRepo(context.TODO(), repo, deployment.GetNumber())
		if err != nil {
			t.Errorf("unable to get deployment %d for repo %d: %v", deployment.GetID(), repo.GetID(), err)
		}
		if diff := cmp.Diff(deployment, got); diff != "" {
			t.Errorf("GetDeploymentForRepo() mismatch (-want +got):\n%s", diff)
		}
	}
	methods["GetDeploymentForRepo"] = true

	// update the deployments
	for _, deployment := range resources.Deployments {
		_, err = db.UpdateDeployment(context.TODO(), deployment)
		if err != nil {
			t.Errorf("unable to update deployment %d: %v", deployment.GetID(), err)
		}

		// lookup the deployment by ID
		got, err := db.GetDeployment(context.TODO(), deployment.GetID())
		if err != nil {
			t.Errorf("unable to get deployment %d by ID: %v", deployment.GetID(), err)
		}
		if diff := cmp.Diff(deployment, got); diff != "" {
			t.Errorf("GetDeployment() mismatch (-want +got):\n%s", diff)
		}
	}
	methods["UpdateDeployment"] = true
	methods["GetDeployment"] = true

	// delete the deployments
	for _, deployment := range resources.Deployments {
		err = db.DeleteDeployment(context.TODO(), deployment)
		if err != nil {
			t.Errorf("unable to delete hook %d: %v", deployment.GetID(), err)
		}
	}
	methods["DeleteDeployment"] = true

	// delete the builds
	for _, build := range resources.Builds {
		err = db.DeleteBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to delete build: %v", err)
		}
	}

	// delete the repos for hook related functions
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo: %v", err)
		}
	}

	// delete the users for the hook related functions
	for _, user := range resources.Users {
		err = db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user: %v", err)
		}
	}

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for deployments", method)
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

	// create the users for hook related functions (owners of repos)
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}

	// create the repos for hook related functions
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	// create the builds for hook related functions
	for _, build := range resources.Builds {
		_, err := db.CreateBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to create build %d: %v", build.GetID(), err)
		}
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
	if diff := cmp.Diff(resources.Hooks, list); diff != "" {
		t.Errorf("ListHooks() mismatch (-want +got):\n%s", diff)
	}
	methods["ListHooks"] = true

	// list the hooks for a repo
	list, count, err = db.ListHooksForRepo(context.TODO(), resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list hooks for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	// only 2 of 3 hooks belong to Repos[0] repo
	if int(count) != len(resources.Hooks)-1 {
		t.Errorf("ListHooksForRepo() is %v, want %v", count, len(resources.Hooks))
	}
	if diff := cmp.Diff([]*api.Hook{resources.Hooks[2], resources.Hooks[0]}, list); diff != "" {
		t.Errorf("ListHooksForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["ListHooksForRepo"] = true

	// lookup the last build by repo
	got, err := db.LastHookForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to get last hook for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if diff := cmp.Diff(resources.Hooks[2], got); diff != "" {
		t.Errorf("LastHookForRepo() mismatch (-want +got):\n%s", diff)
	}
	methods["LastHookForRepo"] = true

	// lookup a hook with matching webhook_id
	got, err = db.GetHookByWebhookID(context.TODO(), resources.Hooks[2].GetWebhookID())
	if err != nil {
		t.Errorf("unable to get last hook for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if diff := cmp.Diff(resources.Hooks[2], got); diff != "" {
		t.Errorf("GetHookByWebhookID() mismatch (-want +got):\n%s", diff)
	}
	methods["GetHookByWebhookID"] = true

	// lookup the hooks by name
	for _, hook := range resources.Hooks {
		got, err = db.GetHookForRepo(context.TODO(), hook.GetRepo(), hook.GetNumber())
		if err != nil {
			t.Errorf("unable to get hook %d for repo %d: %v", hook.GetID(), hook.GetRepo().GetID(), err)
		}
		if diff := cmp.Diff(hook, got); diff != "" {
			t.Errorf("GetHookForRepo() mismatch (-want +got):\n%s", diff)
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
		if diff := cmp.Diff(hook, got); diff != "" {
			t.Errorf("GetHook() mismatch (-want +got):\n%s", diff)
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

	// delete the builds
	for _, build := range resources.Builds {
		err = db.DeleteBuild(context.TODO(), build)
		if err != nil {
			t.Errorf("unable to delete build: %v", err)
		}
	}

	// delete the repos for hook related functions
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo: %v", err)
		}
	}

	// delete the users for the hook related functions
	for _, user := range resources.Users {
		err = db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user: %v", err)
		}
	}

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for hooks", method)
		}
	}
}

func testJWKs(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for jwks
	methods := make(map[string]bool)
	// capture the element type of the jwk interface
	element := reflect.TypeOf(new(dbJWK.JWKInterface)).Elem()
	// iterate through all methods found in the jwk interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for jwks
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	for i := 0; i < resources.JWKs.Len(); i++ {
		jk, _ := resources.JWKs.Key(i)

		jkPub, _ := jk.(jwk.RSAPublicKey)

		err := db.CreateJWK(context.TODO(), jkPub)
		if err != nil {
			t.Errorf("unable to create jwk %s: %v", jkPub.KeyID(), err)
		}
	}
	methods["CreateJWK"] = true

	list, err := db.ListJWKs(context.TODO())
	if err != nil {
		t.Errorf("unable to list jwks: %v", err)
	}

	if !reflect.DeepEqual(resources.JWKs, list) {
		t.Errorf("ListJWKs() mismatch, want %v, got %v", resources.JWKs, list)
	}

	methods["ListJWKs"] = true

	for i := 0; i < resources.JWKs.Len(); i++ {
		jk, _ := resources.JWKs.Key(i)

		jkPub, _ := jk.(jwk.RSAPublicKey)

		got, err := db.GetActiveJWK(context.TODO(), jkPub.KeyID())
		if err != nil {
			t.Errorf("unable to get jwk %s: %v", jkPub.KeyID(), err)
		}

		if !cmp.Equal(jkPub, got, testutils.JwkKeyOpts) {
			t.Errorf("GetJWK() is %v, want %v", got, jkPub)
		}
	}

	methods["GetActiveJWK"] = true

	err = db.RotateKeys(context.TODO())
	if err != nil {
		t.Errorf("unable to rotate keys: %v", err)
	}

	for i := 0; i < resources.JWKs.Len(); i++ {
		jk, _ := resources.JWKs.Key(i)

		jkPub, _ := jk.(jwk.RSAPublicKey)

		_, err := db.GetActiveJWK(context.TODO(), jkPub.KeyID())
		if err == nil {
			t.Errorf("GetActiveJWK() should return err after rotation")
		}
	}

	methods["RotateKeys"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for jwks", method)
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
		err := db.CreateLog(context.TODO(), log)
		if err != nil {
			t.Errorf("unable to create log %d: %v", log.GetID(), err)
		}
	}
	methods["CreateLog"] = true

	// count the logs
	count, err := db.CountLogs(context.TODO())
	if err != nil {
		t.Errorf("unable to count logs: %v", err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(resources.Logs))
	}
	methods["CountLogs"] = true

	// count the logs for a build
	count, err = db.CountLogsForBuild(context.TODO(), resources.Builds[0])
	if err != nil {
		t.Errorf("unable to count logs for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("CountLogs() is %v, want %v", count, len(resources.Logs))
	}
	methods["CountLogsForBuild"] = true

	// list the logs
	list, err := db.ListLogs(context.TODO())
	if err != nil {
		t.Errorf("unable to list logs: %v", err)
	}
	if diff := cmp.Diff(resources.Logs, list); diff != "" {
		t.Errorf("ListLogs() mismatch (-want +got):\n%s", diff)
	}
	methods["ListLogs"] = true

	// list the logs for a build
	list, count, err = db.ListLogsForBuild(context.TODO(), resources.Builds[0], 1, 10)
	if err != nil {
		t.Errorf("unable to list logs for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Logs) {
		t.Errorf("ListLogsForBuild() is %v, want %v", count, len(resources.Logs))
	}
	if diff := cmp.Diff(resources.Logs, list); diff != "" {
		t.Errorf("ListLogsForBuild() mismatch (-want +got):\n%s", diff)
	}
	methods["ListLogsForBuild"] = true

	// lookup the logs by service
	for _, log := range []*api.Log{resources.Logs[0], resources.Logs[1]} {
		service := resources.Services[log.GetServiceID()-1]
		got, err := db.GetLogForService(context.TODO(), service)
		if err != nil {
			t.Errorf("unable to get log %d for service %d: %v", log.GetID(), service.GetID(), err)
		}
		if !cmp.Equal(got, log) {
			t.Errorf("GetLogForService() is %v, want %v", got, log)
		}
	}
	methods["GetLogForService"] = true

	// lookup the logs by service
	for _, log := range []*api.Log{resources.Logs[2], resources.Logs[3]} {
		step := resources.Steps[log.GetStepID()-1]
		got, err := db.GetLogForStep(context.TODO(), step)
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
		err = db.UpdateLog(context.TODO(), log)
		if err != nil {
			t.Errorf("unable to update log %d: %v", log.GetID(), err)
		}

		// lookup the log by ID
		got, err := db.GetLog(context.TODO(), log.GetID())
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
		err = db.DeleteLog(context.TODO(), log)
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

	// create owners
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}

	// create the repos
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
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
	if diff := cmp.Diff(resources.Pipelines, list); diff != "" {
		t.Errorf("ListPipelines() mismatch (-want +got):\n%s", diff)
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
	if diff := cmp.Diff(resources.Pipelines, list); diff != "" {
		t.Errorf("ListPipelines() mismatch (-want +got):\n%s", diff)
	}
	methods["ListPipelinesForRepo"] = true

	// lookup the pipelines by name
	for _, pipeline := range resources.Pipelines {
		repo := resources.Repos[pipeline.GetRepo().GetID()-1]
		got, err := db.GetPipelineForRepo(context.TODO(), pipeline.GetCommit(), repo)
		if err != nil {
			t.Errorf("unable to get pipeline %d for repo %d: %v", pipeline.GetID(), repo.GetID(), err)
		}
		if diff := cmp.Diff(pipeline, got); diff != "" {
			t.Errorf("GetPipelineForRepo() mismatch (-want +got):\n%s", diff)
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
		if diff := cmp.Diff(pipeline, got); diff != "" {
			t.Errorf("GetPipeline() mismatch (-want +got):\n%s", diff)
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

	// delete the repos
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}

	// delete the owners
	for _, user := range resources.Users {
		err := db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}

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

	// create owners
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
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
	if diff := cmp.Diff(resources.Repos, list); diff != "" {
		t.Errorf("ListRepos() mismatch (-want +got):\n%s", diff)
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
	if diff := cmp.Diff(resources.Repos, list); diff != "" {
		t.Errorf("ListReposForOrg() mismatch (-want +got):\n%s", diff)
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
	if diff := cmp.Diff(resources.Repos, list); diff != "" {
		t.Errorf("ListReposForUser() mismatch (-want +got):\n%s", diff)
	}
	methods["ListReposForUser"] = true

	// lookup the repos by name
	for _, repo := range resources.Repos {
		got, err := db.GetRepoForOrg(context.TODO(), repo.GetOrg(), repo.GetName())
		if err != nil {
			t.Errorf("unable to get repo %d by org: %v", repo.GetID(), err)
		}
		if diff := cmp.Diff(repo, got); diff != "" {
			t.Errorf("GetRepoForOrg() mismatch (-want +got):\n%s", diff)
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

	// delete the owners
	for _, user := range resources.Users {
		err := db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}

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

	// create owners
	for _, user := range resources.Users {
		_, err := db.CreateUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to create user %d: %v", user.GetID(), err)
		}
	}

	// create the repos
	for _, repo := range resources.Repos {
		_, err := db.CreateRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to create repo %d: %v", repo.GetID(), err)
		}
	}

	// create the schedules
	for _, schedule := range resources.Schedules {
		_, err := db.CreateSchedule(context.TODO(), schedule)
		if err != nil {
			t.Errorf("unable to create schedule %d: %v", schedule.GetID(), err)
		}
	}
	methods["CreateSchedule"] = true

	// count the schedules
	count, err := db.CountSchedules(context.TODO())
	if err != nil {
		t.Errorf("unable to count schedules: %v", err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("CountSchedules() is %v, want %v", count, len(resources.Schedules))
	}
	methods["CountSchedules"] = true

	// count the schedules for a repo
	count, err = db.CountSchedulesForRepo(context.TODO(), resources.Repos[0])
	if err != nil {
		t.Errorf("unable to count schedules for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("CountSchedulesForRepo() is %v, want %v", count, len(resources.Schedules))
	}
	methods["CountSchedulesForRepo"] = true

	// list the schedules
	list, err := db.ListSchedules(context.TODO())
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !cmp.Equal(list, resources.Schedules, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListSchedules() is %v, want %v", list, resources.Schedules)
	}
	methods["ListSchedules"] = true

	// list the active schedules
	list, err = db.ListActiveSchedules(context.TODO())
	if err != nil {
		t.Errorf("unable to list schedules: %v", err)
	}
	if !cmp.Equal(list, resources.Schedules, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListActiveSchedules() is %v, want %v", list, resources.Schedules)
	}
	methods["ListActiveSchedules"] = true

	// list the schedules for a repo
	list, count, err = db.ListSchedulesForRepo(context.TODO(), resources.Repos[0], 1, 10)
	if err != nil {
		t.Errorf("unable to count schedules for repo %d: %v", resources.Repos[0].GetID(), err)
	}
	if int(count) != len(resources.Schedules) {
		t.Errorf("ListSchedulesForRepo() is %v, want %v", count, len(resources.Schedules))
	}
	if !cmp.Equal(list, []*api.Schedule{resources.Schedules[1], resources.Schedules[0]}, CmpOptApproxUpdatedAt()) {
		t.Errorf("ListSchedulesForRepo() is %v, want %v", list, []*api.Schedule{resources.Schedules[1], resources.Schedules[0]})
	}
	methods["ListSchedulesForRepo"] = true

	// lookup the schedules by name
	for _, schedule := range resources.Schedules {
		repo := resources.Repos[schedule.GetRepo().GetID()-1]
		got, err := db.GetScheduleForRepo(context.TODO(), repo, schedule.GetName())
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
		got, err := db.UpdateSchedule(context.TODO(), schedule, true)
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
		err = db.DeleteSchedule(context.TODO(), schedule)
		if err != nil {
			t.Errorf("unable to delete schedule %d: %v", schedule.GetID(), err)
		}
	}
	methods["DeleteSchedule"] = true

	// delete the repos
	for _, repo := range resources.Repos {
		err = db.DeleteRepo(context.TODO(), repo)
		if err != nil {
			t.Errorf("unable to delete repo %d: %v", repo.GetID(), err)
		}
	}

	// delete the owners
	for _, user := range resources.Users {
		err := db.DeleteUser(context.TODO(), user)
		if err != nil {
			t.Errorf("unable to delete user %d: %v", user.GetID(), err)
		}
	}

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
		_, err := db.CreateSecret(context.TODO(), secret)
		if err != nil {
			t.Errorf("unable to create secret %d: %v", secret.GetID(), err)
		}
	}
	methods["CreateSecret"] = true

	// count the secrets
	count, err := db.CountSecrets(context.TODO())
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
			count, err = db.CountSecretsForOrg(context.TODO(), secret.GetOrg(), nil)
			if err != nil {
				t.Errorf("unable to count secrets for org %s: %v", secret.GetOrg(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForOrg() is %v, want %v", count, 1)
			}
			methods["CountSecretsForOrg"] = true
		case constants.SecretRepo:
			// count the secrets for a repo
			count, err = db.CountSecretsForRepo(context.TODO(), resources.Repos[0], nil)
			if err != nil {
				t.Errorf("unable to count secrets for repo %d: %v", resources.Repos[0].GetID(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForRepo() is %v, want %v", count, 1)
			}
			methods["CountSecretsForRepo"] = true
		case constants.SecretShared:
			// count the secrets for a team
			count, err = db.CountSecretsForTeam(context.TODO(), secret.GetOrg(), secret.GetTeam(), nil)
			if err != nil {
				t.Errorf("unable to count secrets for team %s: %v", secret.GetTeam(), err)
			}
			if int(count) != 1 {
				t.Errorf("CountSecretsForTeam() is %v, want %v", count, 1)
			}
			methods["CountSecretsForTeam"] = true

			// count the secrets for a list of teams
			count, err = db.CountSecretsForTeams(context.TODO(), secret.GetOrg(), []string{secret.GetTeam()}, nil)
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
	list, err := db.ListSecrets(context.TODO())
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
			list, count, err = db.ListSecretsForOrg(context.TODO(), secret.GetOrg(), nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for org %s: %v", secret.GetOrg(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForOrg() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*api.Secret{secret}) {
				t.Errorf("ListSecretsForOrg() is %v, want %v", list, []*api.Secret{secret})
			}
			methods["ListSecretsForOrg"] = true
		case constants.SecretRepo:
			// list the secrets for a repo
			list, count, err = db.ListSecretsForRepo(context.TODO(), resources.Repos[0], nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for repo %d: %v", resources.Repos[0].GetID(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForRepo() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*api.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForRepo() is %v, want %v", list, []*api.Secret{secret})
			}
			methods["ListSecretsForRepo"] = true
		case constants.SecretShared:
			// list the secrets for a team
			list, count, err = db.ListSecretsForTeam(context.TODO(), secret.GetOrg(), secret.GetTeam(), nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for team %s: %v", secret.GetTeam(), err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForTeam() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*api.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForTeam() is %v, want %v", list, []*api.Secret{secret})
			}
			methods["ListSecretsForTeam"] = true

			// list the secrets for a list of teams
			list, count, err = db.ListSecretsForTeams(context.TODO(), secret.GetOrg(), []string{secret.GetTeam()}, nil, 1, 10)
			if err != nil {
				t.Errorf("unable to list secrets for teams %s: %v", []string{secret.GetTeam()}, err)
			}
			if int(count) != 1 {
				t.Errorf("ListSecretsForTeams() is %v, want %v", count, 1)
			}
			if !cmp.Equal(list, []*api.Secret{secret}, CmpOptApproxUpdatedAt()) {
				t.Errorf("ListSecretsForTeams() is %v, want %v", list, []*api.Secret{secret})
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
			got, err := db.GetSecretForOrg(context.TODO(), secret.GetOrg(), secret.GetName())
			if err != nil {
				t.Errorf("unable to get secret %d for org %s: %v", secret.GetID(), secret.GetOrg(), err)
			}
			if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
				t.Errorf("GetSecretForOrg() is %v, want %v", got, secret)
			}
			methods["GetSecretForOrg"] = true
		case constants.SecretRepo:
			// lookup the secret by repo
			got, err := db.GetSecretForRepo(context.TODO(), secret.GetName(), resources.Repos[0])
			if err != nil {
				t.Errorf("unable to get secret %d for repo %d: %v", secret.GetID(), resources.Repos[0].GetID(), err)
			}
			if !cmp.Equal(got, secret, CmpOptApproxUpdatedAt()) {
				t.Errorf("GetSecretForRepo() is %v, want %v", got, secret)
			}
			methods["GetSecretForRepo"] = true
		case constants.SecretShared:
			// lookup the secret by team
			got, err := db.GetSecretForTeam(context.TODO(), secret.GetOrg(), secret.GetTeam(), secret.GetName())
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
		got, err := db.UpdateSecret(context.TODO(), secret)
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
		err = db.DeleteSecret(context.TODO(), secret)
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
		_, err := db.CreateService(context.TODO(), service)
		if err != nil {
			t.Errorf("unable to create service %d: %v", service.GetID(), err)
		}
	}
	methods["CreateService"] = true

	// count the services
	count, err := db.CountServices(context.TODO())
	if err != nil {
		t.Errorf("unable to count services: %v", err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServices() is %v, want %v", count, len(resources.Services))
	}
	methods["CountServices"] = true

	// count the services for a build
	count, err = db.CountServicesForBuild(context.TODO(), resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Services) {
		t.Errorf("CountServicesForBuild() is %v, want %v", count, len(resources.Services))
	}
	methods["CountServicesForBuild"] = true

	// list the services
	list, err := db.ListServices(context.TODO())
	if err != nil {
		t.Errorf("unable to list services: %v", err)
	}
	if !cmp.Equal(list, resources.Services) {
		t.Errorf("ListServices() is %v, want %v", list, resources.Services)
	}
	methods["ListServices"] = true

	// list the services for a build
	list, count, err = db.ListServicesForBuild(context.TODO(), resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list services for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !cmp.Equal(list, []*api.Service{resources.Services[1], resources.Services[0]}) {
		t.Errorf("ListServicesForBuild() is %v, want %v", list, []*api.Service{resources.Services[1], resources.Services[0]})
	}
	if int(count) != len(resources.Services) {
		t.Errorf("ListServicesForBuild() is %v, want %v", count, len(resources.Services))
	}
	methods["ListServicesForBuild"] = true

	expected := map[string]float64{
		"#init":                        1,
		"target/vela-git-slim:v0.12.0": 1,
	}
	images, err := db.ListServiceImageCount(context.TODO())
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
	statuses, err := db.ListServiceStatusCount(context.TODO())
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
		got, err := db.GetServiceForBuild(context.TODO(), build, service.GetNumber())
		if err != nil {
			t.Errorf("unable to get service %d for build %d: %v", service.GetID(), build.GetID(), err)
		}
		if !cmp.Equal(got, service) {
			t.Errorf("GetServiceForBuild() is %v, want %v", got, service)
		}
	}
	methods["GetServiceForBuild"] = true

	// clean the services
	count, err = db.CleanServices(context.TODO(), "integration testing", 1563474090)
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
		got, err := db.UpdateService(context.TODO(), service)
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
		err = db.DeleteService(context.TODO(), service)
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
	ctx := context.TODO()

	// create the steps
	for _, step := range resources.Steps {
		_, err := db.CreateStep(ctx, step)
		if err != nil {
			t.Errorf("unable to create step %d: %v", step.GetID(), err)
		}
	}
	methods["CreateStep"] = true

	// count the steps
	count, err := db.CountSteps(ctx)
	if err != nil {
		t.Errorf("unable to count steps: %v", err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountSteps() is %v, want %v", count, len(resources.Steps))
	}
	methods["CountSteps"] = true

	// count the steps for a build
	count, err = db.CountStepsForBuild(ctx, resources.Builds[0], nil)
	if err != nil {
		t.Errorf("unable to count steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("CountStepsForBuild() is %v, want %v", count, len(resources.Steps))
	}
	methods["CountStepsForBuild"] = true

	// list the steps
	list, err := db.ListSteps(ctx)
	if err != nil {
		t.Errorf("unable to list steps: %v", err)
	}
	if !cmp.Equal(list, resources.Steps) {
		t.Errorf("ListSteps() is %v, want %v", list, resources.Steps)
	}
	methods["ListSteps"] = true

	// list the steps for a build
	list, count, err = db.ListStepsForBuild(ctx, resources.Builds[0], nil, 1, 10)
	if err != nil {
		t.Errorf("unable to list steps for build %d: %v", resources.Builds[0].GetID(), err)
	}
	if !cmp.Equal(list, []*api.Step{resources.Steps[1], resources.Steps[0]}) {
		t.Errorf("ListStepsForBuild() is %v, want %v", list, []*api.Step{resources.Steps[1], resources.Steps[0]})
	}
	if int(count) != len(resources.Steps) {
		t.Errorf("ListStepsForBuild() is %v, want %v", count, len(resources.Steps))
	}
	methods["ListStepsForBuild"] = true

	expected := map[string]float64{
		"#init":                        1,
		"target/vela-git-slim:v0.12.0": 1,
	}
	images, err := db.ListStepImageCount(ctx)
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
	statuses, err := db.ListStepStatusCount(ctx)
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
		got, err := db.GetStepForBuild(ctx, build, step.GetNumber())
		if err != nil {
			t.Errorf("unable to get step %d for build %d: %v", step.GetID(), build.GetID(), err)
		}
		if !cmp.Equal(got, step) {
			t.Errorf("GetStepForBuild() is %v, want %v", got, step)
		}
	}
	methods["GetStepForBuild"] = true

	// clean the steps
	count, err = db.CleanSteps(ctx, "integration testing", 1563474090)
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
		got, err := db.UpdateStep(ctx, step)
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
		err = db.DeleteStep(ctx, step)
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

	userOne := new(api.User)
	userOne.SetID(1)
	userOne.SetName("octocat")
	userOne.SetToken("")
	userOne.SetRefreshToken("")
	userOne.SetFavorites(nil)
	userOne.SetDashboards(nil)
	userOne.SetActive(false)
	userOne.SetAdmin(false)

	userTwo := new(api.User)
	userTwo.SetID(2)
	userTwo.SetName("octokitty")
	userTwo.SetToken("")
	userTwo.SetRefreshToken("")
	userTwo.SetFavorites(nil)
	userTwo.SetDashboards(nil)
	userTwo.SetActive(false)
	userTwo.SetAdmin(false)

	liteUsers := []*api.User{userOne, userTwo}

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
	list, err := db.ListWorkers(context.TODO(), "all", time.Now().Unix(), 0)
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

func testSettings(t *testing.T, db Interface, resources *Resources) {
	// create a variable to track the number of methods called for settings
	methods := make(map[string]bool)
	// capture the element type of the settings interface
	element := reflect.TypeOf(new(dbSettings.SettingsInterface)).Elem()
	// iterate through all methods found in the settings interface
	for i := 0; i < element.NumMethod(); i++ {
		// skip tracking the methods to create indexes and tables for settings
		// since those are already called when the database engine starts
		if strings.Contains(element.Method(i).Name, "Index") ||
			strings.Contains(element.Method(i).Name, "Table") {
			continue
		}

		// add the method name to the list of functions
		methods[element.Method(i).Name] = false
	}

	// create the settings
	for _, s := range resources.Platform {
		_, err := db.CreateSettings(context.TODO(), s)
		if err != nil {
			t.Errorf("unable to create settings %d: %v", s.GetID(), err)
		}
	}
	methods["CreateSettings"] = true

	// update the settings
	for _, s := range resources.Platform {
		s.SetCloneImage("target/vela-git-slim:abc123")
		got, err := db.UpdateSettings(context.TODO(), s)
		if err != nil {
			t.Errorf("unable to update settings %d: %v", s.GetID(), err)
		}

		if !cmp.Equal(got, s) {
			t.Errorf("UpdateSettings() is %v, want %v", got, s)
		}
	}
	methods["UpdateSettings"] = true
	methods["GetSettings"] = true

	// ensure we called all the methods we expected to
	for method, called := range methods {
		if !called {
			t.Errorf("method %s was not called for settings", method)
		}
	}
}

func newResources() *Resources {
	userOne := new(api.User)
	userOne.SetID(1)
	userOne.SetName("octocat")
	userOne.SetToken("superSecretToken")
	userOne.SetRefreshToken("superSecretRefreshToken")
	userOne.SetFavorites([]string{"github/octocat"})
	userOne.SetActive(true)
	userOne.SetAdmin(false)
	userOne.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"})

	userTwo := new(api.User)
	userTwo.SetID(2)
	userTwo.SetName("octokitty")
	userTwo.SetToken("superSecretToken")
	userTwo.SetRefreshToken("superSecretRefreshToken")
	userTwo.SetFavorites([]string{"github/octocat"})
	userTwo.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"})
	userTwo.SetActive(true)
	userTwo.SetAdmin(false)

	repoOne := new(api.Repo)
	repoOne.SetID(1)
	repoOne.SetOwner(userOne.Crop())
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
	repoOne.SetPipelineType("")
	repoOne.SetPreviousName("")
	repoOne.SetApproveBuild(constants.ApproveNever)
	repoOne.SetAllowEvents(api.NewEventsFromMask(1))
	repoOne.SetApprovalTimeout(7)
	repoOne.SetInstallID(0)

	repoTwo := new(api.Repo)
	repoTwo.SetID(2)
	repoTwo.SetOwner(userOne.Crop())
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
	repoTwo.SetPipelineType("")
	repoTwo.SetPreviousName("")
	repoTwo.SetApproveBuild(constants.ApproveForkAlways)
	repoTwo.SetAllowEvents(api.NewEventsFromMask(1))
	repoTwo.SetApprovalTimeout(7)
	repoTwo.SetInstallID(0)

	buildOne := new(api.Build)
	buildOne.SetID(1)
	buildOne.SetRepo(repoOne)
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
	buildOne.SetDeployNumber(0)
	buildOne.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	buildOne.SetClone("https://github.com/github/octocat.git")
	buildOne.SetSource("https://github.com/github/octocat/deployments/1")
	buildOne.SetTitle("push received from https://github.com/github/octocat")
	buildOne.SetMessage("First commit...")
	buildOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	buildOne.SetSender("OctoKitty")
	buildOne.SetSenderSCMID("123")
	buildOne.SetFork(false)
	buildOne.SetAuthor("OctoKitty")
	buildOne.SetEmail("OctoKitty@github.com")
	buildOne.SetLink("https://example.company.com/github/octocat/1")
	buildOne.SetBranch("main")
	buildOne.SetRef("refs/heads/main")
	buildOne.SetBaseRef("")
	buildOne.SetHeadRef("changes")
	buildOne.SetHost("example.company.com")
	buildOne.SetRoute("vela")
	buildOne.SetRuntime("docker")
	buildOne.SetDistribution("linux")
	buildOne.SetApprovedAt(1563474078)
	buildOne.SetApprovedBy("OctoCat")

	buildTwo := new(api.Build)
	buildTwo.SetID(2)
	buildTwo.SetRepo(repoOne)
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
	buildTwo.SetDeployNumber(0)
	buildTwo.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	buildTwo.SetClone("https://github.com/github/octocat.git")
	buildTwo.SetSource("https://github.com/github/octocat/deployments/1")
	buildTwo.SetTitle("pull_request received from https://github.com/github/octocat")
	buildTwo.SetMessage("Second commit...")
	buildTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	buildTwo.SetSender("OctoKitty")
	buildTwo.SetSenderSCMID("123")
	buildTwo.SetFork(false)
	buildTwo.SetAuthor("OctoKitty")
	buildTwo.SetEmail("OctoKitty@github.com")
	buildTwo.SetLink("https://example.company.com/github/octocat/2")
	buildTwo.SetBranch("main")
	buildTwo.SetRef("refs/heads/main")
	buildTwo.SetBaseRef("")
	buildTwo.SetHeadRef("changes")
	buildTwo.SetHost("example.company.com")
	buildTwo.SetRoute("vela")
	buildTwo.SetRuntime("docker")
	buildTwo.SetDistribution("linux")
	buildTwo.SetApprovedAt(1563474078)
	buildTwo.SetApprovedBy("OctoCat")

	dashRepo := new(api.DashboardRepo)
	dashRepo.SetID(1)
	dashRepo.SetName("go-vela/server")
	dashRepo.SetBranches([]string{"main"})
	dashRepo.SetEvents([]string{"push"})

	// crop and set "-" JSON tag fields to nil for dashboard admins
	dashboardAdmins := []*api.User{userOne.Crop(), userTwo.Crop()}
	for _, admin := range dashboardAdmins {
		admin.Token = nil
		admin.RefreshToken = nil
	}

	dashboardOne := new(api.Dashboard)
	dashboardOne.SetID("ba657dab-bc6e-421f-9188-86272bd0069a")
	dashboardOne.SetName("vela")
	dashboardOne.SetCreatedAt(1)
	dashboardOne.SetCreatedBy("octocat")
	dashboardOne.SetUpdatedAt(2)
	dashboardOne.SetUpdatedBy("octokitty")
	dashboardOne.SetAdmins(dashboardAdmins)
	dashboardOne.SetRepos([]*api.DashboardRepo{dashRepo})

	dashboardTwo := new(api.Dashboard)
	dashboardTwo.SetID("45bcf19b-c151-4e2d-b8c6-80a62ba2eae7")
	dashboardTwo.SetName("vela")
	dashboardTwo.SetCreatedAt(1)
	dashboardTwo.SetCreatedBy("octocat")
	dashboardTwo.SetUpdatedAt(2)
	dashboardTwo.SetUpdatedBy("octokitty")
	dashboardTwo.SetAdmins(dashboardAdmins)
	dashboardTwo.SetRepos([]*api.DashboardRepo{dashRepo})

	executableOne := new(api.BuildExecutable)
	executableOne.SetID(1)
	executableOne.SetBuildID(1)
	executableOne.SetData([]byte("foo"))

	executableTwo := new(api.BuildExecutable)
	executableTwo.SetID(2)
	executableTwo.SetBuildID(2)
	executableTwo.SetData([]byte("foo"))

	trimmedBuildOne := *buildOne
	trimmedBuildOne.Repo = &api.Repo{ID: repoOne.ID}

	trimmedBuildTwo := *buildTwo
	trimmedBuildTwo.Repo = &api.Repo{ID: repoOne.ID}

	deploymentOne := new(api.Deployment)
	deploymentOne.SetID(1)
	deploymentOne.SetNumber(1)
	deploymentOne.SetRepo(repoOne)
	deploymentOne.SetURL("https://github.com/github/octocat/deployments/1")
	deploymentOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	deploymentOne.SetRef("refs/heads/main")
	deploymentOne.SetTask("vela-deploy")
	deploymentOne.SetTarget("production")
	deploymentOne.SetDescription("Deployment request from Vela")
	deploymentOne.SetPayload(map[string]string{"foo": "test1"})
	deploymentOne.SetCreatedAt(1)
	deploymentOne.SetCreatedBy("octocat")
	deploymentOne.SetBuilds([]*api.Build{&trimmedBuildOne})

	deploymentTwo := new(api.Deployment)
	deploymentTwo.SetID(2)
	deploymentTwo.SetNumber(2)
	deploymentTwo.SetRepo(repoOne)
	deploymentTwo.SetURL("https://github.com/github/octocat/deployments/2")
	deploymentTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	deploymentTwo.SetRef("refs/heads/main")
	deploymentTwo.SetTask("vela-deploy")
	deploymentTwo.SetTarget("production")
	deploymentTwo.SetDescription("Deployment request from Vela")
	deploymentTwo.SetPayload(map[string]string{"foo": "test1"})
	deploymentTwo.SetCreatedAt(1)
	deploymentTwo.SetCreatedBy("octocat")

	hookOne := new(api.Hook)
	hookOne.SetID(1)
	hookOne.SetRepo(repoOne)
	hookOne.SetBuild(&trimmedBuildOne)
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

	hookTwo := new(api.Hook)
	hookTwo.SetID(2)
	hookTwo.SetRepo(repoTwo)
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

	hookThree := new(api.Hook)
	hookThree.SetID(3)
	hookThree.SetRepo(repoOne)
	hookThree.SetBuild(&trimmedBuildTwo)
	hookThree.SetNumber(3)
	hookThree.SetSourceID("c8da1302-07d6-11ea-882f-6793bca275b8")
	hookThree.SetCreated(time.Now().UTC().Unix())
	hookThree.SetHost("github.com")
	hookThree.SetEvent("push")
	hookThree.SetEventAction("")
	hookThree.SetBranch("main")
	hookThree.SetError("")
	hookThree.SetStatus("success")
	hookThree.SetLink("https://github.com/github/octocat/settings/hooks/1")
	hookThree.SetWebhookID(78910)

	jwkOne := testutils.JWK()
	jwkTwo := testutils.JWK()

	jwkSet := jwk.NewSet()

	_ = jwkSet.AddKey(jwkOne)

	_ = jwkSet.AddKey(jwkTwo)

	logServiceOne := new(api.Log)
	logServiceOne.SetID(1)
	logServiceOne.SetBuildID(1)
	logServiceOne.SetRepoID(1)
	logServiceOne.SetServiceID(1)
	logServiceOne.SetStepID(0)
	logServiceOne.SetData([]byte("foo"))

	logServiceTwo := new(api.Log)
	logServiceTwo.SetID(2)
	logServiceTwo.SetBuildID(1)
	logServiceTwo.SetRepoID(1)
	logServiceTwo.SetServiceID(2)
	logServiceTwo.SetStepID(0)
	logServiceTwo.SetData([]byte("foo"))

	logStepOne := new(api.Log)
	logStepOne.SetID(3)
	logStepOne.SetBuildID(1)
	logStepOne.SetRepoID(1)
	logStepOne.SetServiceID(0)
	logStepOne.SetStepID(1)
	logStepOne.SetData([]byte("foo"))

	logStepTwo := new(api.Log)
	logStepTwo.SetID(4)
	logStepTwo.SetBuildID(1)
	logStepTwo.SetRepoID(1)
	logStepTwo.SetServiceID(0)
	logStepTwo.SetStepID(2)
	logStepTwo.SetData([]byte("foo"))

	pipelineOne := new(api.Pipeline)
	pipelineOne.SetID(1)
	pipelineOne.SetRepo(repoOne)
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
	pipelineOne.SetWarnings([]string{})
	pipelineOne.SetData([]byte("version: 1"))

	pipelineTwo := new(api.Pipeline)
	pipelineTwo.SetID(2)
	pipelineTwo.SetRepo(repoOne)
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
	pipelineTwo.SetWarnings([]string{"42:this is a warning"})
	pipelineTwo.SetData([]byte("version: 1"))

	currTime := time.Now().UTC()
	nextTime, _ := gronx.NextTickAfter("0 0 * * *", currTime, false)

	scheduleOne := new(api.Schedule)
	scheduleOne.SetID(1)
	scheduleOne.SetRepo(repoOne)
	scheduleOne.SetActive(true)
	scheduleOne.SetName("nightly")
	scheduleOne.SetEntry("0 0 * * *")
	scheduleOne.SetCreatedAt(time.Now().UTC().Unix())
	scheduleOne.SetCreatedBy("octocat")
	scheduleOne.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	scheduleOne.SetUpdatedBy("octokitty")
	scheduleOne.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	scheduleOne.SetBranch("main")
	scheduleOne.SetError("no version: YAML property provided")
	scheduleOne.SetNextRun(nextTime.Unix())

	currTime = time.Now().UTC()
	nextTime, _ = gronx.NextTickAfter("0 * * * *", currTime, false)

	scheduleTwo := new(api.Schedule)
	scheduleTwo.SetID(2)
	scheduleTwo.SetRepo(repoOne)
	scheduleTwo.SetActive(true)
	scheduleTwo.SetName("hourly")
	scheduleTwo.SetEntry("0 * * * *")
	scheduleTwo.SetCreatedAt(time.Now().UTC().Unix())
	scheduleTwo.SetCreatedBy("octocat")
	scheduleTwo.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	scheduleTwo.SetUpdatedBy("octokitty")
	scheduleTwo.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	scheduleTwo.SetBranch("main")
	scheduleTwo.SetError("no version: YAML property provided")
	scheduleTwo.SetNextRun(nextTime.Unix())

	secretOrg := new(api.Secret)
	secretOrg.SetID(1)
	secretOrg.SetOrg("github")
	secretOrg.SetRepo("*")
	secretOrg.SetTeam("")
	secretOrg.SetName("foo")
	secretOrg.SetValue("bar")
	secretOrg.SetType("org")
	secretOrg.SetImages([]string{"alpine"})
	secretOrg.SetAllowEvents(api.NewEventsFromMask(1))
	secretOrg.SetAllowCommand(true)
	secretOrg.SetAllowSubstitution(true)
	secretOrg.SetCreatedAt(time.Now().UTC().Unix())
	secretOrg.SetCreatedBy("octocat")
	secretOrg.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretOrg.SetUpdatedBy("octokitty")

	secretRepo := new(api.Secret)
	secretRepo.SetID(2)
	secretRepo.SetOrg("github")
	secretRepo.SetRepo("octocat")
	secretRepo.SetTeam("")
	secretRepo.SetName("foo")
	secretRepo.SetValue("bar")
	secretRepo.SetType("repo")
	secretRepo.SetImages([]string{"alpine"})
	secretRepo.SetAllowEvents(api.NewEventsFromMask(1))
	secretRepo.SetAllowCommand(true)
	secretRepo.SetAllowSubstitution(true)
	secretRepo.SetCreatedAt(time.Now().UTC().Unix())
	secretRepo.SetCreatedBy("octocat")
	secretRepo.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretRepo.SetUpdatedBy("octokitty")

	secretShared := new(api.Secret)
	secretShared.SetID(3)
	secretShared.SetOrg("github")
	secretShared.SetRepo("")
	secretShared.SetTeam("octocat")
	secretShared.SetName("foo")
	secretShared.SetValue("bar")
	secretShared.SetType("shared")
	secretShared.SetImages([]string{"alpine"})
	secretShared.SetAllowCommand(true)
	secretShared.SetAllowSubstitution(true)
	secretShared.SetAllowEvents(api.NewEventsFromMask(1))
	secretShared.SetCreatedAt(time.Now().UTC().Unix())
	secretShared.SetCreatedBy("octocat")
	secretShared.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	secretShared.SetUpdatedBy("octokitty")

	serviceOne := new(api.Service)
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

	serviceTwo := new(api.Service)
	serviceTwo.SetID(2)
	serviceTwo.SetBuildID(1)
	serviceTwo.SetRepoID(1)
	serviceTwo.SetNumber(2)
	serviceTwo.SetName("clone")
	serviceTwo.SetImage("target/vela-git-slim:v0.12.0")
	serviceTwo.SetStatus("pending")
	serviceTwo.SetError("")
	serviceTwo.SetExitCode(0)
	serviceTwo.SetCreated(1563474086)
	serviceTwo.SetStarted(1563474088)
	serviceTwo.SetFinished(1563474089)
	serviceTwo.SetHost("example.company.com")
	serviceTwo.SetRuntime("docker")
	serviceTwo.SetDistribution("linux")

	stepOne := new(api.Step)
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
	stepOne.SetReportAs("")

	stepTwo := new(api.Step)
	stepTwo.SetID(2)
	stepTwo.SetBuildID(1)
	stepTwo.SetRepoID(1)
	stepTwo.SetNumber(2)
	stepTwo.SetName("clone")
	stepTwo.SetImage("target/vela-git-slim:v0.12.0")
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
	stepTwo.SetReportAs("test")

	_bPartialOne := new(api.Build)
	_bPartialOne.SetID(1)

	_bPartialTwo := new(api.Build)
	_bPartialTwo.SetID(2)

	workerOne := new(api.Worker)
	workerOne.SetID(1)
	workerOne.SetHostname("worker-1.example.com")
	workerOne.SetAddress("https://worker-1.example.com")
	workerOne.SetRoutes([]string{"vela"})
	workerOne.SetActive(true)
	workerOne.SetStatus("available")
	workerOne.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	workerOne.SetRunningBuilds([]*api.Build{_bPartialOne})
	workerOne.SetLastBuildStartedAt(time.Now().UTC().Unix())
	workerOne.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	workerOne.SetLastCheckedIn(time.Now().UTC().Unix() - 60)
	workerOne.SetBuildLimit(1)

	workerTwo := new(api.Worker)
	workerTwo.SetID(2)
	workerTwo.SetHostname("worker-2.example.com")
	workerTwo.SetAddress("https://worker-2.example.com")
	workerTwo.SetRoutes([]string{"vela"})
	workerTwo.SetActive(true)
	workerTwo.SetStatus("available")
	workerTwo.SetLastStatusUpdateAt(time.Now().UTC().Unix())
	workerTwo.SetRunningBuilds([]*api.Build{_bPartialTwo})
	workerTwo.SetLastBuildStartedAt(time.Now().UTC().Unix())
	workerTwo.SetLastBuildFinishedAt(time.Now().UTC().Unix())
	workerTwo.SetLastCheckedIn(time.Now().UTC().Unix() - 60)
	workerTwo.SetBuildLimit(1)

	return &Resources{
		Builds:      []*api.Build{buildOne, buildTwo},
		Dashboards:  []*api.Dashboard{dashboardOne, dashboardTwo},
		Deployments: []*api.Deployment{deploymentOne, deploymentTwo},
		Executables: []*api.BuildExecutable{executableOne, executableTwo},
		Hooks:       []*api.Hook{hookOne, hookTwo, hookThree},
		JWKs:        jwkSet,
		Logs:        []*api.Log{logServiceOne, logServiceTwo, logStepOne, logStepTwo},
		Pipelines:   []*api.Pipeline{pipelineOne, pipelineTwo},
		Repos:       []*api.Repo{repoOne, repoTwo},
		Schedules:   []*api.Schedule{scheduleOne, scheduleTwo},
		Secrets:     []*api.Secret{secretOrg, secretRepo, secretShared},
		Services:    []*api.Service{serviceOne, serviceTwo},
		Steps:       []*api.Step{stepOne, stepTwo},
		Users:       []*api.User{userOne, userTwo},
		Workers:     []*api.Worker{workerOne, workerTwo},
	}
}

// CmpOptApproxUpdatedAt is a custom comparator for cmp.Equal
// to reduce flakiness in tests when comparing structs with UpdatedAt field.
func CmpOptApproxUpdatedAt() cmp.Option {
	// Custom Comparer for *int64 fields typically used to store unix timestamps.
	// Will consider time difference of 5s to be equal for sake of tests.
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
