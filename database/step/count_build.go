// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountStepsForBuild gets the count of steps by build ID from the database.
func (e *Engine) CountStepsForBuild(ctx context.Context, b *api.Build, filters map[string]any) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of steps for build %d", b.GetNumber())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
