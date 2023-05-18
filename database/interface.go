// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
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
)

// Interface represents the interface for Vela integrating
// with the different supported Database backends.
type Interface interface {
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

	// HookInterface provides the interface for functionality
	// related to hooks stored in the database.
	hook.HookInterface

	// LogInterface provides the interface for functionality
	// related to logs stored in the database.
	log.LogInterface

	// PipelineInterface provides the interface for functionality
	// related to pipelines stored in the database.
	pipeline.PipelineInterface

	// RepoInterface provides the interface for functionality
	// related to repos stored in the database.
	repo.RepoInterface

	// ScheduleInterface provides the interface for functionality
	// related to schedules stored in the database.
	schedule.ScheduleInterface

	// SecretInterface provides the interface for functionality
	// related to secrets stored in the database.
	secret.SecretInterface

	// ServiceInterface provides the interface for functionality
	// related to services stored in the database.
	service.ServiceInterface

	// StepInterface provides the interface for functionality
	// related to steps stored in the database.
	step.StepInterface

	// UserInterface provides the interface for functionality
	// related to users stored in the database.
	user.UserInterface

	// WorkerInterface provides the interface for functionality
	// related to workers stored in the database.
	worker.WorkerInterface
}
