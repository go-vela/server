// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres org_build_limits table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
org_build_limits (
	id          SERIAL PRIMARY KEY,
	org         VARCHAR(250),
	build_limit INTEGER,
	created_at  BIGINT,
	updated_at  BIGINT,
	updated_by  VARCHAR(250),
	UNIQUE(org)
);
`

	// CreateSqliteTable represents a query to create the Sqlite org_build_limits table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
org_build_limits (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	org         TEXT,
	build_limit INTEGER,
	created_at  INTEGER,
	updated_at  INTEGER,
	updated_by  TEXT,
	UNIQUE(org)
);
`
)

// CreateOrgBuildLimitTable creates the org_build_limits table in the database.
func (e *Engine) CreateOrgBuildLimitTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating org_build_limits table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the org_build_limits table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the org_build_limits table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
