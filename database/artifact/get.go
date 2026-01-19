// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetArtifact gets an artifact by ID from the database.
func (e *Engine) GetArtifact(ctx context.Context, id int64) (*api.Artifact, error) {
	e.logger.WithFields(logrus.Fields{
		"artifact_id": id,
	}).Tracef("getting artifact %d", id)

	// variable to store query results
	r := new(types.Artifact)

	// send query to the database
	err := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, fmt.Errorf("unable to get artifact: %w", err)
	}

	return r.ToAPI(), nil
}
