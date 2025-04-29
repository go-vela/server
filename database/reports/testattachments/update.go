// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// Update updates an existing test report in the database.
func (e *Engine) Update(ctx context.Context, t *api.TestReportAttachments) (*api.TestReportAttachments, error) {
	e.logger.WithFields(logrus.Fields{
		"testattchment": t.GetID(),
	}).Tracef("updating test report attachment %d in the database", t.GetID())

	testReportAttachment := types.TestReportAttachmentFromAPI(t)

	err := testReportAttachment.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Save(testReportAttachment)

	return testReportAttachment.ToAPI(), result.Error
}
