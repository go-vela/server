// SPDX-License-Identifier: Apache-2.0

package log

import "context"

const (
	// CreateBuildIDIndex represents a query to create an
	// index on the logs table for the build_id column.
	CreateBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_build_id
ON logs (build_id);
`
)

// CreateLogIndexes creates the indexes for the logs table in the database.
func (e *engine) CreateLogIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for logs table in the database")

	// create the build_id column index for the logs table
	return e.client.Exec(CreateBuildIDIndex).Error
}
