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
	"github.com/sirupsen/logrus"
)

// ListDeploymentsForRepo gets a list of deployments by repo ID from the database.
func (e *engine) ListDeploymentsForRepo(ctx context.Context, r *library.Repo, page, perPage int) ([]*library.Deployment, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing deployments for repo %s from the database", r.GetFullName())

	// variables to store query results and return value
	d := new([]database.Deployment)
	deployments := []*library.Deployment{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableDeployment).
		Where("repo_id = ?", r.GetID()).
		Order("id DESC").
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

		builds := []*library.Build{}

		for _, a := range tmp.Builds {
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

		// convert query result to library type
		deployments = append(deployments, tmp.ToLibrary(builds))
	}

	return deployments, nil
}
