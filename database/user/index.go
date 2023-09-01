// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
func (e *engine) CreateUserIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for users table in the database")

	// create the refresh_token column index for the users table
	return e.client.Exec(CreateUserRefreshIndex).Error
}
