// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountServicesForBuild gets the count of services by build ID from the database.
func (e *engine) CountServicesForBuild(ctx context.Context, b *library.Build, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of services for build %d from the database", b.GetNumber())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Where("build_id = ?", b.GetID()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
