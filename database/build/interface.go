// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
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
	CreateBuildIndexes() error
	// CreateBuildTable defines a function that creates the builds table.
	CreateBuildTable(string) error

	// Build Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CleanBuilds defines a function that sets pending or running builds to error created before a given time.
	CleanBuilds(string, int64) (int64, error)
	// CountBuilds defines a function that gets the count of all builds.
	CountBuilds() (int64, error)
	// CountBuildsForDeployment defines a function that gets the count of builds by deployment url.
	CountBuildsForDeployment(*library.Deployment, map[string]interface{}) (int64, error)
	// CountBuildsForOrg defines a function that gets the count of builds by org name.
	CountBuildsForOrg(string, map[string]interface{}) (int64, error)
	// CountBuildsForRepo defines a function that gets the count of builds by repo ID.
	CountBuildsForRepo(*library.Repo, map[string]interface{}) (int64, error)
	// CountBuildsForStatus defines a function that gets the count of builds by status.
	CountBuildsForStatus(string, map[string]interface{}) (int64, error)
	// CreateBuild defines a function that creates a new build.
	CreateBuild(*library.Build) error
	// DeleteBuild defines a function that deletes an existing build.
	DeleteBuild(*library.Build) error
	// GetBuild defines a function that gets a build by ID.
	GetBuild(int64) (*library.Build, error)
	// GetBuildForRepo defines a function that gets a build by repo ID and number.
	GetBuildForRepo(*library.Repo, int) (*library.Build, error)
	// LastBuildForRepo defines a function that gets the last build ran by repo ID and branch.
	LastBuildForRepo(*library.Repo, string) (*library.Build, error)
	// ListBuilds defines a function that gets a list of all builds.
	ListBuilds() ([]*library.Build, error)
	// ListBuildsForDeployment defines a function that gets a list of builds by deployment url.
	ListBuildsForDeployment(*library.Deployment, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// ListBuildsForOrg defines a function that gets a list of builds by org name.
	ListBuildsForOrg(string, map[string]interface{}, int, int) ([]*library.Build, int64, error)
	// ListBuildsForRepo defines a function that gets a list of builds by repo ID.
	ListBuildsForRepo(*library.Repo, map[string]interface{}, int64, int64, int, int) ([]*library.Build, int64, error)
	// ListPendingAndRunningBuilds defines a function that gets a list of pending and running builds.
	ListPendingAndRunningBuilds(string) ([]*library.BuildQueue, error)
	// UpdateBuild defines a function that updates an existing build.
	UpdateBuild(*library.Build) error
}
