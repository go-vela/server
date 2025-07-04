// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// BuildInterface represents the Vela interface for build
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type BuildInterface interface {
	// Build Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateBuildIndexes defines a function that creates the indexes for the builds table.
	CreateBuildIndexes(context.Context) error
	// CreateBuildTable defines a function that creates the builds table.
	CreateBuildTable(context.Context, string) error

	// Build Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CleanBuilds defines a function that sets pending or running builds to error created before a given time.
	CleanBuilds(context.Context, string, int64) (int64, error)
	// CountBuilds defines a function that gets the count of all builds.
	CountBuilds(context.Context) (int64, error)
	// CountBuildsForDeployment defines a function that gets the count of builds by deployment url.
	CountBuildsForDeployment(context.Context, *api.Deployment, map[string]interface{}) (int64, error)
	// CountBuildsForOrg defines a function that gets the count of builds by org name.
	CountBuildsForOrg(context.Context, string, map[string]interface{}) (int64, error)
	// CountBuildsForRepo defines a function that gets the count of builds by repo ID.
	CountBuildsForRepo(context.Context, *api.Repo, map[string]interface{}, int64, int64) (int64, error)
	// CountBuildsForStatus defines a function that gets the count of builds by status.
	CountBuildsForStatus(context.Context, string, map[string]interface{}) (int64, error)
	// CreateBuild defines a function that creates a new build.
	CreateBuild(context.Context, *api.Build) (*api.Build, error)
	// DeleteBuild defines a function that deletes an existing build.
	DeleteBuild(context.Context, *api.Build) error
	// GetBuild defines a function that gets a build by ID.
	GetBuild(context.Context, int64) (*api.Build, error)
	// GetBuildForRepo defines a function that gets a build by repo ID and number.
	GetBuildForRepo(context.Context, *api.Repo, int64) (*api.Build, error)
	// LastBuildForRepo defines a function that gets the last build ran by repo ID and branch.
	LastBuildForRepo(context.Context, *api.Repo, string) (*api.Build, error)
	// ListBuilds defines a function that gets a list of all builds.
	ListBuilds(context.Context) ([]*api.Build, error)
	// ListBuildsForOrg defines a function that gets a list of builds by org name.
	ListBuildsForOrg(context.Context, string, map[string]any, map[string]any, int, int) ([]*api.Build, error)
	// ListBuildsForDashboardRepo defines a function that gets a list of builds based on dashboard filters.
	ListBuildsForDashboardRepo(context.Context, *api.Repo, []string, []string) ([]*api.Build, error)
	// ListBuildsForRepo defines a function that gets a list of builds by repo ID.
	ListBuildsForRepo(context.Context, *api.Repo, map[string]interface{}, int64, int64, int, int) ([]*api.Build, error)
	// ListPendingAndRunningBuilds defines a function that gets a list of pending and running builds.
	ListPendingAndRunningBuilds(context.Context, string) ([]*api.QueueBuild, error)
	// ListPendingAndRunningBuildsForRepo defines a function that gets a list of pending and running builds for a repo.
	ListPendingAndRunningBuildsForRepo(context.Context, *api.Repo) ([]*api.Build, error)
	// ListPendingApprovalBuilds defines a function that gets a list of pending approval builds that were created before a given time.
	ListPendingApprovalBuilds(context.Context, string) ([]*api.Build, error)
	// UpdateBuild defines a function that updates an existing build.
	UpdateBuild(context.Context, *api.Build) (*api.Build, error)
}
