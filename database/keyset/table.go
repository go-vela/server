// SPDX-License-Identifier: Apache-2.0

package keyset

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres keysets table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
keysets (
	id     UUID PRIMARY KEY,
	active BOOLEAN,
	key    JSON DEFAULT NULL
);
`

	// CreateSqliteTable represents a query to create the Sqlite keysets table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
keysets (
	id     TEXT PRIMARY KEY,
	active BOOLEAN,
	key    TEXT
);
`
)

// CreateKeySetTable creates the build executables table in the database.
func (e *engine) CreateKeySetTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating keysets table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the keysets table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the keysets table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
