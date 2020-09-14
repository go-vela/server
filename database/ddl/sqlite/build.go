// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// CreateBuildTable represents a query to
	// create the builds table for Vela.
	CreateBuildTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id       INTEGER,
	number        INTEGER,
	parent        INTEGER,
	event         TEXT,
	status        TEXT,
	error         TEXT,
	enqueued      INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	deploy        TEXT,
	clone         TEXT,
	source        TEXT,
	title         TEXT,
	message       TEXT,
	'commit'      TEXT,
	sender        TEXT,
	author        TEXT,
	email         TEXT,
	link          TEXT,
	branch        TEXT,
	ref           TEXT,
	base_ref      TEXT,
	head_ref      TEXT,
	host          TEXT,
	runtime       TEXT,
	distribution  TEXT,
	timestamp     INTEGER,
	UNIQUE(repo_id, number)
);
`

	// CreateBuildRepoIDNumberIndex represents a query to create an
	// index on the builds table for the repo_id and number columns.
	CreateBuildRepoIDNumberIndex = `
CREATE INDEX
IF NOT EXISTS
builds_repo_id_number
ON builds (repo_id, number);
`

	// CreateBuildRepoIDIndex represents a query to create an
	// index on the builds table for the repo_id column.
	CreateBuildRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
builds_repo_id
ON builds (repo_id);
`

	// CreateBuildStatusIndex represents a query to create an
	// index on the builds table for the status column.
	CreateBuildStatusIndex = `
CREATE INDEX
IF NOT EXISTS
builds_status
ON builds (status);
`
)

// createBuildService is a helper function to return
// a service for interacting with the builds table.
func createBuildService() *Service {
	return &Service{
		Create:  CreateBuildTable,
		Indexes: []string{CreateBuildRepoIDIndex, CreateBuildRepoIDNumberIndex, CreateBuildStatusIndex},
	}
}
