// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateHook creates a new hook in the database.
func (e *Engine) CreateHook(ctx context.Context, h *api.Hook) (*api.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("creating hook %d", h.GetNumber())

	var res *api.Hook

	err := e.client.Transaction(func(tx *gorm.DB) error {
		r := h.GetRepo()

		if r.GetID() == 0 {
			return fmt.Errorf("repo ID must be set")
		}

		var next int64

		err := tx.Raw("UPDATE repos SET hook_counter = hook_counter + 1 WHERE id = ? RETURNING hook_counter", h.GetRepo().GetID()).Scan(&next).Error
		if err != nil {
			return err
		}

		h.SetNumber(next)
		r.SetHookCounter(next)

		hook := types.HookFromAPI(h)

		err = hook.Validate()
		if err != nil {
			return err
		}

		err = tx.Table(constants.TableHook).Create(hook).Error
		if err != nil {
			return err
		}

		res = hook.ToAPI()
		res.SetRepo(r)
		res.SetBuild(h.GetBuild())

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
