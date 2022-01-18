// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateBuildTable represents a query to
	// create the builds table for Vela.
	CreateBuildTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id             SERIAL PRIMARY KEY,
	repo_id        INTEGER,
	number         INTEGER,
	parent         INTEGER,
	event          VARCHAR(250),
	status         VARCHAR(250),
	error          VARCHAR(500),
	enqueued       INTEGER,
	created        INTEGER,
	started        INTEGER,
	finished       INTEGER,
	deploy         VARCHAR(500),
	deploy_payload VARCHAR(2000),
	clone          VARCHAR(1000),
	source         VARCHAR(1000),
	title          VARCHAR(1000),
	message        VARCHAR(2000),
	commit         VARCHAR(500),
	sender         VARCHAR(250),
	author         VARCHAR(250),
	email          VARCHAR(500),
	link           VARCHAR(1000),
	branch         VARCHAR(500),
	ref            VARCHAR(500),
	base_ref       VARCHAR(500),
	head_ref       VARCHAR(500),
	host           VARCHAR(250),
	runtime        VARCHAR(250),
	distribution   VARCHAR(250),
	timestamp      INTEGER,
	UNIQUE(repo_id, number)
);
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

	// CreateBuildCreatedIndex represents a query to create an
	// index on the builds table for the created column.
	CreateBuildCreatedIndex = `
CREATE INDEX CONCURRENTLY
IF NOT EXISTS
builds_created
ON builds (created);
`
)
