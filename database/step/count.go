// SPDX-License-Identifier: Apache-2.0

package step

import (
	"github.com/go-vela/types/constants"
)

// CountSteps gets the count of all steps from the database.
func (e *engine) CountSteps() (int64, error) {
	e.logger.Tracef("getting count of all steps from the database")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableStep).
		Count(&s).
		Error

	return s, err
}
