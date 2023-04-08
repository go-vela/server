// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/log"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/step"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
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
	// GetBuildByID defines a function that
	// gets a build by its id.
	GetBuildByID(int64) (*library.Build, error)
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
	GetRepoBuildList(*library.Repo, map[string]interface{}, int64, int64, int, int) ([]*library.Build, int64, error)
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

	// HookService provides the interface for functionality
	// related to hooks stored in the database.
	hook.HookService

	// LogService provides the interface for functionality
	// related to logs stored in the database.
	log.LogService

	// PipelineService provides the interface for functionality
	// related to pipelines stored in the database.
	pipeline.PipelineService

	// RepoService provides the interface for functionality
	// related to repos stored in the database.
	repo.RepoService

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

	// StepService provides the interface for functionality
	// related to steps stored in the database.
	step.StepService

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

	// UserService provides the interface for functionality
	// related to users stored in the database.
	user.UserService

	// WorkerService provides the interface for functionality
	// related to workers stored in the database.
	worker.WorkerService
}
