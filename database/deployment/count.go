// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountDeployments gets the count of all deployments from the database.
func (e *engine) CountDeployments(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all deployments from the database")

	// variable to store query results
	var d int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDeployment).
		Count(&d).
		Error

	return d, err
}
