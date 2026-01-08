// SPDX-License-Identifier: Apache-2.0

package artifact

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountArtifacts gets the count of all artifacts from the database.
func (e *Engine) CountArtifacts(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all artifacts")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableArtifact).
		Count(&s).
		Error

	return s, err
}
