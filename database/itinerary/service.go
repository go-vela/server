// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import "github.com/go-vela/types/library"

// BuildItineraryService represents the Vela interface for build itinerary
// functions with the supported Database backends.
type BuildItineraryService interface {
	// BuildItinerary Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateBuildItineraryTable(string) error

	// // BuildItinerary Data Manipulation Language Functions
	// //
	// // https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateBuildItinerary defines a function that creates a build itinerary.
	CreateBuildItinerary(*library.BuildItinerary) error
	// PopBuildItinerary defines a function that gets and deletes a build itinerary.
	PopBuildItinerary(int64) (*library.BuildItinerary, error)
}
