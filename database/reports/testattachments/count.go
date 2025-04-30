// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"

	"github.com/go-vela/server/constants"
)

// Count gets the count of all test report attachments from the database.
func (e *Engine) Count(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all test report attachments")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Count(&s).
		Error

	return s, err
}
