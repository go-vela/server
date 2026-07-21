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

// UpdateFavoritePosition updates a single favorite position and applies the shift to other favorites.
func (e *Engine) UpdateFavoritePosition(ctx context.Context, u *api.User, r *api.Repo, f *api.Favorite) error {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("updating user %s favorites", u.GetName())

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

		if currentPos.Int64 == f.GetPosition() {
			return nil
		}

		var count int64
		if err := tx.Table(constants.TableFavorite).
			Where("user_id = ?", u.GetID()).
			Count(&count).Error; err != nil {
			return fmt.Errorf("error counting favorites: %w", err)
		}

		// clamp position
		if f.GetPosition() <= 0 {
			f.SetPosition(1)
		}

		if f.GetPosition() > count {
			f.SetPosition(count)
		}

		if currentPos.Int64 > f.GetPosition() {
			// moving up - smaller position
			err := tx.Exec(`
				UPDATE favorites
				SET position = position + 1
				WHERE user_id = ?
				  AND position >= ?
				  AND position < ?
			`, u.GetID(), f.GetPosition(), currentPos.Int64).Error
			if err != nil {
				return fmt.Errorf("error shifting favorites: %w", err)
			}
		} else {
			// moving down - larger position
			err := tx.Exec(`
				UPDATE favorites
				SET position = position - 1
				WHERE user_id = ?
				  AND position <= ?
				  AND position > ?
			`, u.GetID(), f.GetPosition(), currentPos.Int64).Error
			if err != nil {
				return fmt.Errorf("error shifting favorites: %w", err)
			}
		}

		err = tx.Table(constants.TableFavorite).
			Where("user_id = ? AND repo_id = ?", u.GetID(), r.GetID()).
			Update("position", f.GetPosition()).Error
		if err != nil {
			return fmt.Errorf("error updating favorite position: %w", err)
		}

		return nil
	})
}
