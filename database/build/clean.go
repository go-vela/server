// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CleanBuilds updates builds to an error with a provided message with a created timestamp prior to a defined moment.
func (e *engine) CleanBuilds(msg string, before int64) (int64, error) {
	logrus.Tracef("cleaning pending or running builds in the database created prior to %d", before)

	b := new(library.Build)
	b.SetStatus(constants.StatusError)
	b.SetError(msg)

	build := database.BuildFromLibrary(b)

	// send query to the database
	result := e.client.
		Table(constants.TableBuild).
		Where("created < ?", before).
		Where("status = 'running' OR status = 'pending'").
		Updates(build)

	return result.RowsAffected, result.Error
}
