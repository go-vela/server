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

// UpdateHook updates an existing hook in the database.
func (e *engine) UpdateHook(h *library.Hook) (*library.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("updating hook %d in the database", h.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#HookFromLibrary
	hook := database.HookFromLibrary(h)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Hook.Validate
	err := hook.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.Table(constants.TableHook).Save(hook)

	// send query to the database
	return hook.ToLibrary(), result.Error
}
