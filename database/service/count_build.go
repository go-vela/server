// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// CountServicesForBuild gets the count of services by build ID from the database.
func (e *Engine) CountServicesForBuild(ctx context.Context, b *api.Build, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of services for build %d", b.GetNumber())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableService).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
