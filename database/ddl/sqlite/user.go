// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// CreateUserTable represents a query to
	// create the users table for Vela.
	CreateUserTable = `
CREATE TABLE
IF NOT EXISTS
users (
	id             INTEGER PRIMARY KEY AUTOINCREMENT,
	name           TEXT,
	refresh_token  TEXT,
	token          TEXT,
	hash           TEXT,
	favorites      TEXT,
	active         BOOLEAN,
	admin          BOOLEAN,
	UNIQUE(name)
);
`

	// CreateUserNameIndex represents a query to create an
	// index on the users table for the name column.
	CreateUserNameIndex = `
CREATE INDEX
IF NOT EXISTS
users_name
ON users (name);
`

	// CreateRefreshIndex represents a query to create an
	// index on the users table for the refresh_token column.
	CreateRefreshIndex = `
CREATE INDEX
IF NOT EXISTS
users_refresh
ON users (refresh_token);
`
)

// createUserService is a helper function to return
// a service for interacting with the users table.
func createUserService() *Service {
	return &Service{
		Create:  CreateUserTable,
		Indexes: []string{CreateUserNameIndex, CreateRefreshIndex},
	}
}
