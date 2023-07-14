// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

const (
	// CreateRepoIDIndex represents a query to create an
	// index on the pipelines table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
pipelines_repo_id
ON pipelines (repo_id);
`
)

// CreatePipelineIndexes creates the indexes for the pipelines table in the database.
func (e *engine) CreatePipelineIndexes() error {
	e.logger.Tracef("creating indexes for pipelines table in the database")

	// create the repo_id column index for the pipelines table
	return e.client.Exec(CreateRepoIDIndex).Error
}
