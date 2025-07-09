// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres schedules table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
schedules (
	id           BIGSERIAL PRIMARY KEY,
	repo_id      BIGINT,
	active       BOOLEAN,
	name         VARCHAR(100),
	entry        VARCHAR(100),
	created_at   BIGINT,
	created_by   VARCHAR(250),
	updated_at   BIGINT,
	updated_by   VARCHAR(250),
	scheduled_at BIGINT,
	branch       VARCHAR(250),
	error        VARCHAR(250),
	UNIQUE(repo_id, name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite schedules table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
schedules (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id      INTEGER,
	active       BOOLEAN,
	name         TEXT,
	entry        TEXT,
	created_at   INTEGER,
	created_by   TEXT,
	updated_at   INTEGER,
	updated_by   TEXT,
	scheduled_at INTEGER,
	branch       TEXT,
	error        TEXT,
	UNIQUE(repo_id, name)
);
`
)

// CreateScheduleTable creates the schedules table in the database.
func (e *Engine) CreateScheduleTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating schedules table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the schedules table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the schedules table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
