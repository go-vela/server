// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListRepos gets a list of all repos from the database.
func (e *Engine) ListRepos(ctx context.Context) ([]*api.Repo, error) {
	e.logger.Trace("listing all repos")

	// variables to store query results and return value
	r := new([]types.Repo)
	repos := []*api.Repo{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Preload("Owner").
		Find(&r).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// decrypt the fields for the repo
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted repos
			e.logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to API type
		repos = append(repos, tmp.ToAPI())
	}

	return repos, nil
}
