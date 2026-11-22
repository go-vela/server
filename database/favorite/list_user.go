// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListUserFavorites gets a list of all user favorites from the database.
func (e *Engine) ListUserFavorites(ctx context.Context, u *api.User) ([]*api.Favorite, error) {
	e.logger.Trace("listing all user favorites")

	result := []types.Favorite{}

	err := e.client.
		WithContext(ctx).
		Table(constants.TableFavorite+" f").
		Select("r.full_name as repo_name, f.position").
		Joins("JOIN "+constants.TableRepo+" r ON f.repo_id = r.id").
		Where("f.user_id = ?", u.GetID()).
		Order("f.position ASC").
		Find(&result).
		Error
	if err != nil {
		return nil, err
	}

	favorites := make([]*api.Favorite, 0, len(result))
	for _, res := range result {
		favorites = append(favorites, res.ToAPI())
	}

	return favorites, nil
}
