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
WHERE repo_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
`

	// ListRepoBuildsByEvent represents a query to select
	// a build for a repo_id with a specific event type
	// in the database.
	ListRepoBuildsByEvent = `
SELECT *
FROM builds
WHERE repo_id = $1
AND event = $2
ORDER BY number DESC
LIMIT $3
OFFSET $4;
`
	// ListOrgBuildsByEvent represents a joined query
	// between the builds & repos table to select
	// a build for an org with a specific event type
	// in the database.
	ListOrgBuildsByEvent = `
SELECT builds.* 
FROM builds JOIN repos 
ON repos.id=builds.repo_id 
WHERE repos.org = $1
AND builds.event = $2
ORDER BY id DESC
LIMIT $3
OFFSET $4;
`

	// SelectRepoBuild represents a query to select
	// a build for a repo_id in the database.
	SelectRepoBuild = `
SELECT *
FROM builds
WHERE repo_id = $1
AND number = $2
LIMIT 1;
`

	// SelectLastRepoBuild represents a query to select
	// the last build for a repo_id in the database.
	SelectLastRepoBuild = `
SELECT *
FROM builds
WHERE repo_id = $1
ORDER BY number DESC
LIMIT 1;
`
	// ListOrgBuilds represents a joined query
	// between the builds & repos table to select
	// the last build for a org name in the database.
	ListOrgBuilds = `
SELECT builds.*
FROM builds JOIN repos
ON repos.id=builds.repo_id 
WHERE repos.org = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
		`

	// SelectLastRepoBuildByBranch represents a query to
	// select the last build for a repo_id and branch name
	// in the database.
	SelectLastRepoBuildByBranch = `
SELECT *
FROM builds
WHERE repo_id = $1
AND branch = $2
ORDER BY number DESC
LIMIT 1;
`

	// SelectBuildsCount represents a query to select
	// the count of builds in the database.
	SelectBuildsCount = `
SELECT count(*) as count
FROM builds;
`

	// SelectRepoBuildCount represents a query to select
	// the count of builds for a repo_id in the database.
	SelectRepoBuildCount = `
SELECT count(*) as count
FROM builds
WHERE repo_id = $1;
`
	// SelectOrgBuildCount represents a joined query
	// between the builds & repos table to select
	// the count of builds for an org name in the database.
	SelectOrgBuildCount = `
SELECT count(*) as count
FROM builds JOIN repos
ON repos.id = builds.repo_id 
WHERE repos.org = $1;
`
	// SelectRepoBuildCountByEvent represents a query to select
	// the count of builds for by repo and event type in the database.
	SelectRepoBuildCountByEvent = `
SELECT count(*) as count
FROM builds
WHERE repo_id = $1
AND event = $2;
`

	// SelectOrgBuildCountByEvent represents a joined query
	// between the builds & repos table to select
	// the count of builds for by org name and event type in the database.
	SelectOrgBuildCountByEvent = `
SELECT count(*) as count
FROM builds JOIN repos
ON repos.id = builds.repo_id 
WHERE repos.org = $1
AND event = $2;
`

	// SelectBuildsCountByStatus represents a query to select
	// the count of builds for a status in the database.
	SelectBuildsCountByStatus = `
SELECT count(*) as count
FROM builds
WHERE status = $1;
`

	// DeleteBuild represents a query to
	// remove a build from the database.
	DeleteBuild = `
DELETE
FROM builds
WHERE id = $1;
`

	// SelectPendingAndRunningBuilds represents a joined query
	// between the builds & repos table to select
	// the created builds that are in pending or running builds status
	// since the specified timeframe.
	SelectPendingAndRunningBuilds = `
SELECT builds.created, builds.number, builds.status, repos.full_name
FROM builds INNER JOIN repos ON (builds.repo_id = repos.id)
WHERE builds.created > $1
AND builds.status = 'running' or builds.status = 'pending';
`
)
