// SPDX-License-Identifier: Apache-2.0

package installation

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres users table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
installations (
	install_id BIGSERIAL PRIMARY KEY,
	target     VARCHAR(250),
	UNIQUE (target)
);
`

	// CreateSqliteTable represents a query to create the Sqlite installations table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
installations (
	install_id INTEGER PRIMARY KEY AUTOINCREMENT,
	target     TEXT,
	UNIQUE (target)
);
`
)

// CreateInstallationTable creates the installations table in the database.
func (e *Engine) CreateInstallationTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating installations table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the installations table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the installations table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
