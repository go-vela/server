// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetDeploymentForRepo gets a deployment by repo ID and number from the database.
func (e *engine) GetDeploymentForRepo(ctx context.Context, r *api.Repo, number int64) (*api.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": number,
		"org":        r.GetOrg(),
		"repo":       r.GetName(),
	}).Tracef("getting deployment %s/%d", r.GetFullName(), number)

	// variable to store query results
	d := new(types.Deployment)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
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

	result := d.ToAPI(builds)
	result.SetRepo(r)

	return result, nil
}
