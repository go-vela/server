// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountLogsForBuild gets the count of logs by build ID from the database.
func (e *Engine) CountLogsForBuild(ctx context.Context, b *api.Build) (int64, error) {
	e.logger.Tracef("getting count of logs for build %d", b.GetID())

	// variable to store query results
	var l int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Count(&l).
		Error

	return l, err
}
