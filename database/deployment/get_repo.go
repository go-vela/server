// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetDeploymentForRepo gets a deployment by repo ID and number from the database.
func (e *engine) GetDeploymentForRepo(ctx context.Context, r *api.Repo, number int64) (*library.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"deployment": number,
		"org":        r.GetOrg(),
		"repo":       r.GetName(),
	}).Tracef("getting deployment %s/%d", r.GetFullName(), number)

	// variable to store query results
	d := new(database.Deployment)

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

	builds := []*library.Build{}

	for _, a := range d.Builds {
		bID, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			return nil, err
		}
		// variable to store query results
		b := new(database.Build)

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

		builds = append(builds, b.ToLibrary())
	}

	return d.ToLibrary(builds), nil
}
