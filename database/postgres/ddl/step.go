// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	id            SERIAL PRIMARY KEY,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          VARCHAR(250),
	image         VARCHAR(500),
	stage         VARCHAR(250),
	status        VARCHAR(250),
	error         VARCHAR(500),
	exit_code     INTEGER,
	created       INTEGER,
	started       INTEGER,
	finished      INTEGER,
	host          VARCHAR(250),
	runtime       VARCHAR(250),
	distribution  VARCHAR(250),
	UNIQUE(build_id, number)
);
`
)
