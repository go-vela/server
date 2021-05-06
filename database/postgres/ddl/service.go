// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateServiceTable represents a query to
	// create the services table for Vela.
	CreateServiceTable = `
CREATE TABLE
IF NOT EXISTS
services (
	id            SERIAL PRIMARY KEY,
	repo_id       INTEGER,
	build_id      INTEGER,
	number        INTEGER,
	name          VARCHAR(250),
	image         VARCHAR(500),
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
