// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// CreateFavorite creates a new favorite in the database.
func (e *Engine) CreateFavorite(ctx context.Context, u *api.User, f *api.Favorite) error {
	e.logger.WithFields(logrus.Fields{
		"repo": f.GetRepo(),
	}).Tracef("creating favorite for user %s", u.GetName())

	return e.client.
		WithContext(ctx).
		Exec(
			`INSERT INTO favorites (user_id, repo_id, position)
			 SELECT ?, id, ? FROM repos WHERE full_name = ?;`,
			u.GetID(),
			f.GetPosition(),
			f.GetRepo(),
		).Error
}
