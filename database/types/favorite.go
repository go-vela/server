// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"fmt"

	api "github.com/go-vela/server/api/types"
)

type (
	// Favorite represents a user's favorite repository.
	Favorite struct {
		Position sql.NullInt64  `sql:"position"`
		RepoName sql.NullString `sql:"repo_name"`
	}
)

// ToAPI converts the Favorites type
// to a slice of repos.
func (f *Favorite) ToAPI() *api.Favorite {
	favorite := new(api.Favorite)

	favorite.SetPosition(f.Position.Int64)
	favorite.SetRepo(f.RepoName.String)

	return favorite
}

// Validate verifies the necessary fields for
// the Favorite type are populated correctly.
func (f *Favorite) Validate() error {
	// verify the Repo field is populated
	if len(f.RepoName.String) == 0 {
		return fmt.Errorf("empty favorite repo provided")
	}

	return nil
}

func FavoriteFromAPI(f *api.Favorite) *Favorite {
	return &Favorite{
		Position: sql.NullInt64{Int64: f.GetPosition(), Valid: true},
		RepoName: sql.NullString{String: f.GetRepo(), Valid: true},
	}
}
