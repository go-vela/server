// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres services table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
services (
	id            BIGSERIAL PRIMARY KEY,
	repo_id       BIGINT,
	build_id      BIGINT,
	number        INTEGER,
	name          VARCHAR(250),
	image         VARCHAR(500),
	status        VARCHAR(250),
	error         VARCHAR(500),
	exit_code     INTEGER,
	created       BIGINT,
	started       BIGINT,
	finished      BIGINT,
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
func (e *Engine) CreateServiceTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating services table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the services table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the services table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
