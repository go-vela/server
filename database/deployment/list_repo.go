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

// ListDeploymentsForRepo gets a list of deployments by repo ID from the database.
func (e *engine) ListDeploymentsForRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing deployments for repo %s", r.GetFullName())

	// variables to store query results and return value
	d := new([]types.Deployment)
	deployments := []*api.Deployment{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableDeployment).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("repo_id = ?", r.GetID()).
		Order("number DESC").
		Limit(perPage).
		Offset(offset).
		Find(&d).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, deployment := range *d {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := deployment

		builds := []*api.Build{}

		for _, a := range tmp.Builds {
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

		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo %s/%s: %v", r.GetOrg(), r.GetName(), err)
		}

		// convert query result to API type
		deployments = append(deployments, tmp.ToAPI(builds))
	}

	return deployments, nil
}
