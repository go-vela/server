// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteInit deletes an existing init from the database.
func (e *engine) DeleteInit(i *library.Init) error {
	e.logger.WithFields(logrus.Fields{
		"init": i.GetNumber(),
	}).Tracef("deleting init %d in the database", i.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitFromLibrary
	init := database.InitFromLibrary(i)

	// send query to the database
	return e.client.
		Table(constants.TableInit).
		Delete(init).
		Error
}
