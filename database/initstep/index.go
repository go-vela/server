// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

const (
	// CreateBuildIDIndex represents a query to create an
	// index on the inits table for the build_id column.
	CreateBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
initsteps_build_id
ON initsteps (build_id);
`
)

// CreateInitStepIndexes creates the indexes for the inits table in the database.
func (e *engine) CreateInitStepIndexes() error {
	e.logger.Tracef("creating indexes for initsteps table in the database")

	// create the hostname and address columns index for the inits table
	return e.client.Exec(CreateBuildIDIndex).Error
}
