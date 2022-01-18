// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateHookTable represents a query to
	// create the hooks table for Vela.
	CreateHookTable = `
CREATE TABLE
IF NOT EXISTS
hooks (
	id        SERIAL PRIMARY KEY,
	repo_id   INTEGER,
	build_id  INTEGER,
	number    INTEGER,
	source_id VARCHAR(250),
	created   INTEGER,
	host      VARCHAR(250),
	event     VARCHAR(250),
	branch    VARCHAR(500),
	error     VARCHAR(500),
	status    VARCHAR(250),
	link      VARCHAR(1000),
	UNIQUE(repo_id, number)
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
