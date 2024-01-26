// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountStepsForBuild gets the count of steps by build ID from the database.
func (e *engine) CountStepsForBuild(ctx context.Context, b *library.Build, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of steps for build %d from the database", b.GetNumber())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableStep).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
