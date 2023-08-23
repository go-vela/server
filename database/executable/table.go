// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executable

import "github.com/go-vela/types/constants"

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
func (e *engine) CreateBuildExecutableTable(driver string) error {
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
