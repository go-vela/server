// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteTestAttachment deletes an existing test attachment from the database.
func (e *Engine) DeleteTestAttachment(ctx context.Context, r *api.TestAttachment) error {
	e.logger.WithFields(logrus.Fields{
		"test_attachment": r.GetID(),
	}).Tracef("deleting test attachment %d", r.GetID())

	// cast the API type to database type
	attachment := types.TestAttachmentFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Delete(attachment).
		Error
}
