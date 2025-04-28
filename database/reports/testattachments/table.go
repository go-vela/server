// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres testattachments table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
testattachments (
	id             		BIGSERIAL PRIMARY KEY,
	test_report_id		BIGINT,
	created        		BIGINT,
	file_name	  		VARCHAR(1000),
	object_path	  		VARCHAR(1000),
	file_size	  		INTEGER,
	file_type	  		TEXT,
	presigned_url		VARCHAR(2000),
	CONSTRAINT fk_testreport
	FOREIGN KEY (test_report_id)
	REFERENCES testreports(id)
	ON DELETE CASCADE
);
`

	// CreateSqliteTable represents a query to create the Sqlite testattachments table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
testattachments (
	id             	INTEGER PRIMARY KEY AUTOINCREMENT,
	test_report_id	INTEGER,
	created        	INTEGER,
	file_name	   	TEXT,
	object_path	   	TEXT,
	file_size	   	INTEGER,
    file_type 		TEXT,
	presigned_url	VARCHAR(2000),
    FOREIGN KEY (test_report_id) 
    REFERENCES testreports(id)
    ON DELETE CASCADE
);
`
)

// CreateTestAttachmentsTable creates the testattachments table in the database.
func (e *Engine) CreateTestAttachmentsTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating testattachments table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the testattachments table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the testattachments table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
