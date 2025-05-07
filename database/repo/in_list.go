// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetReposInList gets a list of repos from the database from a list of full_names.
func (e *Engine) GetReposInList(ctx context.Context, nameList []string) ([]*api.Repo, error) {
	e.logger.Tracef("getting repos in list %v", nameList)

	// variables to store query results and return value
	r := new([]types.Repo)
	repos := []*api.Repo{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Where("full_name IN (?)", nameList).
		Find(&r).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, repo := range *r {
		tmp := repo

		repos = append(repos, tmp.ToAPI())
	}

	return repos, nil
}
