// SPDX-License-Identifier: Apache-2.0

package artifact

import "context"

const (
	// CreateBuildIDIndex represents a query to create an
	// index on the artifacts table for the build_id column.
	CreateBuildIDIndex = `
CREATE INDEX
IF NOT EXISTS
artifacts_build_id
ON artifacts (build_id);
`
)

// CreateArtifactIndexes creates the indexes for the artifacts table in the database.
func (e *Engine) CreateArtifactIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for artifacts table")

	// create the build_id column index for the artifacts table
	return e.client.
		WithContext(ctx).
		Exec(CreateBuildIDIndex).Error
}
