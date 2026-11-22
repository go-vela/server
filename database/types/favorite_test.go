// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_Favorite_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Favorite)
	want.SetPosition(1)
	want.SetRepo("octocat/Hello-World")

	// run test
	got := testFavorite().ToAPI()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToAPI() mismatch (-want +got):\n%s", diff)
	}
}

func TestTypes_FavoriteFromAPI(t *testing.T) {
	// setup types
	want := &Favorite{
		Position: sql.NullInt64{Int64: 1, Valid: true},
		RepoName: sql.NullString{String: "octocat/Hello-World", Valid: true},
	}

	f := new(api.Favorite)
	f.SetPosition(1)
	f.SetRepo("octocat/Hello-World")

	// run test
	got := FavoriteFromAPI(f)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FavoriteFromAPI is %v, want %v", got, want)
	}
}

// testFavorite is a test helper function to create a Favorite
// type with all fields set to a fake value.
func testFavorite() *Favorite {
	return &Favorite{
		Position: sql.NullInt64{Int64: 1, Valid: true},
		RepoName: sql.NullString{String: "octocat/Hello-World", Valid: true},
		UserID:   sql.NullInt64{Int64: 1, Valid: true},
		RepoID:   sql.NullInt64{Int64: 1, Valid: true},
	}
}
