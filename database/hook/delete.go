// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// DeleteHook deletes an existing hook from the database.
func (e *engine) DeleteHook(h *library.Hook) error {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("deleting hook %d in the database", h.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#HookFromLibrary
	hook := database.HookFromLibrary(h)

	// send query to the database
	return e.client.
		Table(constants.TableHook).
		Delete(hook).
		Error
}
