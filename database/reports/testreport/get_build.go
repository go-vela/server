// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetTestReportForBuild gets a test report by number and build ID from the database.
func (e *Engine) GetTestReportForBuild(ctx context.Context, b *api.Build) (*api.TestReport, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting testreport")

	// variable to store query results
	tr := new(types.TestReport)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReport).
		Where("build_id = ?", b.GetID()).
		Take(tr).
		Error
	if err != nil {
		return nil, err
	}

	return tr.ToAPI(), nil
}
