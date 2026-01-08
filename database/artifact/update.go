// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateArtifact updates an existing artifact in the database.
func (e *Engine) UpdateArtifact(ctx context.Context, t *api.Artifact) (*api.Artifact, error) {
	e.logger.WithFields(logrus.Fields{
		"testattchment": t.GetID(),
	}).Tracef("updating artifact %d in the database", t.GetID())

	artifact := types.ArtifactFromAPI(t)

	err := artifact.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Save(artifact)

	return artifact.ToAPI(), result.Error
}
