// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// ListRepos represents a query to
	// list all repos in the database.
	ListRepos = `
SELECT *
FROM repos;
`

	// ListUserRepos represents a query to list
	// all repos for a user_id in the database.
	ListUserRepos = `
SELECT *
FROM repos
WHERE user_id = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectRepo represents a query to select a
	// repo for an org and name in the database.
	SelectRepo = `
SELECT *
FROM repos
WHERE org = ?
AND name = ?
LIMIT 1;
`

	// SelectUserReposCount represents a query to select
	// the count of repos for a user_id in the database.
	SelectUserReposCount = `
SELECT count(*) as count
FROM repos
WHERE user_id = ?;
`

	// SelectReposCount represents a query to select
	// the count of repos in the database.
	SelectReposCount = `
SELECT count(*) as count
FROM repos;
`

	// DeleteRepo represents a query to
	// remove a repo from the database.
	DeleteRepo = `
DELETE
FROM repos
WHERE id = ?;
`
)

// createRepoService is a helper function to return
// a service for interacting with the repos table.
func createRepoService() *Service {
	return &Service{
		List: map[string]string{
			"all":  ListRepos,
			"user": ListUserRepos,
		},
		Select: map[string]string{
			"repo":        SelectRepo,
			"count":       SelectReposCount,
			"countByUser": SelectUserReposCount,
		},
		Delete: DeleteRepo,
	}
}
