// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres dashboards table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
dashboards (
	id            UUID PRIMARY KEY,
	name          VARCHAR(250),
	created_at    BIGINT,
	created_by    VARCHAR(250),
	updated_at    BIGINT,
	updated_by    VARCHAR(250),
	admins        JSON DEFAULT NULL,
	repos         JSON DEFAULT NULL
);
`

	// CreateSqliteTable represents a query to create the Sqlite dashboards table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
dashboards (
	id            TEXT PRIMARY KEY,
	name          TEXT,
	created_at    INTEGER,
	created_by	  TEXT,
	updated_at    INTEGER,
	updated_by    TEXT,
	admins        TEXT,
	repos         TEXT
);
`
)

// CreateDashboardTable creates the dashboards table in the database.
func (e *Engine) CreateDashboardTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating dashboards table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the dashboards table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the dashboards table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
