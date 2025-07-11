// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres builds table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id             BIGSERIAL PRIMARY KEY,
	repo_id        BIGINT,
	pipeline_id    BIGINT,
	number         BIGINT,
	parent         BIGINT,
	event          VARCHAR(250),
	event_action   VARCHAR(250),
	status         VARCHAR(250),
	error          VARCHAR(1000),
	enqueued       BIGINT,
	created        BIGINT,
	started        BIGINT,
	finished       BIGINT,
	deploy         VARCHAR(500),
	deploy_number  BIGINT,
	deploy_payload VARCHAR(2000),
	clone          VARCHAR(1000),
	source         VARCHAR(1000),
	title          VARCHAR(1000),
	message        VARCHAR(2000),
	commit         VARCHAR(500),
	sender         VARCHAR(250),
	sender_scm_id  VARCHAR(250),
	fork           BOOLEAN,
	author         VARCHAR(250),
	email          VARCHAR(500),
	link           VARCHAR(1000),
	branch         VARCHAR(500),
	ref            VARCHAR(500),
	base_ref       VARCHAR(500),
	head_ref       VARCHAR(500),
	host           VARCHAR(250),
	route          VARCHAR(250),
	runtime        VARCHAR(250),
	distribution   VARCHAR(250),
	approved_at    BIGINT,
	approved_by    VARCHAR(250),
	timestamp      BIGINT,
	UNIQUE(repo_id, number)
);
`

	// CreateSqliteTable represents a query to create the Sqlite builds table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
builds (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id        INTEGER,
	pipeline_id    INTEGER,
	number         INTEGER,
	parent         INTEGER,
	event          TEXT,
	event_action   TEXT,
	status         TEXT,
	error          TEXT,
	enqueued       INTEGER,
	created        INTEGER,
	started        INTEGER,
	finished       INTEGER,
	deploy         TEXT,
	deploy_number  INTEGER,
	deploy_payload TEXT,
	clone          TEXT,
	source         TEXT,
	title          TEXT,
	message        TEXT,
	'commit'       TEXT,
	sender         TEXT,
	sender_scm_id  TEXT,
	fork           BOOLEAN,
	author         TEXT,
	email          TEXT,
	link           TEXT,
	branch         TEXT,
	ref            TEXT,
	base_ref       TEXT,
	head_ref       TEXT,
	host           TEXT,
	route          TEXT,
	runtime        TEXT,
	distribution   TEXT,
	approved_at    INTEGER,
	approved_by    TEXT,
	timestamp      INTEGER,
	UNIQUE(repo_id, number)
);
`
)

// CreateBuildTable creates the builds table in the database.
func (e *Engine) CreateBuildTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating builds table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the builds table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the builds table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
