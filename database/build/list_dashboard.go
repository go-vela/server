// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListBuildsForDashboardRepo gets a list of builds by repo ID from the database.
func (e *engine) ListBuildsForDashboardRepo(ctx context.Context, r *api.Repo, branches, events []string) ([]*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing builds for repo %s", r.GetFullName())

	// variables to store query results and return values
	b := new([]types.Build)
	builds := []*api.Build{}

	query := e.client.Table(constants.TableBuild).Where("repo_id = ?", r.GetID())

	if len(branches) > 0 {
		query = query.Where("branch IN (?)", branches)
	}

	if len(events) > 0 {
		query = query.Where("event IN (?)", events)
	}

	err := query.
		Order("number DESC").
		Limit(5).
		Find(&b).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		builds = append(builds, tmp.ToAPI())
	}

	return builds, nil
}
