// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// CreateRepoTable represents a query to
	// create the repos table for Vela.
	// todo - look into what do for db changes
	CreateRepoTable = `
CREATE TABLE
IF NOT EXISTS
repos (
	id            SERIAL PRIMARY KEY,
	user_id       INTEGER,
	hash          VARCHAR(500),
	org           VARCHAR(250),
	name          VARCHAR(250),
	full_name     VARCHAR(500),
	link          VARCHAR(1000),
	clone         VARCHAR(1000),
	branch        VARCHAR(250),
	timeout       INTEGER,
	max_builds    INTEGER, 
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

	// CreateRepoFullNameIndex represents a query to create an
	// index on the repos table for the full_name column.
	CreateRepoFullNameIndex = `
CREATE INDEX
IF NOT EXISTS
repos_full_name
ON repos (full_name);
`
)

// createRepoService is a helper function to return
// a service for interacting with the repos table.
func createRepoService() *Service {
	return &Service{
		Create:  CreateRepoTable,
		Indexes: []string{CreateRepoOrgNameIndex, CreateRepoFullNameIndex},
	}
}
