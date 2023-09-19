// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListReposForOrg gets a list of repos by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListReposForOrg(ctx context.Context, org, sortBy string, filters map[string]interface{}, page, perPage int) ([]*library.Repo, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing repos for org %s from the database", org)

	// variables to store query results and return values
	count := int64(0)
	r := new([]database.Repo)
	repos := []*library.Repo{}

	// count the results
	count, err := e.CountReposForOrg(ctx, org, filters)
	if err != nil {
		return repos, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return repos, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	switch sortBy {
	case "latest":
		query := e.client.
			Table(constants.TableBuild).
			Select("repos.id, MAX(builds.created) AS latest_build").
			Joins("INNER JOIN repos repos ON builds.repo_id = repos.id").
			Where("repos.org = ?", org).
			Group("repos.id")

		err = e.client.
			Table(constants.TableRepo).
			Select("repos.*").
			Joins("LEFT JOIN (?) t on repos.id = t.id", query).
			Order("latest_build DESC NULLS LAST").
			Limit(perPage).
			Offset(offset).
			Find(&r).
			Error
		if err != nil {
			return nil, count, err
		}
	case "name":
		fallthrough
	default:
		err = e.client.
			Table(constants.TableRepo).
			Where("org = ?", org).
			Where(filters).
			Order("name").
			Limit(perPage).
			Offset(offset).
			Find(&r).
			Error
		if err != nil {
			return nil, count, err
		}
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

	return repos, count, nil
}
