// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// TestReportInterface represents the Vela interface for testreports
// functions with the supported Database backends.
//

type TestReportInterface interface {
	// Artifacts Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateTestReportIndexes defines a function that creates the indexes for the testreports table.
	CreateTestReportIndexes(context.Context) error
	// CreateTestReportTable defines a function that creates the testreports table.
	CreateTestReportTable(context.Context, string) error

	// Artifacts Management Functions

	// CountTestReports returns the count of all test reports.
	CountTestReports(context.Context) (int64, error)

	// CountTestReportsByBuild returns the count of test reports by build ID.
	CountTestReportsByBuild(context.Context, *api.Build, map[string]interface{}, int64, int64) (int64, error)

	// CountTestReportsByRepo returns the count of test reports by repo ID.
	CountTestReportsByRepo(context.Context, *api.Repo, map[string]interface{}) (int64, error)

	// CreateTestReport creates a new test report.
	CreateTestReport(context.Context, *api.TestReport) (*api.TestReport, error)

	// DeleteTestReport removes a test report by ID.
	DeleteTestReport(context.Context, *api.TestReport) error

	// GetTestReport returns a test report by ID.
	GetTestReport(context.Context, int64) (*api.TestReport, error)

	// GetTestReportForBuild defines a function that gets a test report by number and build ID.
	GetTestReportForBuild(context.Context, *api.Build) (*api.TestReport, error)

	// ListTestReports returns a list of all test reports.
	ListTestReports(context.Context) ([]*api.TestReport, error)

	// ListTestReportsByBuild returns a list of test reports by build ID.
	ListTestReportsByBuild(context.Context, *api.Build, int, int) ([]*api.TestReport, error) // TODO atm,  there will only be 1 test report per build.

	// ListTestReportsByRepo returns a list of test reports by repo ID.
	ListTestReportsByRepo(context.Context, *api.Repo, int, int) ([]*api.TestReport, error)

	// UpdateTestReport updates a test report by ID.
	UpdateTestReport(context.Context, *api.TestReport) (*api.TestReport, error)
}
