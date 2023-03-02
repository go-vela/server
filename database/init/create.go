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

// CreateInit creates a new init in the database.
func (e *engine) CreateInit(i *library.Init) error {
	e.logger.WithFields(logrus.Fields{
		"init": i.GetNumber(),
	}).Tracef("creating init %d in the database", i.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#InitFromLibrary
	init := database.InitFromLibrary(i)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Init.Validate
	err := init.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		Table(constants.TableInit).
		Create(init).
		Error
}
