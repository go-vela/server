// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

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

	// CreateStepBuildIDNumberIndex represents a query to create an
	// index on the steps table for the build_id and number columns.
	CreateStepBuildIDNumberIndex = `
CREATE INDEX
IF NOT EXISTS
steps_build_id_number
ON steps (build_id, number);
`
)

// createStepService is a helper function to return
// a service for interacting with the steps table.
func createStepService() *Service {
	return &Service{
		Create:  CreateStepTable,
		Indexes: []string{CreateStepBuildIDNumberIndex},
	}
}
