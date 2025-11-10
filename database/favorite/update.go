// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// UpdateFavorites updates a user's favorites in the database.
func (e *Engine) UpdateFavorites(ctx context.Context, u *api.User, favs []*api.Favorite) error {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("updating user %s favorites", u.GetName())

	if len(favs) == 0 {
		return nil
	}

	switch e.config.Driver {
	case constants.DriverPostgres:
		valueStrings := make([]string, 0, len(favs))
		valueArgs := make([]any, 0, len(favs)*2)

		for _, fav := range favs {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, fav.GetRepo(), fav.GetPosition())
		}

		args := []any{u.GetID()}
		args = append(args, valueArgs...)

		query := fmt.Sprintf(`
WITH input(repo_name, position) AS (
  VALUES %s
)
INSERT INTO favorites (user_id, repo_id, position)
SELECT ?, r.id, input.position
FROM input
JOIN repos r ON r.full_name = input.repo_name
ON CONFLICT (user_id, repo_id) DO UPDATE
  SET position = excluded.position;
`, strings.Join(valueStrings, ", "))

		return e.client.
			WithContext(ctx).
			Exec(query, args...).
			Error

	default:
		return fmt.Errorf("unsupported database driver: %s", e.config.Driver)
	}
}
