// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	// ListBuilds represents a query to
	// list all builds in the database.
	ListBuilds = `
SELECT *
FROM builds;
`

	// SelectBuildByID represents a query to select
	// a build for its id in the database
	SelectBuildByID = `
SELECT *
FROM builds
WHERE id = ?
LIMIT 1;
`

	// SelectRepoBuild represents a query to select
	// a build for a repo_id in the database.
	SelectRepoBuild = `
SELECT *
FROM builds
WHERE repo_id = ?
AND number = ?
LIMIT 1;
`

	// SelectLastRepoBuild represents a query to select
	// the last build for a repo_id in the database.
	SelectLastRepoBuild = `
SELECT *
FROM builds
WHERE repo_id = ?
ORDER BY number DESC
LIMIT 1;
`

	// SelectLastRepoBuildByBranch represents a query to
	// select the last build for a repo_id and branch name
	// in the database.
	SelectLastRepoBuildByBranch = `
SELECT *
FROM builds
WHERE repo_id = ?
AND branch = ?
ORDER BY number DESC
LIMIT 1;
`

	// SelectBuildsCount represents a query to select
	// the count of builds in the database.
	SelectBuildsCount = `
SELECT count(*) as count
FROM builds;
`

	// SelectBuildsCountByStatus represents a query to select
	// the count of builds for a status in the database.
	SelectBuildsCountByStatus = `
SELECT count(*) as count
FROM builds
WHERE status = ?;
`

	// DeleteBuild represents a query to
	// remove a build from the database.
	DeleteBuild = `
DELETE
FROM builds
WHERE id = ?;
`

	// SelectPendingAndRunningBuilds represents a joined query
	// between the builds & repos table to select
	// the created builds that are in pending or running builds status
	// since the specified timeframe.
	SelectPendingAndRunningBuilds = `
SELECT builds.created, builds.number, builds.status, repos.full_name
FROM builds INNER JOIN repos
ON builds.repo_id = repos.id
WHERE builds.created > ?
AND (builds.status = 'running' OR builds.status = 'pending');
`
)
