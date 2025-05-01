// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateTestAttachment creates a new test attachment in the database.
func (e *Engine) CreateTestAttachment(ctx context.Context, r *api.TestAttachment) (*api.TestAttachment, error) {
	e.logger.WithFields(logrus.Fields{
		"test_report": r.GetID(),
	}).Tracef("creating test attachment %d", r.GetID())

	attachment := types.TestAttachmentFromAPI(r)

	err := attachment.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Create(attachment).Error
	if err != nil {
		return nil, err
	}

	result := attachment.ToAPI()
	result.SetTestReportID(r.GetTestReportID())

	return result, nil
}
