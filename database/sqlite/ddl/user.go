// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

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

	// CreateUserRefreshIndex represents a query to create an
	// index on the users table for the refresh_token column.
	CreateUserRefreshIndex = `
CREATE INDEX
IF NOT EXISTS
users_refresh
ON users (refresh_token);
`
)
