// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	// ListRepoBuilds represents a query to list
	// all builds for a repo_id in the database.
	ListRepoBuilds = `
SELECT *
FROM builds
WHERE repo_id = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// ListRepoBuildsByEvent represents a query to select
	// a build for a repo_id with a specific event type
	// in the database.
	ListRepoBuildsByEvent = `
SELECT *
FROM builds
WHERE repo_id = ?
AND event = ?
ORDER BY number DESC
LIMIT ?
OFFSET ?;
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

	// ListOrgBuilds represents a joined query
	// between the builds & repos table to select
	// the last build for a org name in the database.
	ListOrgBuilds = `
SELECT builds.*
FROM builds JOIN repos
ON builds.repo_id = repos.id
WHERE repos.org = ?
AND repo_id NOT IN (?)
ORDER BY id DESC
LIMIT ?
OFFSET ?;
		`

	// ListOrgBuildsByEvent represents a joined query
	// between the builds & repos table to select
	// a build for an org with a specific event type
	// in the database.
	ListOrgBuildsByEvent = `
SELECT builds.* 
FROM builds JOIN repos 
ON builds.repo_id = repos.id
WHERE repos.org = ?
AND repo_id NOT IN (?)
AND builds.event = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectBuildsCount represents a query to select
	// the count of builds in the database.
	SelectBuildsCount = `
SELECT count(*) as count
FROM builds;
`

	// SelectOrgBuildCount represents a joined query
	// between the builds & repos table to select
	// the count of builds for an org name in the database.
	SelectOrgBuildCount = `
SELECT count(*) as count
FROM builds JOIN repos
ON builds.repo_id = repos.id
WHERE repos.org = ?
AND repo_id NOT IN (?);
`

	// SelectOrgBuildCountByEvent represents a joined query
	// between the builds & repos table to select
	// the count of builds for by org name and event type in the database.
	SelectOrgBuildCountByEvent = `
SELECT count(*) as count
FROM builds JOIN repos
ON builds.repo_id = repos.id
WHERE repos.org = ?
AND repo_id NOT IN (?)
AND builds.event = ?;
`

	// SelectRepoBuildCount represents a query to select
	// the count of builds for a repo_id in the database.
	SelectRepoBuildCount = `
SELECT count(*) as count
FROM builds
WHERE repo_id = ?;
`

	// SelectRepoBuildCountByEvent represents a query to select
	// the count of builds for by repo and event type in the database.
	SelectRepoBuildCountByEvent = `
SELECT count(*) as count
FROM builds
WHERE repo_id = ?
AND event = ?;
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
