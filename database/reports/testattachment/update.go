// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateTestAttachment updates an existing test report in the database.
func (e *Engine) UpdateTestAttachment(ctx context.Context, t *api.TestAttachment) (*api.TestAttachment, error) {
	e.logger.WithFields(logrus.Fields{
		"testattchment": t.GetID(),
	}).Tracef("updating test attachment %d in the database", t.GetID())

	testAttachment := types.TestAttachmentFromAPI(t)

	err := testAttachment.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Save(testAttachment)

	return testAttachment.ToAPI(), result.Error
}
