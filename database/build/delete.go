// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteBuild deletes an existing build from the database.
func (e *Engine) DeleteBuild(ctx context.Context, b *api.Build) error {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("deleting build %d", b.GetNumber())

	build := types.BuildFromAPI(b)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Delete(build).
		Error
}
