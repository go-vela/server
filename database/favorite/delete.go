// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// DeleteFavorite deletes a user favorite in the database.
func (e *Engine) DeleteFavorite(ctx context.Context, u *api.User, r *api.Repo) error {
	e.logger.WithFields(logrus.Fields{
		"repo": r.GetFullName(),
	}).Tracef("deleting favorite for user %s", u.GetName())

	return e.client.Transaction(func(tx *gorm.DB) error {
		currentPos := sql.NullInt64{}

		err := tx.Table(constants.TableFavorite).
			Select("position").
			Where("user_id = ? AND repo_id = ?", u.GetID(), r.GetID()).
			Scan(&currentPos).Error
		if err != nil {
			return fmt.Errorf("error getting current favorite position: %w", err)
		}

		if !currentPos.Valid {
			return fmt.Errorf("favorite not found for repo: %s", r.GetFullName())
		}

		err = tx.Exec(`
			DELETE FROM favorites
			 WHERE user_id = ? AND repo_id = ?`,
			u.GetID(),
			r.GetID(),
		).Error
		if err != nil {
			return fmt.Errorf("error deleting favorite: %w", err)
		}

		// shift favorites up to fill gap
		err = tx.Exec(`
				UPDATE favorites
				SET position = position - 1
				WHERE user_id = ?
				  AND position > ?
			`, u.GetID(), currentPos.Int64).Error
		if err != nil {
			return fmt.Errorf("error shifting favorites: %w", err)
		}

		return nil
	})
}
