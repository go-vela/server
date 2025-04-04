// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres secrets table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
secrets (
	id                 BIGSERIAL PRIMARY KEY,
	type               VARCHAR(100),
	org                VARCHAR(250),
	repo               VARCHAR(250),
	team               VARCHAR(250),
	name               VARCHAR(250),
	value              BYTEA,
	images             VARCHAR(1000),
	allow_events       BIGINT,
	allow_command      BOOLEAN,
	allow_substitution BOOLEAN,
	created_at         BIGINT,
	created_by         VARCHAR(250),
	updated_at         BIGINT,
	updated_by         VARCHAR(250),
	UNIQUE(type, org, repo, name),
	UNIQUE(type, org, team, name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite secrets table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
secrets (
	id                 INTEGER PRIMARY KEY AUTOINCREMENT,
	type               TEXT,
	org                TEXT,
	repo               TEXT,
	team               TEXT,
	name               TEXT,
	value              TEXT,
	images             TEXT,
	allow_events       INTEGER,
	allow_command      BOOLEAN,
	allow_substitution BOOLEAN,
	created_at         INTEGER,
	created_by	       TEXT,
	updated_at         INTEGER,
	updated_by         TEXT,
	UNIQUE(type, org, repo, name),
	UNIQUE(type, org, team, name)
);
`
)

// CreateSecretTable creates the secrets table in the database.
func (e *engine) CreateSecretTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating secrets table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the secrets table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the secrets table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
