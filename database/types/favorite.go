// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"

	api "github.com/go-vela/server/api/types"
)

type (
	// Favorite represents a user's favorite repository.
	Favorite struct {
		Position sql.NullInt64  `sql:"position"`
		RepoName sql.NullString `sql:"repo_name"`
		UserID   sql.NullInt64  `sql:"user_id"`
		RepoID   sql.NullInt64  `sql:"repo_id"`
	}
)

// ToAPI converts the Favorites type
// to a slice of repos.
func (f *Favorite) ToAPI() *api.Favorite {
	favorite := new(api.Favorite)

	if f.Position.Valid {
		favorite.SetPosition(f.Position.Int64)
	}

	favorite.SetRepo(f.RepoName.String)

	return favorite
}

func FavoriteFromAPI(f *api.Favorite) *Favorite {
	return &Favorite{
		Position: sql.NullInt64{Int64: f.GetPosition(), Valid: true},
		RepoName: sql.NullString{String: f.GetRepo(), Valid: true},
	}
}
