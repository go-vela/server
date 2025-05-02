// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListReposForOrg gets a list of repos by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ReposInList(ctx context.Context, nameList []string) ([]*api.Repo, error) {
	// variables to store query results and return values
	r := new([]types.Repo)
	repos := []*api.Repo{}

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
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// convert query result to API type
		repos = append(repos, tmp.ToAPI())
	}

	return repos, nil
}
