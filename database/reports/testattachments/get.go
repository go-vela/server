// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// Get gets a test report attachment by ID from the database.
func (e *Engine) GetTestReportAttachment(ctx context.Context, id int64) (*api.TestReportAttachments, error) {
	e.logger.WithFields(logrus.Fields{
		"test_attachment_id": id,
	}).Tracef("getting test report attachment %d", id)

	// variable to store query results
	r := new(types.TestReportAttachment)

	// send query to the database
	err := e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, fmt.Errorf("unable to get test report attachment: %w", err)
	}

	return r.ToAPI(), nil
}
