// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// CreateRepoTable represents a query to
	// create the repos table for Vela.
	CreateRepoTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id       INTEGER,
	hash          TEXT,
	org           TEXT,
	name          TEXT,
	full_name     TEXT,
	link          TEXT,
	clone         TEXT,
	branch        TEXT,
	timeout       INTEGER,
	visibility    TEXT,
	private       BOOLEAN,
	trusted       BOOLEAN,
	active        BOOLEAN,
	allow_pull    BOOLEAN,
	allow_push    BOOLEAN,
	allow_deploy  BOOLEAN,
	allow_tag     BOOLEAN,
	allow_comment BOOLEAN,
	UNIQUE(full_name)
);
`

	// CreateRepoOrgNameIndex represents a query to create an
	// index on the repos table for the org and name columns.
	CreateRepoOrgNameIndex = `
CREATE INDEX
IF NOT EXISTS
repos_org_name
ON repos (org, name);
`
)

// createRepoService is a helper function to return
// a service for interacting with the repos table.
func createRepoService() *Service {
	return &Service{
		Create:  CreateRepoTable,
		Indexes: []string{CreateRepoOrgNameIndex},
	}
}
