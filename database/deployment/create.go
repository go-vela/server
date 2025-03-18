// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateDeployment creates a new deployment in the database.
func (e *engine) CreateDeployment(ctx context.Context, d *api.Deployment) (*api.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("creating deployment %d", d.GetID())

	// cast the API type to database type
	deployment := types.DeploymentFromAPI(d)

	// validate the necessary fields are populated
	err := deployment.Validate()
	if err != nil {
		return nil, err
	}

	err = e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Create(deployment).Error
	if err != nil {
		return nil, err
	}

	result := deployment.ToAPI(d.Builds)
	result.SetRepo(d.GetRepo())

	return result, nil
}
