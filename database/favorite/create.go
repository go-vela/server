// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// CreateFavorite creates a new favorite in the database.
func (e *Engine) CreateFavorite(ctx context.Context, u *api.User, f *api.Favorite) error {
	e.logger.WithFields(logrus.Fields{
		"repo": f.GetRepo(),
	}).Tracef("creating favorite for user %s", u.GetName())

	res := e.client.
		WithContext(ctx).
		Exec(`
			INSERT INTO favorites (user_id, repo_id, position)
			SELECT ?, r.id,
			       (SELECT COALESCE(MAX(position), 0) + 1 FROM favorites WHERE user_id = ?)
			FROM repos r
			WHERE r.full_name = ?;
		`, u.GetID(), u.GetID(), f.GetRepo())

	e.logger.Infof("result err: %s, rows: %d", res.Error, res.RowsAffected)

	if res.Error != nil {
		return res.Error
	}

	// no repo found
	if res.RowsAffected == 0 {
		return fmt.Errorf("repo not found: %s", f.GetRepo())
	}

	return nil
}
