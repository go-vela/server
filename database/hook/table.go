// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres hooks table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
hooks (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	build_id     INTEGER,
	number       INTEGER,
	source_id    VARCHAR(250),
	created      INTEGER,
	host         VARCHAR(250),
	event        VARCHAR(250),
	event_action VARCHAR(250),
	branch       VARCHAR(500),
	error        VARCHAR(500),
	status       VARCHAR(250),
	link         VARCHAR(1000),
	webhook_id   INTEGER,
	UNIQUE(repo_id, number)
);
`

	// CreateSqliteTable represents a query to create the Sqlite hooks table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
hooks (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id      INTEGER,
	build_id     INTEGER,
	number       INTEGER,
	source_id    TEXT,
	created      INTEGER,
	host         TEXT,
	event        TEXT,
	event_action TEXT,
	branch       TEXT,
	error        TEXT,
	status       TEXT,
	link         TEXT,
	webhook_id   INTEGER,
	UNIQUE(repo_id, number)
);
`
)

// CreateHookTable creates the hooks table in the database.
func (e *engine) CreateHookTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating hooks table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the hooks table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the hooks table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
