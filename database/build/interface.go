// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/library"
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
	CountBuildsForDeployment(context.Context, *library.Deployment, map[string]interface{}) (int64, error)
	// CountBuildsForOrg defines a function that gets the count of builds by org name.
	CountBuildsForOrg(context.Context, string, map[string]interface{}) (int64, error)
	// CountBuildsForRepo defines a function that gets the count of builds by repo ID.
	CountBuildsForRepo(context.Context, *library.Repo, map[string]interface{}) (int64, error)
	// CountBuildsForSender defines a function that gets the count of builds by sender.
	CountBuildsForSender(context.Context, string, map[string]interface{}) (int64, error)
	// CountBuildsForStatus defines a function that gets the count of builds by status.
	CountBuildsForStatus(context.Context, string, map[string]interface{}) (int64, error)
	// CreateBuild defines a function that creates a new build.
	CreateBuild(context.Context, *library.Build) (*library.Build, error)
	// DeleteBuild defines a function that deletes an existing build.
	DeleteBuild(context.Context, *library.Build) error
	// GetBuild defines a function that gets a build by ID.
	GetBuild(context.Context, int64) (*library.Build, error)
	// GetBuildForRepo defines a function that gets a build by repo ID and number.
	GetBuildForRepo(context.Context, *library.Repo, int) (*library.Build, error)
	// LastBuildForRepo defines a function that gets the last build ran by repo ID and branch.
	LastBuildForRepo(context.Context, *library.Repo, string) (*library.Build, error)
	// ListBuilds defines a function that gets a list of all builds.
	ListBuilds(context.Context) ([]*library.Build, error)
	// ListBuildsForDeployment defines a function that gets a list of builds by deployment url.
	ListBuildsForDeployment(context.Context, *library.Deployment, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// ListBuildsForOrg defines a function that gets a list of builds by org name.
	ListBuildsForOrg(context.Context, string, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// ListBuildsForRepo defines a function that gets a list of builds by repo ID.
	ListBuildsForRepo(context.Context, *library.Repo, map[string]interface{}, int64, int64, int, int) ([]*library.Build, int64, error)
	// ListBuildsForSender defines a function that gets a list of builds by sender.
	ListBuildsForSender(context.Context, string, map[string]interface{}, int64, int64, int, int) ([]*library.Build, int64, error)
	// ListPendingAndRunningBuilds defines a function that gets a list of pending and running builds.
	ListPendingAndRunningBuilds(context.Context, string) ([]*library.BuildQueue, error)
	// UpdateBuild defines a function that updates an existing build.
	UpdateBuild(context.Context, *library.Build) (*library.Build, error)
}
