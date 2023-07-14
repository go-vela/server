// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListRepos gets a list of all repos from the database.
func (e *engine) ListRepos() ([]*library.Repo, error) {
	e.logger.Trace("listing all repos from the database")

	// variables to store query results and return value
	count := int64(0)
	r := new([]database.Repo)
	repos := []*library.Repo{}

	// count the results
	count, err := e.CountRepos()
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
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted repos
			e.logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, nil
}
