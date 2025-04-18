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

	// CreatePostgresAllowlistTable represents a query to create the Postgres secrets_repo_allowlist table.
	CreatePostgresAllowlistTable = `
CREATE TABLE
IF NOT EXISTS
secret_repo_allowlist (
	id                 BIGSERIAL PRIMARY KEY,
	secret_id          BIGINT,
	repo               VARCHAR(500),
	UNIQUE(secret_id, repo)
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

	// CreateSqliteAllowlistTable represents a query to create the Sqlite secrets_allowlist table.
	CreateSqliteAllowlistTable = `
CREATE TABLE
IF NOT EXISTS
secret_repo_allowlist (
	id                 BIGSERIAL PRIMARY KEY,
	secret_id          INTEGER,
	repo               TEXT,
	UNIQUE(secret_id, repo)
);
`
)

// CreateSecretTables creates the secrets and secret_repo_allowlist tables in the database.
func (e *Engine) CreateSecretTables(ctx context.Context, driver string) error {
	e.logger.Tracef("creating secrets table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the secrets allowlist table for Postgres
		err := e.client.
			WithContext(ctx).
			Exec(CreatePostgresAllowlistTable).Error
		if err != nil {
			return err
		}

		// create the secrets table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the secrets allowlist table for Sqlite
		err := e.client.
			WithContext(ctx).
			Exec(CreateSqliteAllowlistTable).Error
		if err != nil {
			return err
		}

		// create the secrets table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
