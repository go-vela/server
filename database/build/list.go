// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListBuilds gets a list of all builds from the database.
func (e *engine) ListBuilds() ([]*library.Build, error) {
	e.logger.Trace("listing all builds from the database")

	// variables to store query results and return value
	count := int64(0)
	r := new([]database.Build)
	builds := []*library.Build{}

	// count the results
	count, err := e.CountBuilds()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		Table(constants.TableBuild).
		Find(&r).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, build := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// decrypt the fields for the build
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted builds
			e.logger.Errorf("unable to decrypt build %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, nil
}
