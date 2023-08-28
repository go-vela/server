// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executable

import (
	"context"

	"github.com/go-vela/types/library"
)

// BuildExecutableInterface represents the Vela interface for build executable
// functions with the supported Database backends.
type BuildExecutableInterface interface {
	// BuildExecutable Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateBuildExecutableTable(context.Context, string) error

	// BuildExecutable Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CreateBuildExecutable defines a function that creates a build executable.
	CreateBuildExecutable(context.Context, *library.BuildExecutable) error
	// PopBuildExecutable defines a function that gets and deletes a build executable.
	PopBuildExecutable(context.Context, int64) (*library.BuildExecutable, error)
}
