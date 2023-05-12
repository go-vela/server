// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import "github.com/go-vela/types/constants"

const (
	// CreatePostgresTable represents a query to create the Postgres build_itineraries table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
build_itineraries (
	id               SERIAL PRIMARY KEY,
	build_id         INTEGER,
	data             BYTEA,
	UNIQUE(build_id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite build_itineraries table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
build_itineraries (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id         INTEGER,
	data             BLOB,
	UNIQUE(build_id)
);
`
)

// CreateBuildItineraryTable creates the build itineraries table in the database.
func (e *engine) CreateBuildItineraryTable(driver string) error {
	e.logger.Tracef("creating build_itineraries table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the build_itineraries table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the build_itineraries table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
