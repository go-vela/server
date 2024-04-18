// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListBuildsForOrg gets a list of builds by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListBuildsForOrg(ctx context.Context, org string, filters map[string]interface{}, page, perPage int) ([]*api.Build, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing builds for org %s from the database", org)

	// variables to store query results and return values
	count := int64(0)
	b := new([]types.Build)
	builds := []*api.Build{}

	// count the results
	count, err := e.CountBuildsForOrg(ctx, org, filters)
	if err != nil {
		return builds, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return builds, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err = e.client.
		Table(constants.TableBuild).
		Preload("Repo").
		Preload("Repo.Owner").
		Select("builds.*").
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repos.org = ?", org).
		Where(filters).
		Order("created DESC").
		Order("id").
		Limit(perPage).
		Offset(offset).
		Find(&b).
		Error
	if err != nil {
		return nil, count, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		err = tmp.Repo.Decrypt(e.config.EncryptionKey)
		if err != nil {
			e.logger.Errorf("unable to decrypt repo: %v", err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Build.ToLibrary
		builds = append(builds, tmp.ToAPI())
	}

	return builds, count, nil
}
