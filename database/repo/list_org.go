// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListReposForOrg gets a list of repos by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListReposForOrg(ctx context.Context, org, sortBy string, filters map[string]interface{}, page, perPage int) ([]*api.Repo, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing repos for org %s", org)

	// variables to store query results and return values
	count := int64(0)
	r := new([]types.Repo)
	repos := []*api.Repo{}

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
			WithContext(ctx).
			Table(constants.TableBuild).
			Select("repos.id, MAX(builds.created) AS latest_build").
			Joins("INNER JOIN repos repos ON builds.repo_id = repos.id").
			Where("repos.org = ?", org).
			Group("repos.id")

		err = e.client.
			WithContext(ctx).
			Table(constants.TableRepo).
			Preload("Owner").
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
			WithContext(ctx).
			Table(constants.TableRepo).
			Preload("Owner").
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

	return repos, count, nil
}
