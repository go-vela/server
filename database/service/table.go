// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres services table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
services (
	id            SERIAL PRIMARY KEY,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          VARCHAR(250),
	image         VARCHAR(500),
	status        VARCHAR(250),
	error         VARCHAR(500),
	exit_code     INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	host          VARCHAR(250),
	runtime       VARCHAR(250),
	distribution  VARCHAR(250),
	UNIQUE(build_id, number)
);
`

	// CreateSqliteTable represents a query to create the Sqlite services table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
services (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          TEXT,
	image         TEXT,
	status        TEXT,
	error         TEXT,
	exit_code     INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	host          TEXT,
	runtime       TEXT,
	distribution  TEXT,
	UNIQUE(build_id, number)
);
`
)

// CreateServiceTable creates the services table in the database.
func (e *engine) CreateServiceTable(driver string) error {
	e.logger.Tracef("creating services table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the services table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the services table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
