// SPDX-License-Identifier: Apache-2.0

package favorite

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// FavoriteInterface represents the Vela interface for user favorite
// functions with the supported Database backends.
type FavoriteInterface interface {
	// Favorite Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language
	CreateFavoritesTable(context.Context, string) error

	CreateFavoritesIndexes(context.Context) error

	// Favorite Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	CreateFavorite(context.Context, *api.User, *api.Favorite) error

	DeleteFavorite(context.Context, *api.User, *api.Favorite) error

	ListUserFavorites(context.Context, *api.User) ([]*api.Favorite, error)

	UpdateFavorites(context.Context, *api.User, []*api.Favorite) error
}
