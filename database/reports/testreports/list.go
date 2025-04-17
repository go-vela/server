// SPDX-License-Identifier: Apache-2.0

package testreports

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// List returns a list of test reports from the database.
func (e *Engine) List(ctx context.Context, page, perPage int) ([]*api.TestReport, int64, error) {
	e.logger.Trace("listing test reports from the database")

	// variables to store query results and return value
	t := new([]types.TestReport)
	var reports []*api.TestReport

	// calculate offset for pagination
	offset := (page - 1) * perPage

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Order("created DESC").
		Limit(perPage).
		Offset(offset).
		Find(&t).
		Error
	if err != nil {
		return nil, 0, err
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	// get the total count of reports
	var count int64
	err = e.client.
		WithContext(ctx).
		Table(constants.TableTestReports).
		Count(&count).
		Error
	if err != nil {
		return nil, 0, err
	}

	return reports, count, nil
}
