// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteInitStep deletes an existing init step from the database.
func (e *engine) DeleteInitStep(i *library.InitStep) error {
	e.logger.WithFields(logrus.Fields{
		"initstep": i.GetNumber(),
	}).Tracef("deleting init step %d in the database", i.GetID())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitStepFromLibrary
	initStep := database.InitStepFromLibrary(i)

	// send query to the database
	return e.client.
		Table(constants.TableInitStep).
		Delete(initStep).
		Error
}
