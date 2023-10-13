// SPDX-License-Identifier: Apache-2.0

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

	// CleanBuildExecutables defines a function that deletes errored builds' corresponding executables.
	CleanBuildExecutables(context.Context) error
	// CreateBuildExecutable defines a function that creates a build executable.
	CreateBuildExecutable(context.Context, *library.BuildExecutable) error
	// PopBuildExecutable defines a function that gets and deletes a build executable.
	PopBuildExecutable(context.Context, int64) (*library.BuildExecutable, error)
}
