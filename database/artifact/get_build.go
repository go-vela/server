// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetArtifactForBuild gets an artifact by number and build ID from the database.
func (e *Engine) GetArtifactForBuild(ctx context.Context, b *api.Build) (*api.Artifact, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting artifact")

	// variable to store query results
	tr := new(types.Artifact)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Where("build_id = ?", b.GetID()).
		Take(tr).
		Error
	if err != nil {
		return nil, err
	}

	return tr.ToAPI(), nil
}
