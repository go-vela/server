// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// DeleteHook deletes an existing hook from the database.
func (e *engine) DeleteHook(ctx context.Context, h *library.Hook) error {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("deleting hook %d", h.GetNumber())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#HookFromLibrary
	hook := database.HookFromLibrary(h)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Delete(hook).
		Error
}
