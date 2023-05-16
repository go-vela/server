// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with create.go
package build

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateBuild updates an existing build in the database.
func (e *engine) UpdateBuild(r *library.Build) error {
	e.logger.WithFields(logrus.Fields{
		"org":   r.GetOrg(),
		"build": r.GetName(),
	}).Tracef("creating build %s in the database", r.GetFullName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildFromLibrary
	build := database.BuildFromLibrary(r)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Build.Validate
	err := build.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the build
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Build.Encrypt
	err = build.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt build %s: %w", r.GetFullName(), err)
	}

	// send query to the database
	return e.client.
		Table(constants.TableBuild).
		Save(build).
		Error
}
