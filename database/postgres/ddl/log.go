// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package ddl

const (
	// CreateLogTable represents a query to
	// create the logs table for Vela.
	CreateLogTable = `
CREATE TABLE
IF NOT EXISTS
logs (
	id            SERIAL PRIMARY KEY,
	build_id      INTEGER,
	repo_id       INTEGER,
	service_id    INTEGER,
	step_id       INTEGER,
	data          BYTEA,
	UNIQUE(step_id),
	UNIQUE(service_id)
);
`

	// CreateLogBuildIDIndex represents a query to create an
	// index on the logs table for the build_id column.
	CreateLogBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_build_id
ON logs (build_id);
`
)
