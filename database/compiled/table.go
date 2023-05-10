// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import "github.com/go-vela/types/constants"

const (
	// CreatePostgresTable represents a query to create the Postgres pipelines table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
compiled (
	id               SERIAL PRIMARY KEY,
	build_id         INTEGER,
	data             BYTEA,
	UNIQUE(build_id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite pipelines table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
compiled (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id         INTEGER,
	data             BLOB,
	UNIQUE(build_id)
);
`
)

// CreatePipelineTable creates the pipelines table in the database.
func (e *engine) CreateCompiledTable(driver string) error {
	e.logger.Tracef("creating pipelines table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the pipelines table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the pipelines table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
