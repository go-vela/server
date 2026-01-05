// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListTestAttachmentsByBuildID returns a list of test attachments for a specific build ID from the database.
func (e *Engine) ListTestAttachmentsByBuildID(ctx context.Context, buildID int64) ([]*api.TestAttachment, error) {
	e.logger.WithFields(logrus.Fields{
		"build_id": buildID,
	}).Trace("listing test attachments for build from the database")

	// variables to store query results and return value
	t := new([]types.TestAttachment)

	var attachments []*api.TestAttachment

	// send query to the database and store result in variable
	// join with testreports table since testattachments references testreports via test_report_id
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Joins("JOIN testreports ON testattachments.test_report_id = testreports.id").
		Where("testreports.build_id = ?", buildID).
		Order("testattachments.created_at DESC").
		Find(&t).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, attachment := range *t {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := attachment

		attachments = append(attachments, tmp.ToAPI())
	}

	return attachments, nil
}
