// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateTestReport creates a new test report in the database.
func (e *Engine) CreateTestReport(ctx context.Context, r *api.TestReport) (*api.TestReport, error) {
	e.logger.WithFields(logrus.Fields{
		"test_report": r.GetID(),
	}).Tracef("creating test report %d", r.GetID())

	report := types.TestReportFromAPI(r)

	err := report.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableTestReport).
		Create(report).Error
	if err != nil {
		return nil, err
	}

	result := report.ToAPI()
	result.SetBuildID(r.GetBuildID())

	return result, nil
}
