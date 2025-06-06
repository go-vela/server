// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres users table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
users (
	id             BIGSERIAL PRIMARY KEY,
	name           VARCHAR(250),
	refresh_token  VARCHAR(500),
	token          VARCHAR(500),
	favorites      VARCHAR(5000),
	active         BOOLEAN,
	admin          BOOLEAN,
	dashboards     VARCHAR(5000),
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
	favorites      TEXT,
	active         BOOLEAN,
	admin          BOOLEAN,
	dashboards     TEXT,
	UNIQUE(name)
);
`
)

// CreateUserTable creates the users table in the database.
func (e *Engine) CreateUserTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating users table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the users table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the users table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
