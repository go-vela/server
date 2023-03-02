// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres inits table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
initsteps (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	build_id     INTEGER,
	number       INTEGER,
	name         VARCHAR(250),
	mimetype     VARCHAR(250),
	reporter     VARCHAR(250),
	UNIQUE(build_id, number)
);
`

	// CreateSqliteTable represents a query to create the Sqlite inits table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
initsteps (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id      INTEGER,
	build_id     INTEGER,
	number       INTEGER,
	name         TEXT,
	mimetype     TEXT,
	reporter     TEXT,
	UNIQUE(build_id, number)
);
`
)

// CreateInitStepTable creates the inits table in the database.
func (e *engine) CreateInitStepTable(driver string) error {
	e.logger.Tracef("creating initsteps table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the inits table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the inits table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
