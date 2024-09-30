// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// DeleteDeployment deletes an existing deployment from the database.
func (e *engine) DeleteDeployment(ctx context.Context, d *api.Deployment) error {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("deleting deployment %d", d.GetID())

	// cast the library type to database type
	deployment := types.DeploymentFromAPI(d)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Delete(deployment).
		Error
}
