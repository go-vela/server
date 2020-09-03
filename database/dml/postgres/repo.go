// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

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
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
`

	// ListOrgRepos represents a query to list
	// all repos for an org in the database.
	ListOrgRepos = `
SELECT *
FROM repos
WHERE org = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
`

	// SelectRepo represents a query to select a
	// repo for an org and name in the database.

	//[here] This is the query for badges.
	SelectRepo = `
SELECT *
FROM repos
WHERE org = $1
AND name = $2
LIMIT 1;
`
	//[here] Note: Rename this query in the future to avoid confusion.
	SelectRepoOrg = `
SELECT *
FROM repos
WHERE org = $1
LIMIT 1;
	`

	// SelectUserReposCount represents a query to select
	// the count of repos for a user_id in the database.
	SelectUserReposCount = `
SELECT count(*) as count
FROM repos
WHERE user_id = $1;
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
WHERE id = $1;
`
)

// createRepoService is a helper function to return
// a service for interacting with the repos table.
func createRepoService() *Service {
	return &Service{
		List: map[string]string{
			"all":  ListRepos,
			"user": ListUserRepos,
			"org":  ListOrgRepos,
		},
		Select: map[string]string{
			"repo":        SelectRepo,
			"count":       SelectReposCount,
			"countByUser": SelectUserReposCount,
			"rename":      SelectRepoOrg,
		},
		Delete: DeleteRepo,
	}
}
