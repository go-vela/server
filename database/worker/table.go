// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres workers table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id                     BIGSERIAL PRIMARY KEY,
	hostname               VARCHAR(250),
	address                VARCHAR(250),
	routes                 VARCHAR(1000),
	active                 BOOLEAN,
	status                 VARCHAR(50),
	last_status_update_at  BIGINT,
	running_build_ids      VARCHAR(500),
	last_build_started_at  BIGINT,
	last_build_finished_at BIGINT,
	last_checked_in        BIGINT,
	build_limit            INTEGER,
	UNIQUE(hostname)
);
`
	// CreateSqliteTable represents a query to create the Sqlite workers table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
workers (
	id                     INTEGER PRIMARY KEY AUTOINCREMENT,
	hostname               TEXT,
	address                TEXT,
	routes                 TEXT,
	active                 BOOLEAN,
	status                 VARCHAR(50),
	last_status_update_at  INTEGER,
	running_build_ids      VARCHAR(500),
	last_build_started_at  INTEGER,
	last_build_finished_at INTEGER,
	last_checked_in	       INTEGER,
	build_limit            INTEGER,
	UNIQUE(hostname)
);
`
)

// CreateWorkerTable creates the workers table in the database.
func (e *Engine) CreateWorkerTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating workers table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the workers table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the workers table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
