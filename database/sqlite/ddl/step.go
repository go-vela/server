// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateStepTable represents a query to
	// create the steps table for Vela.
	CreateStepTable = `
CREATE TABLE
IF NOT EXISTS
steps (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          TEXT,
	image         TEXT,
	stage         TEXT,
	status        TEXT,
	error         TEXT,
	exit_code     INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	host          TEXT,
	runtime       TEXT,
	distribution  TEXT,
	UNIQUE(build_id, number)
);
`
)
