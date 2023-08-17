// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteRepo deletes an existing repo from the database.
func (e *engine) DeleteRepo(ctx context.Context, r *library.Repo) error {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("deleting repo %s from the database", r.GetFullName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#RepoFromLibrary
	repo := database.RepoFromLibrary(r)

	// send query to the database
	return e.client.
		Table(constants.TableRepo).
		Delete(repo).
		Error
}
