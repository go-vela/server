// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateBuildItinerary creates a new build itinerary in the database.
func (e *engine) CreateBuildItinerary(b *library.BuildItinerary) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetBuildID(),
	}).Tracef("creating build itinerary for build %d in the database", b.GetBuildID())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildItineraryFromLibrary
	compiled := database.BuildItineraryFromLibrary(b)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildItinerary.Validate
	err := compiled.Validate()
	if err != nil {
		return err
	}

	// compress data for the build itinerary
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildItinerary.Compress
	err = compiled.Compress(e.config.CompressionLevel)
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableBuildItinerary).
		Create(compiled).
		Error
}
