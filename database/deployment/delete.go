// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeleteDeployment deletes an existing deployment from the database.
func (e *engine) DeleteDeployment(ctx context.Context, d *library.Deployment) error {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("deleting deployment %d", d.GetID())

	// cast the library type to database type
	deployment := database.DeploymentFromLibrary(d)

	// send query to the database
	return e.client.
		Table(constants.TableDeployment).
		Delete(deployment).
		Error
}
