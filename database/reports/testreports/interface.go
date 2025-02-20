// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
)

// TestReportsInterface represents the Vela interface for testreports
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type TestReportsInterface interface {
	// TestReports Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateTestReportsIndexes defines a function that creates the indexes for the testreports table.
	CreateTestReportsIndexes(context.Context) error
	// CreateTestReportsTable defines a function that creates the testreports table.
	CreateTestReportsTable(context.Context, string) error
}
