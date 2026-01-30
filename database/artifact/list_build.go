// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListArtifactsByBuildID returns a list of artifacts for a specific build ID from the database.
func (e *Engine) ListArtifactsByBuildID(ctx context.Context, buildID int64) ([]*api.Artifact, error) {
	e.logger.WithFields(logrus.Fields{
		"build_id": buildID,
	}).Trace("listing artifacts for build from the database")

	// variables to store query results and return value
	t := new([]types.Artifact)

	var artifacts []*api.Artifact

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Where("build_id = ?", buildID).
		Order("artifacts.created_at DESC").
		Find(&t).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, artifact := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := artifact

		artifacts = append(artifacts, tmp.ToAPI())
	}

	return artifacts, nil
}
