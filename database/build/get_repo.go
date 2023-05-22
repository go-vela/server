// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetBuildForRepo gets a build by repo ID and number from the database.
func (e *engine) GetBuildForRepo(r *library.Repo, number int) (*library.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": number,
		"org":   r.GetOrg(),
		"repo":  r.GetName(),
	}).Tracef("getting build %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("number = ?", number).
		Take(b).
		Error
	if err != nil {
		return nil, err
	}

	return b.ToLibrary(), nil
}
