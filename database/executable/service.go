// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executable

import "github.com/go-vela/types/library"

// BuildExecutableService represents the Vela interface for build executable
// functions with the supported Database backends.
type BuildExecutableService interface {
	// BuildExecutable Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateBuildExecutableTable(string) error

	// // BuildExecutable Data Manipulation Language Functions
	// //
	// // https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateBuildExecutable defines a function that creates a build executable.
	CreateBuildExecutable(*library.BuildExecutable) error
	// PopBuildExecutable defines a function that gets and deletes a build executable.
	PopBuildExecutable(int64) (*library.BuildExecutable, error)
}
