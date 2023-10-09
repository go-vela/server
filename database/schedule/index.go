// SPDX-License-Identifier: Apache-2.0

package schedule

import "context"

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
func (e *engine) CreateScheduleIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for schedules table in the database")

	// create the repo_id column index for the schedules table
	return e.client.Exec(CreateRepoIDIndex).Error
}
