// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountSteps gets the count of all steps from the database.
func (e *Engine) CountSteps(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all steps")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Count(&s).
		Error

	return s, err
}
