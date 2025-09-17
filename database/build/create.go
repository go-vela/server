// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateBuild creates a new build in the database.
func (e *Engine) CreateBuild(ctx context.Context, b *api.Build) (*api.Build, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("creating build %d", b.GetNumber())

	var res *api.Build

	err := e.client.Transaction(func(tx *gorm.DB) error {
		r := b.GetRepo()

		if r.GetID() == 0 {
			return fmt.Errorf("repo ID must be set")
		}

		var next int64

		err := tx.Raw("UPDATE repos SET counter = counter + 1 WHERE id = ? RETURNING counter", b.GetRepo().GetID()).Scan(&next).Error
		if err != nil {
			return err
		}

		b.SetNumber(next)
		r.SetCounter(next)

		build := types.BuildFromAPI(b)

		err = build.Validate()
		if err != nil {
			return err
		}

		// crop build if any columns are too large
		build = build.Crop()

		err = tx.Table(constants.TableBuild).Create(build).Error
		if err != nil {
			return err
		}

		res = build.ToAPI()
		res.SetRepo(r)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
