// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// Update updates an existing test report in the database.
func (e *Engine) Update(ctx context.Context, t *api.TestReport) (*api.TestReport, error) {
	e.logger.WithFields(logrus.Fields{
		"testreport": t.GetID(),
	}).Tracef("updating test report %d in the database", t.GetID())

	testReport := types.TestReportFromAPI(t)

	err := testReport.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Save(testReport)

	return testReport.ToAPI(), result.Error
}
