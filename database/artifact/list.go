// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListArtifacts returns a list of artifacts from the database.
func (e *Engine) ListArtifacts(ctx context.Context) ([]*api.Artifact, error) {
	e.logger.Trace("listing artifacts from the database")

	// variables to store query results and return value
	t := new([]types.Artifact)

	var reports []*api.Artifact

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Order("created_at DESC").
		Find(&t).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	return reports, nil
}
