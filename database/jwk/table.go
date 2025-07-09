// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres jwks table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
jwks (
	id     UUID PRIMARY KEY,
	active BOOLEAN,
	key    JSON DEFAULT NULL
);
`

	// CreateSqliteTable represents a query to create the Sqlite jwks table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
jwks (
	id     TEXT PRIMARY KEY,
	active BOOLEAN,
	key    TEXT
);
`
)

// CreateJWKTable creates the jwks table in the database.
func (e *Engine) CreateJWKTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating jwks table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the jwks table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the jwks table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
