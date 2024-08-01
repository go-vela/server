// SPDX-License-Identifier: Apache-2.0

package pipeline

import "context"

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
func (e *engine) CreatePipelineIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for pipelines table in the database")

	// create the repo_id column index for the pipelines table
	return e.client.Exec(CreateRepoIDIndex).Error
}
