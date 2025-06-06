// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountBuilds gets the count of all builds from the database.
func (e *Engine) CountBuilds(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all builds")

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Count(&b).
		Error

	return b, err
}
