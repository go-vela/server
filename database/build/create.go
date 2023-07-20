// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with update.go
package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateBuild creates a new build in the database.
func (e *engine) CreateBuild(b *library.Build) (*library.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("creating build %d in the database", b.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildFromLibrary
	build := database.BuildFromLibrary(b)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Build.Validate
	err := build.Validate()
	if err != nil {
		return nil, err
	}

	// crop build if any columns are too large
	build = build.Crop()

	// send query to the database
	result := e.client.Table(constants.TableBuild).Create(build)

	return build.ToLibrary(), result.Error
}
