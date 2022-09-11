// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

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
func (e *engine) CreateRepoIndexes() error {
	e.logger.Tracef("creating indexes for repos table in the database")

	// create the repo_id column index for the repos table
	return e.client.Exec(CreateOrgNameIndex).Error
}
