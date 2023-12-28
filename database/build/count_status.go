// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountBuildsForStatus gets the count of builds by status from the database.
func (e *engine) CountBuildsForStatus(ctx context.Context, status string, filters map[string]interface{}) (int64, error) {
	e.logger.Tracef("getting count of builds for status %s from the database", status)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("status = ?", status).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
