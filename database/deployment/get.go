// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"strconv"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetDeployment gets a deployment by ID from the database.
func (e *engine) GetDeployment(ctx context.Context, id int64) (*api.Deployment, error) {
	e.logger.Tracef("getting deployment %d", id)

	// variable to store query results
	d := new(types.Deployment)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("id = ?", id).
		Take(d).
		Error
	if err != nil {
		return nil, err
	}

	builds := []*api.Build{}

	for _, a := range d.Builds {
		bID, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			return nil, err
		}
		// variable to store query results
		b := new(types.Build)

		// send query to the database and store result in variable
		err = e.client.
			WithContext(ctx).
			Table(constants.TableBuild).
			Where("id = ?", bID).
			Take(b).
			Error
		if err != nil {
			return nil, err
		}

		builds = append(builds, b.ToAPI())
	}

	err = d.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo: %v", err)
	}

	// return the deployment
	return d.ToAPI(builds), nil
}
