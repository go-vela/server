// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import "github.com/go-vela/types/constants"

const (
	// CreatePostgresTable represents a query to create the Postgres repos table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id            SERIAL PRIMARY KEY,
	user_id       INTEGER,
	hash          VARCHAR(500),
	org           VARCHAR(250),
	name          VARCHAR(250),
	full_name     VARCHAR(500),
	link          VARCHAR(1000),
	clone         VARCHAR(1000),
	branch        VARCHAR(250),
	topics        VARCHAR(1020),
	build_limit   INTEGER,
	timeout       INTEGER,
	counter       INTEGER,
	visibility    TEXT,
	private       BOOLEAN,
	trusted       BOOLEAN,
	active        BOOLEAN,
	allow_events  INTEGER,
	pipeline_type TEXT,
	previous_name VARCHAR(100),
	UNIQUE(full_name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite repos table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id       INTEGER,
	hash          TEXT,
	org           TEXT,
	name          TEXT,
	full_name     TEXT,
	link          TEXT,
	clone         TEXT,
	branch        TEXT,
	topics        TEXT,
	build_limit   INTEGER,
	timeout       INTEGER,
	counter       INTEGER,
	visibility    TEXT,
	private       BOOLEAN,
	trusted       BOOLEAN,
	active        BOOLEAN,
	allow_events  INTEGER,
	pipeline_type TEXT,
	previous_name TEXT,
	UNIQUE(full_name)
);
`
)

// CreateRepoTable creates the repos table in the database.
func (e *engine) CreateRepoTable(driver string) error {
	e.logger.Tracef("creating repos table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the repos table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the repos table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
