// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteHook deletes an existing hook from the database.
func (e *engine) DeleteHook(ctx context.Context, h *api.Hook) error {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("deleting hook %d", h.GetNumber())

	hook := types.HookFromAPI(h)

	// send query to the database
	return e.client.
		Table(constants.TableHook).
		Delete(hook).
		Error
}
