// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

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

	// CreateLogStepIDIndex represents a query to create an
	// index on the logs table for the step_id column.
	CreateLogStepIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_step_id
ON logs (step_id);
`

	// CreateLogServiceIDIndex represents a query to create an
	// index on the logs table for the service_id column.
	CreateLogServiceIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_service_id
ON logs (service_id);
`
)

// createLogService is a helper function to return
// a service for interacting with the logs table.
func createLogService() *Service {
	return &Service{
		Create:  CreateLogTable,
		Indexes: []string{CreateLogBuildIDIndex, CreateLogStepIDIndex, CreateLogServiceIDIndex},
	}
}
