// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

const (
	// CreateBuildIDIndex represents a query to create an
	// index on the inits table for the build_id column.
	CreateBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
inits_build_id
ON inits (build_id);
`
)

// CreateHookIndexes creates the indexes for the inits table in the database.
func (e *engine) CreateInitsIndexes() error {
	e.logger.Tracef("creating indexes for inits table in the database")

	// create the hostname and address columns index for the inits table
	return e.client.Exec(CreateBuildIDIndex).Error
}
