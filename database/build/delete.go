// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteBuild deletes an existing build from the database.
func (e *engine) DeleteBuild(ctx context.Context, b *library.Build) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("deleting build %d from the database", b.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildFromLibrary
	build := database.BuildFromLibrary(b)

	// send query to the database
	return e.client.
		Table(constants.TableBuild).
		Delete(build).
		Error
}
