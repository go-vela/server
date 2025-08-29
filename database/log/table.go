// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/server/constants"
)

// Consider partitioning as the logs table will be running hot with the current
// access pattern - primarily driven by the workers. The table could suffer
// from excessive bloat. Partitioning will help spread the load.
//
// CREATE TABLE logs (
//     id            BIGSERIAL,
//     build_id      BIGINT,
//     repo_id       BIGINT,
//     service_id    BIGINT,
//     step_id       BIGINT,
//     data          BYTEA,
//     created_at    BIGINT NOT NULL,
//     PRIMARY KEY (id)
// ) PARTITION BY HASH (id);
//
// We do not need unique indices as you would have to create a unique index
// combined with 'id' (partition key) which will not provide any advantages.
// However, you would want to consider regular indices on 'build_id', 'service_id',
// 'step_id', and 'created_at' to help with query performance.
//
// Example paritioning query:
// DO $$
// BEGIN
//     FOR i IN 0..19 LOOP
//         EXECUTE format('CREATE TABLE logs_partition_%s PARTITION OF logs
//                        FOR VALUES WITH (modulus 20, remainder %s)', i, i);
//     END LOOP;
// END $$;
//
// Then, create indices:
// CREATE INDEX IF NOT EXISTS logs_build_id ON logs (build_id);
// CREATE INDEX IF NOT EXISTS logs_service_id ON logs (service_id);
// CREATE INDEX IF NOT EXISTS logs_step_id ON logs (step_id);
// CREATE INDEX IF NOT EXISTS logs_created_at ON logs (created_at);
//
// Note: SQLite does not support partitioning, so this is not an option.

const (
	// CreatePostgresTable represents a query to create the Postgres logs table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            BIGSERIAL PRIMARY KEY,
	build_id      BIGINT,
	repo_id       BIGINT,
	service_id    BIGINT,
	step_id       BIGINT,
	data          BYTEA,
	created_at    BIGINT,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite logs table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id      INTEGER,
	repo_id       INTEGER,
	service_id    INTEGER,
	step_id       INTEGER,
	data          BLOB,
	created_at    INTEGER,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`
)

// CreateLogTable creates the logs table in the database.
func (e *Engine) CreateLogTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating logs table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the logs table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the logs table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
