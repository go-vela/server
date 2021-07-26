// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	// ListOrgRepos represents a query to list
	// all repos for an org in the database.
	ListOrgRepos = `
SELECT *
FROM repos
WHERE org = ?
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
	// SelectOrgReposCount represents a query to select
	// the count of repos for a user_id in the database.
	SelectOrgReposCount = `
SELECT count(*) as count
FROM repos
WHERE org = ?;
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
