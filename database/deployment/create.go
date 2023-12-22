// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateDeployment creates a new deployment in the database.
func (e *engine) CreateDeployment(ctx context.Context, d *library.Deployment) (*library.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("creating deployment %d in the database", d.GetID())

	// cast the library type to database type
	deployment := database.DeploymentFromLibrary(d)

	// validate the necessary fields are populated
	err := deployment.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.Table(constants.TableDeployment).Create(deployment)

	// send query to the database
	return deployment.ToLibrary(d.Builds), result.Error
}
