// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListByRepo returns a list of test reports by repo ID from the database.
func (e *Engine) ListByRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.TestReport, error) {
	e.logger.WithFields(logrus.Fields{
		"repo_id": r.GetID(),
	}).Tracef("listing test reports by repo ID %d", r.GetID())

	// variables to store query results and return value
	t := new([]types.TestReport)
	reports := []*api.TestReport{}

	// calculate offset for pagination
	offset := (page - 1) * perPage

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Preload("Build").
		Preload("Build.Repo").
		Preload("Build.Repo.Owner").
		Select("testreports.*").
		Joins("JOIN builds ON testreports.build_id = builds.id").
		Joins("JOIN repos ON builds.repo_id = repos.id").
		Where("repo_id = ?", r.GetID()).
		Order("created DESC").
		Limit(perPage).
		Offset(offset).
		Find(&t).
		Error
	if err != nil {
		return nil, fmt.Errorf("unable to list test reports by repo ID: %w", err)
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	return reports, nil
}
