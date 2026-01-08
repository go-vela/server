// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteArtifact deletes an existing artifact from the database.
func (e *Engine) DeleteArtifact(ctx context.Context, r *api.Artifact) error {
	e.logger.WithFields(logrus.Fields{
		"artifact": r.GetID(),
	}).Tracef("deleting artifact %d", r.GetID())

	// cast the API type to database type
	artifact := types.ArtifactFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Delete(artifact).
		Error
}
