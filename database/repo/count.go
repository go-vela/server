// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountRepos gets the count of all repos from the database.
func (e *Engine) CountRepos(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all repos")

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Count(&r).
		Error

	return r, err
}
