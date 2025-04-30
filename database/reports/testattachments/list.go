// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// List returns a list of test reports from the database.
func (e *Engine) List(ctx context.Context) ([]*api.TestReportAttachments, error) {
	e.logger.Trace("listing test report attachments from the database")

	// variables to store query results and return value
	t := new([]types.TestReportAttachment)

	var reports []*api.TestReportAttachments

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Order("created DESC").
		Find(&t).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, report := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := report

		reports = append(reports, tmp.ToAPI())
	}

	return reports, nil
}
