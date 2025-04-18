// SPDX-License-Identifier: Apache-2.0

package testattachments

import "context"

const (
	// CreateTestReportIDIndex represents a query to create an
	// index on the testattachments table for the testreport_id column.
	CreateTestReportIDIndex = `
CREATE INDEX
IF NOT EXISTS
testattachments_testreport_id
ON testattachments (testreport_id);
`
)


// CreateTestAttachmentsIndexes creates the indexes for the testattachments table in the database.
func (e *engine) CreateTestAttachmentsIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for testattachments table")

	// create the testreport_id column index for the testattachments table
	return e.client.
		WithContext(ctx).
		Exec(CreateTestReportIDIndex).Error

}
