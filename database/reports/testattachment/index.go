// SPDX-License-Identifier: Apache-2.0

package testattachment

import "context"

const (
	// CreateTestReportIDIndex represents a query to create an
	// index on the testattachments table for the test_report_id column.
	CreateTestReportIDIndex = `
CREATE INDEX
IF NOT EXISTS
testattachments_test_report_id
ON testattachments (test_report_id);
`
)

// CreateTestAttachmentIndexes creates the indexes for the testattachments table in the database.
func (e *Engine) CreateTestAttachmentIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for testattachments table")

	// create the test_report_id column index for the testattachments table
	return e.client.
		WithContext(ctx).
		Exec(CreateTestReportIDIndex).Error
}
