// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres deployments table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
deployments (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	number       INTEGER,
	url          VARCHAR(500),
	commit       VARCHAR(500),
	ref          VARCHAR(500),
	task         VARCHAR(500),
	target       VARCHAR(500),
	description  VARCHAR(2500),
	payload      VARCHAR(2500),
	created_at   INTEGER,
	created_by   VARCHAR(250),
	builds       VARCHAR(500),
	UNIQUE(repo_id, number)
);
`

	// CreateSqliteTable represents a query to create the Sqlite deployments table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
deployments (
	id           SERIAL PRIMARY KEY,
	repo_id      INTEGER,
	number       INTEGER,	
	url     	 VARCHAR(1000),
	"commit"     VARCHAR(500),
	ref          VARCHAR(500),
	task         VARCHAR(500),
	target       VARCHAR(500),
	description  VARCHAR(2500),
	payload      VARCHAR(2500),
	created_at   INTEGER,
	created_by   VARCHAR(250),
	builds       VARCHAR(50),
	UNIQUE(repo_id, number)
);
`
)

// CreateDeploymentTable creates the deployments table in the database.
func (e *engine) CreateDeploymentTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating deployments table in the database")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the deployments table for Postgres
		return e.client.Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the deployments table for Sqlite
		return e.client.Exec(CreateSqliteTable).Error
	}
}
