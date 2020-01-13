// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
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
	id        SERIAL PRIMARY KEY,
	name      VARCHAR(250),
	token     VARCHAR(500),
	hash      VARCHAR(500),
	favorites VARCHAR(1000),
	active    BOOLEAN,
	admin     BOOLEAN,
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
)

// createUserService is a helper function to return
// a service for interacting with the users table.
func createUserService() *Service {
	return &Service{
		Create:  CreateUserTable,
		Indexes: []string{CreateUserNameIndex},
	}
}
