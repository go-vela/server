// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	"fmt"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// ListByRepo returns a list of test reports by repo ID from the database.
func (e *Engine) ListByRepo(ctx context.Context, r *api.Repo, page, perPage int) ([]*api.TestReport, int64, error) {
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
		Where("repo_id = ?", r.GetID()).
		Order("created DESC").
		Limit(perPage).
		Offset(offset).
		Find(&t).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list test reports by repo ID: %w", err)
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	// get the total count of reports for this repo
	var count int64
	err = e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Where("repo_id = ?", r.GetID()).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count test reports by repo ID: %w", err)
	}

	return reports, count, nil
}
