// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres builds table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id             SERIAL PRIMARY KEY,
	repo_id        INTEGER,
	pipeline_id    INTEGER,
	number         INTEGER,
	parent         INTEGER,
	event          VARCHAR(250),
	event_action   VARCHAR(250),
	status         VARCHAR(250),
	error          VARCHAR(1000),
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

	// CreateSqliteTable represents a query to create the Sqlite builds table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id        INTEGER,
	pipeline_id    INTEGER,
	number         INTEGER,
	parent         INTEGER,
	event          TEXT,
	event_action   TEXT,
	status         TEXT,
	error          TEXT,
	enqueued       INTEGER,
	created        INTEGER,
	started        INTEGER,
	finished       INTEGER,
	deploy         TEXT,
	deploy_payload TEXT,
	clone          TEXT,
	source         TEXT,
	title          TEXT,
	message        TEXT,
	'commit'       TEXT,
	sender         TEXT,
	author         TEXT,
	email          TEXT,
	link           TEXT,
	branch         TEXT,
	ref            TEXT,
	base_ref       TEXT,
	head_ref       TEXT,
	host           TEXT,
	runtime        TEXT,
	distribution   TEXT,
	timestamp      INTEGER,
	UNIQUE(repo_id, number)
);
`
)

// CreateBuildTable creates the builds table in the database.
func (e *engine) CreateBuildTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating builds table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the builds table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the builds table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
