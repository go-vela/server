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

// DeleteBuild deletes an existing build from the database.
func (e *engine) DeleteBuild(r *library.Build) error {
	e.logger.WithFields(logrus.Fields{
		"org":   r.GetOrg(),
		"build": r.GetName(),
	}).Tracef("deleting build %s from the database", r.GetFullName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildFromLibrary
	build := database.BuildFromLibrary(r)

	// send query to the database
	return e.client.
		Table(constants.TableBuild).
		Delete(build).
		Error
}
