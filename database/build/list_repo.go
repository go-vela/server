// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListBuildsForRepo gets a list of builds by repo ID from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ListBuildsForRepo(ctx context.Context, r *api.Repo, filters map[string]any, before, after int64, page, perPage int) ([]*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing builds for repo %s", r.GetFullName())

	// variables to store query results and return values
	b := new([]types.Build)
	builds := []*api.Build{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	err := e.client.
		WithContext(ctx).
		Table(constants.TableBuild).
		Where("repo_id = ?", r.GetID()).
		Where("created < ?", before).
		Where("created > ?", after).
		Where(filters).
		Order("number DESC").
		Limit(perPage).
		Offset(offset).
		Find(&b).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		result := tmp.ToAPI()
		result.SetRepo(r)

		builds = append(builds, result)
	}

	return builds, nil
}
