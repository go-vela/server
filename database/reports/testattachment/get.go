// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetTestAttachment gets a test attachment by ID from the database.
func (e *Engine) GetTestAttachment(ctx context.Context, id int64) (*api.TestAttachment, error) {
	e.logger.WithFields(logrus.Fields{
		"test_attachment_id": id,
	}).Tracef("getting test attachment %d", id)

	// variable to store query results
	r := new(types.TestAttachment)

	// send query to the database
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, fmt.Errorf("unable to get test attachment: %w", err)
	}

	return r.ToAPI(), nil
}
