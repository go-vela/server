// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

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

	// SelectRepoByID represents a query to select a
	// repo for an id in the database.
	SelectRepoByID = `
SELECT *
FROM repos
WHERE id = ?
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

	// ListReposByLastUpdate represents a query to list
	// all repos in an org, ordered by latest activity.
	// In this case, latest activity is synonymous with
	// the created timestamp of the last build for the repo.
	ListReposByLastUpdate = `
SELECT r.*
FROM repos r LEFT JOIN (
	SELECT repos.id, MAX(builds.created) as latest_build
	FROM builds INNER JOIN repos
	ON builds.repo_id = repos.id
	WHERE repos.org = ?
	GROUP BY repos.id) t
ON r.id = t.id
ORDER BY latest_build DESC NULLS LAST
LIMIT ?
OFFSET ?;
`
)
