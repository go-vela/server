// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"context"

	"github.com/go-vela/types/library"
)

// RepoInterface represents the Vela interface for repo
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
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
	CountReposForUser(context.Context, *library.User, map[string]interface{}) (int64, error)
	// CreateRepo defines a function that creates a new repo.
	CreateRepo(context.Context, *library.Repo) (*library.Repo, error)
	// DeleteRepo defines a function that deletes an existing repo.
	DeleteRepo(context.Context, *library.Repo) error
	// GetRepo defines a function that gets a repo by ID.
	GetRepo(context.Context, int64) (*library.Repo, error)
	// GetRepoForOrg defines a function that gets a repo by org and repo name.
	GetRepoForOrg(context.Context, string, string) (*library.Repo, error)
	// ListRepos defines a function that gets a list of all repos.
	ListRepos(context.Context) ([]*library.Repo, error)
	// ListReposForOrg defines a function that gets a list of repos by org name.
	ListReposForOrg(context.Context, string, string, map[string]interface{}, int, int) ([]*library.Repo, int64, error)
	// ListReposForUser defines a function that gets a list of repos by user ID.
	ListReposForUser(context.Context, *library.User, string, map[string]interface{}, int, int) ([]*library.Repo, int64, error)
	// UpdateRepo defines a function that updates an existing repo.
	UpdateRepo(context.Context, *library.Repo) (*library.Repo, error)
}
