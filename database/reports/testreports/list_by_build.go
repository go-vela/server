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

// ListByBuild returns a list of test reports by build ID from the database.
func (e *Engine) ListByBuild(ctx context.Context, b *api.Build, page, perPage int) ([]*api.TestReport, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build_id": b.GetID(),
	}).Tracef("listing test reports by build number %v", b.GetNumber())

	// variables to store query results and return value
	t := new([]types.TestReport)
	reports := []*api.TestReport{}

	// calculate offset for pagination
	offset := (page - 1) * perPage

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Where("build_id = ?", b.GetID()).
		Order("created DESC").
		Limit(perPage).
		Offset(offset).
		Find(&t).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to list test reports by build ID: %w", err)
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	// get the total count of reports for this build
	var count int64
	err = e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Where("build_id = ?", b.GetID()).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, fmt.Errorf("unable to count test reports by build ID: %w", err)
	}

	return reports, count, nil
}
