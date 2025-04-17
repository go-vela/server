// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres testreports table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
testreports (
	id             SERIAL PRIMARY KEY,
	build_id       INTEGER,
	created        INTEGER
);
`

	// CreateSqliteTable represents a query to create the Sqlite testreports table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
testreports (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id       INTEGER,
	created        INTEGER
);
`
)

// CreateTestReportsTable creates the testreports table in the database.
func (e *Engine) CreateTestReportsTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating testreports table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the testreports table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the testreports table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
