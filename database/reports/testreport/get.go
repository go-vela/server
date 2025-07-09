// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetTestReport gets a test report by ID from the database.
func (e *Engine) GetTestReport(ctx context.Context, id int64) (*api.TestReport, error) {
	e.logger.WithFields(logrus.Fields{
		"test_report": id,
	}).Tracef("getting test report %d", id)

	// variable to store query results
	r := new(types.TestReport)

	// send query to the database
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReport).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, fmt.Errorf("unable to get test report: %w", err)
	}

	return r.ToAPI(), nil
}
