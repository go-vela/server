// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// RepoInterface represents the Vela interface for repo
// functions with the supported Database backends.
//

type RepoInterface interface {
	// Repo Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateRepoIndexes defines a function that creates the indexes for the repos table.
	CreateRepoIndexes(context.Context) error
	// CreateRepoTable defines a function that creates the repos table.
	CreateRepoTable(context.Context, string) error

	// Repo Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountRepos defines a function that gets the count of all repos.
	CountRepos(context.Context) (int64, error)
	// CountReposForOrg defines a function that gets the count of repos by org name.
	CountReposForOrg(context.Context, string, map[string]interface{}) (int64, error)
	// CountReposForUser defines a function that gets the count of repos by user ID.
	CountReposForUser(context.Context, *api.User, map[string]interface{}) (int64, error)
	// CreateRepo defines a function that creates a new repo.
	CreateRepo(context.Context, *api.Repo) (*api.Repo, error)
	// DeleteRepo defines a function that deletes an existing repo.
	DeleteRepo(context.Context, *api.Repo) error
	// GetRepo defines a function that gets a repo by ID.
	GetRepo(context.Context, int64) (*api.Repo, error)
	// GetRepoForOrg defines a function that gets a repo by org and repo name.
	GetRepoForOrg(context.Context, string, string) (*api.Repo, error)
	// GetReposInList defines a function that gets a list of repos from a list of full names.
	GetReposInList(context.Context, []string) ([]*api.Repo, error)
	// ListRepos defines a function that gets a list of all repos.
	ListRepos(context.Context) ([]*api.Repo, error)
	// ListReposForOrg defines a function that gets a list of repos by org name.
	ListReposForOrg(context.Context, string, string, map[string]interface{}, int, int) ([]*api.Repo, error)
	// ListReposForUser defines a function that gets a list of repos by user ID.
	ListReposForUser(context.Context, *api.User, string, map[string]interface{}, int, int) ([]*api.Repo, error)
	// UpdateRepo defines a function that updates an existing repo.
	UpdateRepo(context.Context, *api.Repo) (*api.Repo, error)
}
