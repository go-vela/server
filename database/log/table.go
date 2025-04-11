// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres logs table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            BIGSERIAL PRIMARY KEY,
	build_id      BIGINT,
	repo_id       BIGINT,
	service_id    BIGINT,
	step_id       BIGINT,
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
	data          BLOB,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`
)

// CreateLogTable creates the logs table in the database.
func (e *Engine) CreateLogTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating logs table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the logs table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the logs table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
