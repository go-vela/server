// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// DeleteBuild deletes an existing build from the database.
func (e *engine) DeleteBuild(ctx context.Context, b *api.Build) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("deleting build %d from the database", b.GetNumber())

	build := FromAPI(b)

	// send query to the database
	return e.client.
		Table(constants.TableBuild).
		Delete(build).
		Error
}
