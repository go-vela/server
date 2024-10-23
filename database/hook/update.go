// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateHook updates an existing hook in the database.
func (e *engine) UpdateHook(ctx context.Context, h *api.Hook) (*api.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("updating hook %d", h.GetNumber())

	hook := types.HookFromAPI(h)

	err := hook.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.WithContext(ctx).Table(constants.TableHook).Save(hook).Error
	if err != nil {
		return nil, err
	}

	result := hook.ToAPI()
	result.SetRepo(h.GetRepo())
	result.SetBuild(h.GetBuild())

	return result, nil
}
