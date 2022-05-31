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
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id      INTEGER,
	build_id     INTEGER,
	number       INTEGER,
	source_id    TEXT,
	created      INTEGER,
	host         TEXT,
	event        TEXT,
	event_action TEXT,
	branch       TEXT,
	error        TEXT,
	status       TEXT,
	link         TEXT,
	webhook_id   INTEGER,
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
