// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

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
func (e *engine) CreateHookIndexes() error {
	e.logger.Tracef("creating indexes for hooks table in the database")

	// create the hostname and address columns index for the hooks table
	return e.client.Exec(CreateRepoIDIndex).Error
}
