// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres users table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
favorites (
  user_id    BIGINT      NOT NULL,
  repo_id    BIGINT      NOT NULL,
  position   INTEGER,
  PRIMARY KEY (user_id, repo_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (repo_id) REFERENCES repos(id)
);
`

	// CreateSqliteTable represents a query to create the Sqlite users table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
favorites (
  user_id    INTEGER     NOT NULL,
  repo_id    INTEGER     NOT NULL,
  position   INTEGER,
  PRIMARY KEY (user_id, repo_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (repo_id) REFERENCES repositories(id)
);
`
)

// CreateFavoritesTable creates the favorites table in the database.
func (e *Engine) CreateFavoritesTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating favorites table")

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the favorites table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the favorites table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
