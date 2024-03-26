// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// ListRepos gets a list of all repos from the database.
func (e *engine) ListRepos(ctx context.Context) ([]*api.Repo, error) {
	e.logger.Trace("listing all repos from the database")

	// variables to store query results and return value
	count := int64(0)
	r := new([]Repo)
	repos := []*api.Repo{}

	// count the results
	count, err := e.CountRepos(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return repos, nil
	}

	// send query to the database and store result in variable
	err = e.client.
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

		// decrypt owner fields
		err = tmp.Owner.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo owner %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		repos = append(repos, tmp.ToAPI())
	}

	return repos, nil
}
