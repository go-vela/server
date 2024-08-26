// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// UpdateDeployment updates an existing deployment in the database.
func (e *engine) UpdateDeployment(ctx context.Context, d *library.Deployment) (*library.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("updating deployment %d", d.GetID())

	// cast the library type to database type
	deployment := database.DeploymentFromLibrary(d)

	// validate the necessary fields are populated
	err := deployment.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Save(deployment)

	// send query to the database
	return deployment.ToLibrary(d.Builds), result.Error
}
