// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres users table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
users (
	id             SERIAL PRIMARY KEY,
	name           VARCHAR(250),
	refresh_token  VARCHAR(1000),
	token          VARCHAR(1000),
	hash           VARCHAR(500),
	favorites      VARCHAR(5000),
	active         BOOLEAN,
	admin          BOOLEAN,
	UNIQUE(name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite users table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
users (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	name           TEXT,
	refresh_token  TEXT,
	token          TEXT,
	hash           TEXT,
	favorites      TEXT,
	active         BOOLEAN,
	admin          BOOLEAN,
	UNIQUE(name)
);
`
)

// CreateUserTable creates the users table in the database.
func (e *engine) CreateUserTable(driver string) error {
	e.logger.Tracef("creating users table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the users table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the users table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
