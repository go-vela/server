// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	api "github.com/go-vela/server/api/types"
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

	// TestReport Management Functions

	// Count returns the count of all reports.
	Count(context.Context) (int64, error)

	// CountByBuild returns the count of reports by build ID.
	CountByBuild(context.Context, *api.Build, map[string]interface{}) (int64, error)

	// CountByRepo returns the count of reports by repo ID.
	CountByRepo(context.Context, *api.Repo, map[string]interface{}, int64, int64) (int64, error)

	// Create creates a new report.
	Create(context.Context, *api.TestReport) (*api.TestReport, error)

	// DeleteByID removes a report by ID.
	DeleteByID(context.Context, *api.TestReport) error

	// Get returns a report by ID.
	Get(context.Context, int64) (*api.TestReport, error)

	// List returns a list of all reports.
	List(context.Context, int, int) ([]*api.TestReport, int64, error)

	// ListByBuild returns a list of reports by build ID.
	ListByBuild(context.Context, *api.Build, int, int) ([]*api.TestReport, int64, error)

	// ListByRepo returns a list of reports by repo ID.
	ListByRepo(context.Context, *api.Repo, int, int) ([]*api.TestReport, int64, error)

	// Update updates a report.
	Update(context.Context, *api.TestReport) (*api.TestReport, error)
}
