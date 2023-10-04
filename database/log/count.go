// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountLogs gets the count of all logs from the database.
func (e *engine) CountLogs(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all logs from the database")

	// variable to store query results
	var l int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableLog).
		Count(&l).
		Error

	return l, err
}
