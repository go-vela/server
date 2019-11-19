// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/library"
)

// Service represents the interface for Vela integrating
// with the different supported Database backends.
type Service interface {
	// Build Database Interface Functions

	// GetBuild defines a function that
	// gets a build by number and repo ID.
	GetBuild(int, *library.Repo) (*library.Build, error)
	// GetLastBuild defines a function that
	// gets the last build ran by repo ID.
	GetLastBuild(*library.Repo) (*library.Build, error)
	// GetBuildCount defines a function that
	// gets the count of builds.
	GetBuildCount() (int64, error)
	// GetBuildCountByStatus defines a function that
	// gets a the count of builds by status.
	GetBuildCountByStatus(string) (int64, error)
	// GetBuildList defines a function that gets
	// a list of all builds.
	GetBuildList() ([]*library.Build, error)
	// GetRepoBuildList defines a function that
	// gets a list of builds by repo ID.
	GetRepoBuildList(*library.Repo, int, int) ([]*library.Build, error)
	// GetRepoBuildCount defines a function that
	// gets the count of builds by repo ID.
	GetRepoBuildCount(*library.Repo) (int64, error)
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
	GetTypeSecretList(string, string, string, int, int) ([]*library.Secret, error)
	// GetTypeSecretCount defines a function that gets a count
	// of secrets by type, owner, and name (repo or team).
	GetTypeSecretCount(string, string, string) (int64, error)
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
}
