// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

const (
	// CreateRepoIDIndex represents a query to create an
	// index on the schedules table for the repo_id column.
	CreateRepoIDIndex = `
CREATE INDEX
IF NOT EXISTS
schedules_repo_id
ON schedules (repo_id);
`
)

// CreateScheduleIndexes creates the indexes for the schedules table in the database.
func (e *engine) CreateScheduleIndexes() error {
	e.logger.Tracef("creating indexes for schedules table in the database")

	// create the repo_id column index for the schedules table
	return e.client.Exec(CreateRepoIDIndex).Error
}
