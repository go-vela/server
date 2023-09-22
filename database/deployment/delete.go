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

// DeleteDeployment deletes an existing deployment from the database.
func (e *engine) DeleteDeployment(ctx context.Context, d *library.Deployment) error {
	e.logger.WithFields(logrus.Fields{
		"deployment": d.GetID(),
	}).Tracef("deleting deployment %d in the database", d.GetID())

	// cast the library type to database type
	deployment := database.DeploymentFromLibrary(d)

	// send query to the database
	return e.client.
		Table(constants.TableDeployment).
		Delete(deployment).
		Error
}
