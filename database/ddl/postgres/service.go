// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

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
	UNIQUE(build_id, number)
);
`

	// CreateServiceBuildIDNumberIndex represents a query to create an
	// index on the services table for the build_id and number columns.
	CreateServiceBuildIDNumberIndex = `
CREATE INDEX
IF NOT EXISTS
services_build_id_number
ON services (build_id, number);
`
)

// createServiceService is a helper function to return
// a service for interacting with the services table.
func createServiceService() *Service {
	return &Service{
		Create:  CreateServiceTable,
		Indexes: []string{CreateServiceBuildIDNumberIndex},
	}
}
