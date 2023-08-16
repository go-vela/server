// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import "context"

const (
	// CreateCreatedIndex represents a query to create an
	// index on the builds table for the created column.
	CreateCreatedIndex = `
CREATE INDEX
IF NOT EXISTS
builds_created
ON builds (created);
`

	// CreateRepoIDIndex represents a query to create an
	// index on the builds table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
builds_repo_id
ON builds (repo_id);
`

	// CreateSourceIndex represents a query to create an
	// index on the builds table for the source column.
	CreateSourceIndex = `
CREATE INDEX
IF NOT EXISTS
builds_source
ON builds (source);
`

	// CreateStatusIndex represents a query to create an
	// index on the builds table for the status column.
	CreateStatusIndex = `
CREATE INDEX
IF NOT EXISTS
builds_status
ON builds (status);
`
)

// CreateBuildIndexes creates the indexes for the builds table in the database.
func (e *engine) CreateBuildIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for builds table in the database")

	// create the created column index for the builds table
	err := e.client.Exec(CreateCreatedIndex).Error
	if err != nil {
		return err
	}

	// create the repo_id column index for the builds table
	err = e.client.Exec(CreateRepoIDIndex).Error
	if err != nil {
		return err
	}

	// create the source column index for the builds table
	err = e.client.Exec(CreateSourceIndex).Error
	if err != nil {
		return err
	}

	// create the status column index for the builds table
	return e.client.Exec(CreateStatusIndex).Error
}
