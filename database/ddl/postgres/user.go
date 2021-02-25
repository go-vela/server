// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// CreateUserTable represents a query to
	// create the users table for Vela.
	CreateUserTable = `
CREATE TABLE
IF NOT EXISTS
users (
	id             SERIAL PRIMARY KEY,
	name           VARCHAR(250),
	refresh_token  VARCHAR(500),
	token          VARCHAR(500),
	hash           VARCHAR(500),
	favorites      VARCHAR(5000),
	active         BOOLEAN,
	admin          BOOLEAN,
	UNIQUE(name)
);
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
		Indexes: []string{CreateRefreshIndex},
	}
}
