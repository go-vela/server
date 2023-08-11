// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"context"
	"errors"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// LastBuildForRepo gets the last build by repo ID and branch from the database.
func (e *engine) LastBuildForRepo(ctx context.Context, r *library.Repo, branch string) (*library.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last build for repo %s from the database", r.GetFullName())

	// variable to store query results
	b := new(database.Build)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("branch = ?", branch).
		Order("number DESC").
		Take(b).
		Error
	if err != nil {
		// check if the query returned a record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// the record will not exist if it is a new repo
			return nil, nil
		}

		return nil, err
	}

	return b.ToLibrary(), nil
}
