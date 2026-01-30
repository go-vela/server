// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres artifacts table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
artifacts (
	id             		BIGSERIAL PRIMARY KEY,
	build_id       		BIGINT,
	created_at        	BIGINT,
	file_name	  		VARCHAR(1000),
	object_path	  		VARCHAR(1000),
	file_size	  		INTEGER,
	file_type	  		TEXT,
	presigned_url		VARCHAR(2000)
);
`

	// 	// CreatePostgresTable represents a query to create the Postgres artifacts table.
	// 	CreatePostgresTable = `
	// CREATE TABLE
	// IF NOT EXISTS
	// artifacts (
	// 	id             		BIGSERIAL PRIMARY KEY,
	// 	build_id		BIGINT,
	// 	created_at        	BIGINT,
	// 	file_name	  		VARCHAR(1000),
	// 	object_path	  		VARCHAR(1000),
	// 	file_size	  		INTEGER,
	// 	file_type	  		TEXT,
	// 	presigned_url		VARCHAR(2000),
	// );
	// `.

	// CreateSqliteTable represents a query to create the Sqlite artifacts table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
artifacts (
	id             	INTEGER PRIMARY KEY AUTOINCREMENT,
	build_id        BIGINT,
	created_at      BIGINT,
	file_name	   	TEXT,
	object_path	   	TEXT,
	file_size	   	INTEGER,
    file_type 		TEXT,
	presigned_url	VARCHAR(2000)
);
`

//	// CreateSqliteTable represents a query to create the Sqlite artifacts table.
//	CreateSqliteTable = `
//
// CREATE TABLE
// IF NOT EXISTS
// artifacts (
//
//		id             	INTEGER PRIMARY KEY AUTOINCREMENT,
//		build_id	INTEGER,
//		created_at      INTEGER,
//		file_name	   	TEXT,
//		object_path	   	TEXT,
//		file_size	   	INTEGER,
//	    file_type 		TEXT,
//		presigned_url	VARCHAR(2000),
//
// );
// `
)

// CreateArtifactTable creates the artifacts table in the database.
func (e *Engine) CreateArtifactTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating artifacts table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the artifacts table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the artifacts table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
