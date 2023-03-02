// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres logs table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            SERIAL PRIMARY KEY,
	build_id      INTEGER,
	repo_id       INTEGER,
	service_id    INTEGER,
	step_id       INTEGER,
	initstep_id   INTEGER,
	data          BYTEA,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite logs table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id      INTEGER,
	repo_id       INTEGER,
	service_id    INTEGER,
	step_id       INTEGER,
	initstep_id   INTEGER,
	data          BLOB,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`
)

// CreateLogTable creates the logs table in the database.
func (e *engine) CreateLogTable(driver string) error {
	e.logger.Tracef("creating logs table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the logs table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the logs table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
