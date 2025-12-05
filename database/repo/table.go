// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres repos table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id               BIGSERIAL PRIMARY KEY,
	user_id          BIGINT,
	hash             VARCHAR(500),
	org              VARCHAR(250),
	name             VARCHAR(250),
	full_name        VARCHAR(500),
	link             VARCHAR(1000),
	clone            VARCHAR(1000),
	branch           VARCHAR(250),
	topics           VARCHAR(1020),
	build_limit      INTEGER,
	timeout          INTEGER,
	counter          BIGINT,
	hook_counter     BIGINT,
	visibility       TEXT,
	private          BOOLEAN,
	trusted          BOOLEAN,
	active           BOOLEAN,
	allow_events     BIGINT,
	pipeline_type    TEXT,
	previous_name    VARCHAR(100),
	approve_build    VARCHAR(20),
	approval_timeout INTEGER,
	install_id       BIGINT,
	custom_props     JSON DEFAULT NULL,
	UNIQUE(full_name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite repos table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id               INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id          INTEGER,
	hash             TEXT,
	org              TEXT,
	name             TEXT,
	full_name        TEXT,
	link             TEXT,
	clone            TEXT,
	branch           TEXT,
	topics           TEXT,
	build_limit      INTEGER,
	timeout          INTEGER,
	counter          INTEGER,
	hook_counter     INTEGER,
	visibility       TEXT,
	private          BOOLEAN,
	trusted          BOOLEAN,
	active           BOOLEAN,
	allow_events     INTEGER,
	pipeline_type    TEXT,
	previous_name    TEXT,
	approve_build    TEXT,
	approval_timeout INTEGER,
	install_id       INTEGER,
	custom_props     TEXT,
	UNIQUE(full_name)
);
`
)

// CreateRepoTable creates the repos table in the database.
func (e *Engine) CreateRepoTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating repos table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the repos table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the repos table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
