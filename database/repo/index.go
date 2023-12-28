// SPDX-License-Identifier: Apache-2.0

package repo

import "context"

const (
	// CreateOrgNameIndex represents a query to create an
	// index on the repos table for the org and name columns.
	CreateOrgNameIndex = `
CREATE INDEX
IF NOT EXISTS
repos_org_name
ON repos (org, name);
`
)

// CreateRepoIndexes creates the indexes for the repos table in the database.
func (e *engine) CreateRepoIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for repos table in the database")

	// create the org and name columns index for the repos table
	return e.client.Exec(CreateOrgNameIndex).Error
}
