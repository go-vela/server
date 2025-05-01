// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListTestAttachments returns a list of test reports from the database.
func (e *Engine) ListTestAttachments(ctx context.Context) ([]*api.TestAttachment, error) {
	e.logger.Trace("listing test attachments from the database")

	// variables to store query results and return value
	t := new([]types.TestAttachment)

	var reports []*api.TestAttachment

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Order("created_at DESC").
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
