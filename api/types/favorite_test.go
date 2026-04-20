// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Favorite_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		favorite *Favorite
		want     *Favorite
	}{
		{
			favorite: testFavorite(),
			want:     testFavorite(),
		},
		{
			favorite: new(Favorite),
			want:     new(Favorite),
		},
	}

	// run tests
	for _, test := range tests {
		if test.favorite.GetPosition() != test.want.GetPosition() {
			t.Errorf("GetPosition is %v, want %v", test.favorite.GetPosition(), test.want.GetPosition())
		}

		if test.favorite.GetRepo() != test.want.GetRepo() {
			t.Errorf("GetRepo is %v, want %v", test.favorite.GetRepo(), test.want.GetRepo())
		}
	}
}

func TestTypes_Favorite_Setters(t *testing.T) {
	// setup types
	var f *Favorite

	// setup tests
	tests := []struct {
		favorite *Favorite
		want     *Favorite
	}{
		{
			favorite: testFavorite(),
			want:     testFavorite(),
		},
		{
			favorite: f,
			want:     new(Favorite),
		},
	}

	// run tests
	for _, test := range tests {
		test.favorite.SetPosition(test.want.GetPosition())
		test.favorite.SetRepo(test.want.GetRepo())

		if test.favorite.GetPosition() != test.want.GetPosition() {
			t.Errorf("SetPosition is %v, want %v", test.favorite.GetPosition(), test.want.GetPosition())
		}

		if test.favorite.GetRepo() != test.want.GetRepo() {
			t.Errorf("SetRepo is %v, want %v", test.favorite.GetRepo(), test.want.GetRepo())
		}
	}
}

func TestTypes_Favorite_String(t *testing.T) {
	// setup types
	f := testFavorite()

	want := fmt.Sprintf(`{
  Position: %d,
  Repo: %s,
}`,
		f.GetPosition(),
		f.GetRepo(),
	)

	// run test
	got := f.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testFavorite is a test helper function to create a Favorite
// type with all fields set to a fake value.
func testFavorite() *Favorite {
	f := new(Favorite)

	f.SetPosition(1)
	f.SetRepo("octocat/Hello-World")

	return f
}
