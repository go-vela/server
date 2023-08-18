// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres pipelines table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
pipelines (
	id               SERIAL PRIMARY KEY,
	repo_id          INTEGER,
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
	data             BLOB,
	UNIQUE(repo_id, 'commit')
);
`
)

// CreatePipelineTable creates the pipelines table in the database.
func (e *engine) CreatePipelineTable(ctx context.Context, driver string) error {
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
