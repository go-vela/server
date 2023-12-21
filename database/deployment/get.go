// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"context"
	"strconv"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetDeployment gets a deployment by ID from the database.
func (e *engine) GetDeployment(ctx context.Context, id int64) (*library.Deployment, error) {
	e.logger.Tracef("getting deployment %d from the database", id)

	// variable to store query results
	d := new(database.Deployment)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDeployment).
		Where("id = ?", id).
		Take(d).
		Error
	if err != nil {
		return nil, err
	}

	builds := []*library.Build{}

	for _, a := range d.Builds {
		bID, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			return nil, err
		}
		// variable to store query results
		b := new(database.Build)

		// send query to the database and store result in variable
		err2 := e.client.
			Table(constants.TableBuild).
			Where("id = ?", bID).
			Take(b).
			Error
		if err2 != nil {
			return nil, err
		}

		builds = append(builds, b.ToLibrary())
	}

	// return the deployment
	return d.ToLibrary(builds), nil
}
