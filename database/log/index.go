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

	// CreateServiceIDIndex represents a query to create an
	// index on the logs table for the service_id column.
	CreateServiceIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_service_id
ON logs (service_id);
`

	// CreateStepIDIndex represents a query to create an
	// index on the logs table for the step_id column.
	CreateStepIDIndex = `
CREATE INDEX
IF NOT EXISTS
logs_step_id
ON logs (step_id);
`
)

// CreateLogIndexes creates the indexes for the logs table in the database.
func (e *Engine) CreateLogIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for logs table")

	indices := []string{
		CreateBuildIDIndex,
		CreateCreatedAtIndex,
		CreateServiceIDIndex,
		CreateStepIDIndex,
	}

	for _, index := range indices {
		if err := e.client.WithContext(ctx).Exec(index).Error; err != nil {
			return err
		}
	}

	return nil
}
