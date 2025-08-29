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

	// CreateCreatedAtIndex represents a query to create an
	// index on the logs table for the created_at column.
	CreateCreatedAtIndex = `
CREATE INDEX
IF NOT EXISTS
logs_created_at
ON logs (created_at);
`
)

// CreateLogIndexes creates the indexes for the logs table in the database.
func (e *Engine) CreateLogIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for logs table")

	indices := []string{
		CreateBuildIDIndex,
		CreateCreatedAtIndex,
	}

	for _, index := range indices {
		if err := e.client.WithContext(ctx).Exec(index).Error; err != nil {
			return err
		}
	}

	return nil
}
