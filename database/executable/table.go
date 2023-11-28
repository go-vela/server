// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres build_executables table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
build_executables (
	id               SERIAL PRIMARY KEY,
	build_id         INTEGER,
	data             BYTEA,
	UNIQUE(build_id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite build_executables table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
build_executables (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id         INTEGER,
	data             BLOB,
	UNIQUE(build_id)
);
`
)

// CreateBuildExecutableTable creates the build executables table in the database.
func (e *engine) CreateBuildExecutableTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating build_executables table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the build_executables table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the build_executables table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
