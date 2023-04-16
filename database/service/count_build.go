// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountServicesForBuild gets the count of services by build ID from the database.
func (e *engine) CountServicesForBuild(b *library.Build, filters map[string]interface{}) (int64, error) {
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
