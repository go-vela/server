// SPDX-License-Identifier: Apache-2.0

package step

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres steps table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
steps (
	id            SERIAL PRIMARY KEY,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          VARCHAR(250),
	image         VARCHAR(500),
	stage         VARCHAR(250),
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

	// CreateSqliteTable represents a query to create the Sqlite steps table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
steps (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          TEXT,
	image         TEXT,
	stage         TEXT,
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

// CreateStepTable creates the steps table in the database.
func (e *engine) CreateStepTable(driver string) error {
	e.logger.Tracef("creating steps table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the steps table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the steps table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
