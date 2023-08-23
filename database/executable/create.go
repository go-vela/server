// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executable

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateBuildExecutable creates a new build executable in the database.
func (e *engine) CreateBuildExecutable(b *library.BuildExecutable) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetBuildID(),
	}).Tracef("creating build executable for build %d in the database", b.GetBuildID())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutableFromLibrary
	executable := database.BuildExecutableFromLibrary(b)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutable.Validate
	err := executable.Validate()
	if err != nil {
		return err
	}

	// compress data for the build executable
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutable.Compress
	err = executable.Compress(e.config.CompressionLevel)
	if err != nil {
		return err
	}

	// encrypt the data field for the build executable
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutable.Encrypt
	err = executable.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt build executable for build %d: %w", b.GetBuildID(), err)
	}

	// send query to the database
	return e.client.
		Table(constants.TableBuildExecutable).
		Create(executable).
		Error
}
