// SPDX-License-Identifier: Apache-2.0

package testreport

import "context"

const (
	// CreateCreatedIndex represents a query to create an
	// index on the testreports table for the created_at column.
	CreateCreatedIndex = `
CREATE INDEX
IF NOT EXISTS
testreports_created_at
ON testreports (created_at);
`

	// CreateBuildIDIndex represents a query to create an
	// index on the testreports table for the build_id column.
	CreateBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
testreports_build_id
ON testreports (build_id);
`
)

// CreateTestReportIndexes creates the indexes for the testreports table in the database.
func (e *Engine) CreateTestReportIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for testreports table")

	// create the created_at column index for the testreports table
	err := e.client.
		WithContext(ctx).
		Exec(CreateCreatedIndex).Error
	if err != nil {
		return err
	}

	// create the build_id column index for the testreports table
	return e.client.
		WithContext(ctx).
		Exec(CreateBuildIDIndex).Error
}
