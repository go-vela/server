// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

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

	// SelectRepoBuildCountByEvent represents a query to select
	// the count of builds for by repo and event type in the database.
	SelectRepoBuildCountByEvent = `
SELECT count(*) as count
FROM builds
WHERE repo_id = $1
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
)

// createBuildService is a helper function to return
// a service for interacting with the builds table.
func createBuildService() *Service {
	return &Service{
		List: map[string]string{
			"all":         ListBuilds,
			"repo":        ListRepoBuilds,
			"repoByEvent": ListRepoBuildsByEvent,
		},
		Select: map[string]string{
			"repo":                SelectRepoBuild,
			"last":                SelectLastRepoBuild,
			"count":               SelectBuildsCount,
			"countByStatus":       SelectBuildsCountByStatus,
			"countByRepo":         SelectRepoBuildCount,
			"countByRepoAndEvent": SelectRepoBuildCountByEvent,
		},
		Delete: DeleteBuild,
	}
}
