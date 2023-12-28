// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// CountLogsForBuild gets the count of logs by build ID from the database.
func (e *engine) CountLogsForBuild(ctx context.Context, b *library.Build) (int64, error) {
	e.logger.Tracef("getting count of logs for build %d from the database", b.GetID())

	// variable to store query results
	var l int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Count(&l).
		Error

	return l, err
}
