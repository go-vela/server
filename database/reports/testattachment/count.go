// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/go-vela/server/constants"
)

// CountTestAttachments gets the count of all test attachments from the database.
func (e *Engine) CountTestAttachments(ctx context.Context) (int64, error) {
	e.logger.Tracef("getting count of all test attachments")

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Count(&s).
		Error

	return s, err
}
