// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CountBuildsForSender gets the count of builds by sender from the database.
func (e *engine) CountBuildsForSender(ctx context.Context, sender string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"sender": sender,
	}).Tracef("getting count of builds for sender %s from the database", sender)

	// variable to store query results
	var b int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableBuild).
		Where("sender = ?", sender).
		Where(filters).
		Count(&b).
		Error

	return b, err
}
