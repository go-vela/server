// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres settings table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
settings (
	id                  SERIAL PRIMARY KEY,
	compiler            JSON DEFAULT NULL,
	queue               JSON DEFAULT NULL,
	repo_allowlist      VARCHAR(1000),
	schedule_allowlist  VARCHAR(1000),
	max_dashboard_repos INTEGER,
	created_at          BIGINT,
	updated_at          BIGINT,
	updated_by          VARCHAR(250)
);
`

	// CreateSqliteTable represents a query to create the Sqlite settings table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
settings (
	id                      INTEGER PRIMARY KEY AUTOINCREMENT,
	compiler                TEXT,
	queue                   TEXT,
	repo_allowlist          VARCHAR(1000),
	schedule_allowlist      VARCHAR(1000),
	max_dashboard_repos     INTEGER,
	created_at              INTEGER,
	updated_at              INTEGER,
	updated_by              TEXT
);
`
)

// CreateSettingsTable creates the settings table in the database.
func (e *Engine) CreateSettingsTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating settings table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the steps table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the steps table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
