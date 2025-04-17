// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteByID deletes an existing test report from the database.
func (e *Engine) DeleteByID(ctx context.Context, r *api.TestReport) error {
	e.logger.WithFields(logrus.Fields{
		"test_report": r.GetID(),
	}).Tracef("deleting test report %d", r.GetID())

	// cast the API type to database type
	report := types.TestReportFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Delete(report).
		Error
}
