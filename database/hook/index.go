// SPDX-License-Identifier: Apache-2.0

package hook

import "context"

const (
	// CreateRepoIDIndex represents a query to create an
	// index on the hooks table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
hooks_repo_id
ON hooks (repo_id);
`
)

// CreateHookIndexes creates the indexes for the hooks table in the database.
func (e *engine) CreateHookIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for hooks table in the database")

	// create the repo_id column index for the hooks table
	return e.client.Exec(CreateRepoIDIndex).Error
}
