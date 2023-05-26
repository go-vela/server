// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CleanServices updates services to an error with a created timestamp prior to a defined moment.
func (e *engine) CleanServices(before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running steps in the database created prior to %d", before)

	s := new(library.Service)
	s.SetStatus(constants.StatusError)

	service := database.ServiceFromLibrary(s)

	// send query to the database
	result := e.client.
		Table(constants.TableService).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(service)

	return result.RowsAffected, result.Error
}
