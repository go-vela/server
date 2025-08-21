// SPDX-License-Identifier: Apache-2.0

package testreport

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
	id             BIGSERIAL PRIMARY KEY,
	build_id       BIGINT,
	created_at     BIGINT
);
`

	// CreateSqliteTable represents a query to create the Sqlite testreports table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
testreports (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id       BIGINT,
	created_at     BIGINT
);
`
)

// CreateTestReportTable creates the testreports table in the database.
func (e *Engine) CreateTestReportTable(ctx context.Context, driver string) error {
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
