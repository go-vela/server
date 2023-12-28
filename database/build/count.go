// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountBuilds gets the count of all builds from the database.
func (e *engine) CountBuilds(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all builds from the database")

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Count(&b).
		Error

	return b, err
}
