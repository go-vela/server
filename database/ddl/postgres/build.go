// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// CreateBuildTable represents a query to
	// create the builds table for Vela.
	CreateBuildTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id            SERIAL PRIMARY KEY,
	repo_id       INTEGER,
	number        INTEGER,
	parent        INTEGER,
	event         VARCHAR(250),
	status        VARCHAR(250),
	error         VARCHAR(500),
	enqueued      INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	deploy        VARCHAR(500),
	clone         VARCHAR(1000),
	source        VARCHAR(1000),
	title         VARCHAR(1000),
	message       VARCHAR(2000),
	commit        VARCHAR(500),
	sender        VARCHAR(250),
	author        VARCHAR(250),
	email         VARCHAR(500),
	link          VARCHAR(1000),
	branch        VARCHAR(500),
	ref           VARCHAR(500),
	base_ref      VARCHAR(500),
	host          VARCHAR(250),
	runtime       VARCHAR(250),
	distribution  VARCHAR(250),
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
