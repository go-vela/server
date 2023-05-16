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

// GetBuildForOrg gets a build by org and build name from the database.
func (e *engine) GetBuildForOrg(org, name string) (*library.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":   org,
		"build": name,
	}).Tracef("getting build %s/%s from the database", org, name)

	// variable to store query results
	r := new(database.Build)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("org = ?", org).
		Where("name = ?", name).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the build
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Build.Decrypt
	err = r.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted builds
		e.logger.Errorf("unable to decrypt build %s/%s: %v", org, name, err)

		// return the unencrypted build
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		return r.ToLibrary(), nil
	}

	// return the decrypted build
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
	return r.ToLibrary(), nil
}
