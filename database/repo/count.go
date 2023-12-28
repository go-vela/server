// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/go-vela/types/constants"
)

// CountRepos gets the count of all repos from the database.
func (e *engine) CountRepos(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all repos from the database")

	// variable to store query results
	var r int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Count(&r).
		Error

	return r, err
}
