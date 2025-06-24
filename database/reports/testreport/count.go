// SPDX-License-Identifier: Apache-2.0

package testreport

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountTestReports gets the count of all test reports from the database.
func (e *Engine) CountTestReports(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all test reports")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReport).
		Count(&s).
		Error

	return s, err
}
