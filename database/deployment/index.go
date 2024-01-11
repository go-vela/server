// SPDX-License-Identifier: Apache-2.0

package deployment

import "context"

const (
	// CreateRepoIDIndex represents a query to create an
	// index on the deployments table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
deployments_repo_id
ON deployments (repo_id);
`
)

// CreateDeploymetsIndexes creates the indexes for the deployments table in the database.
func (e *engine) CreateDeploymentIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for deployments table in the database")

	// create the repo_id column index for the deployments table
	return e.client.Exec(CreateRepoIDIndex).Error
}
