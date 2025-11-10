// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// DeleteFavorite deletes a user favorite in the database.
func (e *Engine) DeleteFavorite(ctx context.Context, u *api.User, f *api.Favorite) error {
	e.logger.WithFields(logrus.Fields{
		"repo": f.GetRepo(),
	}).Tracef("deleting favorite for user %s", u.GetName())

	return e.client.
		WithContext(ctx).
		Exec(
			`DELETE FROM favorites
			 WHERE user_id = ? AND repo_id = (SELECT id FROM repos WHERE full_name = ?)`,
			u.GetID(),
			f.GetRepo(),
		).Error
}
