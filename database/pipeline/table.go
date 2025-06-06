// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres pipelines table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
pipelines (
	id               BIGSERIAL PRIMARY KEY,
	repo_id          BIGINT,
	commit           VARCHAR(500),
	flavor           VARCHAR(100),
	platform         VARCHAR(100),
	ref              VARCHAR(500),
	type             VARCHAR(100),
	version          VARCHAR(50),
	external_secrets BOOLEAN,
	internal_secrets BOOLEAN,
	services         BOOLEAN,
	stages           BOOLEAN,
	steps            BOOLEAN,
	templates        BOOLEAN,
	warnings         VARCHAR(5000),
	data             BYTEA,
	UNIQUE(repo_id, commit)
);
`

	// CreateSqliteTable represents a query to create the Sqlite pipelines table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
pipelines (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id          INTEGER,
	'commit'         TEXT,
	flavor           TEXT,
	platform         TEXT,
	ref              TEXT,
	type             TEXT,
	version          TEXT,
	external_secrets BOOLEAN,
	internal_secrets BOOLEAN,
	services         BOOLEAN,
	stages           BOOLEAN,
	steps            BOOLEAN,
	templates        BOOLEAN,
	warnings         TEXT,
	data             BLOB,
	UNIQUE(repo_id, 'commit')
);
`
)

// CreatePipelineTable creates the pipelines table in the database.
func (e *Engine) CreatePipelineTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating pipelines table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the pipelines table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the pipelines table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
