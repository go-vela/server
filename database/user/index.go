// SPDX-License-Identifier: Apache-2.0

package user

import "context"

const (
	// CreateUserRefreshIndex represents a query to create an
	// index on the users table for the refresh_token column.
	CreateUserRefreshIndex = `
CREATE INDEX
IF NOT EXISTS
users_refresh
ON users (refresh_token);
`
)

// CreateUserIndexes creates the indexes for the users table in the database.
func (e *Engine) CreateUserIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for users table")

	// create the refresh_token column index for the users table
	return e.client.
		WithContext(ctx).
		Exec(CreateUserRefreshIndex).Error
}
