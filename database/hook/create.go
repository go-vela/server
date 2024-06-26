// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// CreateHook creates a new hook in the database.
func (e *engine) CreateHook(ctx context.Context, h *library.Hook) (*library.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("creating hook %d", h.GetNumber())

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

	result := e.client.Table(constants.TableHook).Create(hook)

	// send query to the database
	return hook.ToLibrary(), result.Error
}
