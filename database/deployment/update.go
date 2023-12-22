// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateDeployment updates an existing deployment in the database.
func (e *engine) UpdateDeployment(ctx context.Context, d *library.Deployment) (*library.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("updating deployment %d in the database", d.GetID())

	// cast the library type to database type
	deployment := database.DeploymentFromLibrary(d)

	// validate the necessary fields are populated
	err := deployment.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.Table(constants.TableDeployment).Save(deployment)

	// send query to the database
	return deployment.ToLibrary(d.Builds), result.Error
}
