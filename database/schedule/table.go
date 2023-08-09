// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres schedules table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
schedules (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	active       BOOLEAN,
	name         VARCHAR(100),
	entry        VARCHAR(100),
	created_at   INTEGER,
	created_by   VARCHAR(250),
	updated_at   INTEGER,
	updated_by   VARCHAR(250),
	scheduled_at INTEGER,
	UNIQUE(repo_id, name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite schedules table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
schedules (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id      INTEGER,
	active       BOOLEAN,
	name         TEXT,
	entry        TEXT,
	created_at   INTEGER,
	created_by   TEXT,
	updated_at   INTEGER,
	updated_by   TEXT,
	scheduled_at INTEGER,
	UNIQUE(repo_id, name)
);
`
)

// CreateScheduleTable creates the schedules table in the database.
func (e *engine) CreateScheduleTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating schedules table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the schedules table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the schedules table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
