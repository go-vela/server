// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateArtifact creates a new artifact in the database.
func (e *Engine) CreateArtifact(ctx context.Context, r *api.Artifact) (*api.Artifact, error) {
	e.logger.WithFields(logrus.Fields{
		"build": r.GetID(),
	}).Tracef("creating artifact %d", r.GetID())

	artifact := types.ArtifactFromAPI(r)

	err := artifact.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Create(artifact).Error
	if err != nil {
		return nil, err
	}

	result := artifact.ToAPI()
	result.SetBuildID(r.GetBuildID())

	return result, nil
}
