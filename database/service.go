// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/library"
)

// Service represents the interface for Vela integrating
// with the different supported Database backends.
type Service interface {
	// Database Interface Functions

	// Driver defines a function that outputs
	// the configured database driver.
	Driver() string

	// Build Database Interface Functions

	// GetBuild defines a function that
	// gets a build by number and repo ID.
	GetBuild(int, *library.Repo) (*library.Build, error)
	// GetLastBuild defines a function that
	// gets the last build ran by repo ID.
	GetLastBuild(*library.Repo) (*library.Build, error)
	// GetLastBuildByBranch defines a function that
	// gets the last build ran by repo ID and branch.
	GetLastBuildByBranch(*library.Repo, string) (*library.Build, error)
	// GetBuildCount defines a function that
	// gets the count of builds.
	GetBuildCount() (int64, error)
	// GetBuildCountByStatus defines a function that
	// gets a the count of builds by status.
	GetBuildCountByStatus(string) (int64, error)
	// GetBuildList defines a function that gets
	// a list of all builds.
	GetBuildList() ([]*library.Build, error)
	// GetDeploymentBuildList defines a function that gets
	// a list of builds related to a deployment.
	GetDeploymentBuildList(string) ([]*library.Build, error)
	// GetRepoBuildList defines a function that
	// gets a list of builds by repo ID.
	GetRepoBuildList(*library.Repo, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// GetOrgBuildList defines a function that
	// gets a list of builds by org.
	GetOrgBuildList(string, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// GetRepoBuildCount defines a function that
	// gets the count of builds by repo ID.
	GetRepoBuildCount(*library.Repo, map[string]interface{}) (int64, error)
	// GetOrgBuildCount defines a function that
	// gets the count of builds by org.
	GetOrgBuildCount(string, map[string]interface{}) (int64, error)
	// GetPendingAndRunningBuilds defines a function that
	// gets the list of pending and running builds.
	GetPendingAndRunningBuilds(string) ([]*library.BuildQueue, error)
	// CreateBuild defines a function that
	// creates a new build.
	CreateBuild(*library.Build) error
	// UpdateBuild defines a function that
	// updates a build.
	UpdateBuild(*library.Build) error
	// DeleteBuild defines a function that
	// deletes a build by unique ID.
	DeleteBuild(int64) error

	// Hook Database Interface Functions

	// GetHook defines a function that
	// gets a webhook by number and repo ID.
	GetHook(int, *library.Repo) (*library.Hook, error)
	// GetLastHook defines a function that
	// gets the last hook by repo ID.
	GetLastHook(*library.Repo) (*library.Hook, error)
	// GetHookList defines a function that gets
	// a list of all webhooks.
	GetHookList() ([]*library.Hook, error)
	// GetRepoHookList defines a function that
	// gets a list of webhooks by repo ID.
	GetRepoHookList(*library.Repo, int, int) ([]*library.Hook, error)
	// GetRepoHookCount defines a function that
	// gets the count of webhooks by repo ID.
	GetRepoHookCount(*library.Repo) (int64, error)
	// CreateHook defines a function that
	// creates a new webhook.
	CreateHook(*library.Hook) error
	// UpdateHook defines a function that
	// updates a webhook.
	UpdateHook(*library.Hook) error
	// DeleteHook defines a function that
	// deletes a webhook by unique ID.
	DeleteHook(int64) error

	// Log Database Interface Functions

	// GetStepLog defines a function that
	// gets a step log by unique ID.
	GetStepLog(int64) (*library.Log, error)
	// GetServiceLog defines a function that
	// gets a service log by unique ID.
	GetServiceLog(int64) (*library.Log, error)
	// GetBuildLogs defines a function that
	// gets a list of logs by build ID.
	GetBuildLogs(int64) ([]*library.Log, error)
	// CreateLog defines a function that
	// creates a new log.
	CreateLog(*library.Log) error
	// UpdateLog defines a function that
	// updates a log.
	UpdateLog(*library.Log) error
	// DeleteLog defines a function that
	// deletes a log by unique ID.
	DeleteLog(int64) error

	// Repo Database Interface Functions

	// GetRepo defines a function that
	// gets a repo by org and name.
	GetRepo(string, string) (*library.Repo, error)
	// GetRepoList defines a function that
	// gets a list of all repos.
	GetRepoList() ([]*library.Repo, error)
	// GetOrgRepoList defines a function that
	// gets a list of all repos by org excluding repos specified.
	GetOrgRepoList(string, map[string]string, int, int) ([]*library.Repo, error)
	// GetOrgRepoCount defines a function that
	// gets the count of repos for an org.
	GetOrgRepoCount(string, map[string]string) (int64, error)
	// GetRepoCount defines a function that
	// gets the count of repos.
	GetRepoCount() (int64, error)
	// GetUserRepoList defines a function
	// that gets a list of repos by user ID.
	GetUserRepoList(*library.User, int, int) ([]*library.Repo, error)
	// GetUserRepoCount defines a function that
	// gets the count of repos for a user.
	GetUserRepoCount(*library.User) (int64, error)
	// CreateRepo defines a function that
	// creates a new repo.
	CreateRepo(*library.Repo) error
	// UpdateRepo defines a function that
	// updates a repo.
	UpdateRepo(*library.Repo) error
	// DeleteRepo defines a function that
	// deletes a repo by unique ID.
	DeleteRepo(int64) error

	// Secret Database Interface Functions

	// GetSecret defines a function that gets a secret
	// by type, org, name (repo or team) and secret name.
	GetSecret(string, string, string, string) (*library.Secret, error)
	// GetSecretList defines a function that
	// gets a list of all secrets.
	GetSecretList() ([]*library.Secret, error)
	// GetTypeSecretList defines a function that gets a list
	// of secrets by type, owner, and name (repo or team).
	GetTypeSecretList(string, string, string, int, int, []string) ([]*library.Secret, error)
	// GetTypeSecretCount defines a function that gets a count
	// of secrets by type, owner, and name (repo or team).
	GetTypeSecretCount(string, string, string, []string) (int64, error)
	// CreateSecret defines a function that
	// creates a new secret.
	CreateSecret(*library.Secret) error
	// UpdateSecret defines a function that
	// updates a secret.
	UpdateSecret(*library.Secret) error
	// DeleteSecret defines a function that
	// deletes a secret by unique ID.
	DeleteSecret(int64) error

	// Step Database Interface Functions

	// GetStep defines a function that
	// gets a step by number and build ID.
	GetStep(int, *library.Build) (*library.Step, error)
	// GetStepList defines a function that
	// gets a list of all steps.
	GetStepList() ([]*library.Step, error)
	// GetBuildStepList defines a function that
	// gets a list of steps by build ID.
	GetBuildStepList(*library.Build, int, int) ([]*library.Step, error)
	// GetBuildStepCount defines a function that
	// gets the count of steps by build ID.
	GetBuildStepCount(*library.Build) (int64, error)
	// GetStepImageCount defines a function that
	// gets a list of all step images and the
	// count of their occurrence.
	GetStepImageCount() (map[string]float64, error)
	// GetStepStatusCount defines a function that
	// gets a list of all step statuses and the
	// count of their occurrence.
	GetStepStatusCount() (map[string]float64, error)
	// CreateStep defines a function that
	// creates a new step.
	CreateStep(*library.Step) error
	// UpdateStep defines a function that
	// updates a step.
	UpdateStep(*library.Step) error
	// DeleteStep defines a function that
	// deletes a step by unique ID.
	DeleteStep(int64) error

	// Service Database Interface Functions

	// GetService defines a function that
	// gets a step by number and build ID.
	GetService(int, *library.Build) (*library.Service, error)
	// GetServiceList defines a function that
	// gets a list of all steps.
	GetServiceList() ([]*library.Service, error)
	// GetBuildServiceList defines a function
	// that gets a list of steps by build ID.
	GetBuildServiceList(*library.Build, int, int) ([]*library.Service, error)
	// GetBuildServiceCount defines a function
	// that gets the count of steps by build ID.
	GetBuildServiceCount(*library.Build) (int64, error)
	// GetServiceImageCount defines a function that
	// gets a list of all service images and the
	// count of their occurrence.
	GetServiceImageCount() (map[string]float64, error)
	// GetServiceStatusCount defines a function that
	// gets a list of all service statuses and the
	// count of their occurrence.
	GetServiceStatusCount() (map[string]float64, error)
	// CreateService defines a function that
	// creates a new step.
	CreateService(*library.Service) error
	// UpdateService defines a function that
	// updates a step.
	UpdateService(*library.Service) error
	// DeleteService defines a function that
	// deletes a step by unique ID.
	DeleteService(int64) error

	// User Database Interface Functions

	// GetUser defines a function that
	// gets a user by unique ID.
	GetUser(int64) (*library.User, error)
	// GetUserName defines a function that
	// gets a user by name.
	GetUserName(string) (*library.User, error)
	// GetUserList defines a function that
	// gets a list of all users.
	GetUserList() ([]*library.User, error)
	// GetUserCount defines a function that
	// gets the count of users.
	GetUserCount() (int64, error)
	// GetUserLiteList defines a function
	// that gets a lite list of users.
	GetUserLiteList(int, int) ([]*library.User, error)
	// CreateUser defines a function that
	// creates a new user.
	CreateUser(*library.User) error
	// UpdateUser defines a function that
	// updates a user.
	UpdateUser(*library.User) error
	// DeleteUser defines a function that
	// deletes a user by unique ID.
	DeleteUser(int64) error

	// Worker Database Interface Functions

	// GetWorker defines a function that
	// gets a worker by hostname.
	GetWorker(string) (*library.Worker, error)
	// GetWorkerAddress defines a function that
	// gets a worker by address.
	GetWorkerByAddress(string) (*library.Worker, error)
	// GetWorkerList defines a function that
	// gets a list of all workers.
	GetWorkerList() ([]*library.Worker, error)
	// GetWorkerCount defines a function that
	// gets the count of workers.
	GetWorkerCount() (int64, error)
	// CreateWorker defines a function that
	// creates a new worker.
	CreateWorker(*library.Worker) error
	// UpdateWorker defines a function that
	// updates a worker by unique ID.
	UpdateWorker(*library.Worker) error
	// DeleteWorker defines a function that
	// deletes a worker by hostname.
	DeleteWorker(int64) error
}
