// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres settings table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
settings (
	id            SERIAL PRIMARY KEY,
	foo_num       INTEGER,
	foo_str       VARCHAR(250)
);
`

	// CreateSqliteTable represents a query to create the Sqlite settings table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
settings (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	foo_num      INTEGER,
	foo_str          TEXT
);
`
)

// CreateSettingsTable creates the settings table in the database.
func (e *engine) CreateSettingsTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating settings table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the steps table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the steps table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
