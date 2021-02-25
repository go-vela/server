// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// CreateHookTable represents a query to
	// create the hooks table for Vela.
	CreateHookTable = `
CREATE TABLE
IF NOT EXISTS
hooks (
	id        INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id   INTEGER,
	build_id  INTEGER,
	number    INTEGER,
	source_id TEXT,
	created   INTEGER,
	host      TEXT,
	event     TEXT,
	branch    TEXT,
	error     TEXT,
	status    TEXT,
	link      TEXT,
	UNIQUE(repo_id, build_id)
);
`

	// CreateHookRepoIDIndex represents a query to create an
	// index on the hooks table for the repo_id column.
	CreateHookRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
hooks_repo_id
ON hooks (repo_id);
`
)

// createHookService is a helper function to return
// a service for interacting with the hooks table.
func createHookService() *Service {
	return &Service{
		Create:  CreateHookTable,
		Indexes: []string{CreateHookRepoIDIndex},
	}
}
