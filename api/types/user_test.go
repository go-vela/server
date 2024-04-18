// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestTypes_User_Sanitize(t *testing.T) {
	// setup types
	u := testUser()

	want := new(User)
	want.SetID(1)
	want.SetName("octocat")
	want.SetActive(true)
	want.SetToken(constants.SecretMask)
	want.SetRefreshToken(constants.SecretMask)

	// run test
	got := u.Sanitize()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sanitize is %v, want %v", got, want)
	}
}

func TestTypes_User_Environment(t *testing.T) {
	// setup types
	want := map[string]string{
		"VELA_USER_ACTIVE":    "true",
		"VELA_USER_ADMIN":     "false",
		"VELA_USER_FAVORITES": "[\"github/octocat\"]",
		"VELA_USER_NAME":      "octocat",
	}

	// run test
	got := testUser().Environment()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Environment is %v, want %v", got, want)
	}
}

func TestTypes_User_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		user *User
		want *User
	}{
		{
			user: testUser(),
			want: testUser(),
		},
		{
			user: new(User),
			want: new(User),
		},
	}

	// run tests
	for _, test := range tests {
		if test.user.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.user.GetID(), test.want.GetID())
		}

		if test.user.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.user.GetName(), test.want.GetName())
		}

		if test.user.GetRefreshToken() != test.want.GetRefreshToken() {
			t.Errorf("GetRefreshToken is %v, want %v", test.user.GetRefreshToken(), test.want.GetRefreshToken())
		}

		if test.user.GetToken() != test.want.GetToken() {
			t.Errorf("GetToken is %v, want %v", test.user.GetToken(), test.want.GetToken())
		}

		if !reflect.DeepEqual(test.user.GetFavorites(), test.want.GetFavorites()) {
			t.Errorf("GetFavorites is %v, want %v", test.user.GetFavorites(), test.want.GetFavorites())
		}

		if test.user.GetActive() != test.want.GetActive() {
			t.Errorf("GetActive is %v, want %v", test.user.GetActive(), test.want.GetActive())
		}

		if test.user.GetAdmin() != test.want.GetAdmin() {
			t.Errorf("GetAdmin is %v, want %v", test.user.GetAdmin(), test.want.GetAdmin())
		}

		if !reflect.DeepEqual(test.user.GetDashboards(), test.want.GetDashboards()) {
			t.Errorf("GetDashboards is %v, want %v", test.user.GetDashboards(), test.want.GetDashboards())
		}
	}
}

func TestTypes_User_Setters(t *testing.T) {
	// setup types
	var u *User

	// setup tests
	tests := []struct {
		user *User
		want *User
	}{
		{
			user: testUser(),
			want: testUser(),
		},
		{
			user: u,
			want: new(User),
		},
	}

	// run tests
	for _, test := range tests {
		test.user.SetID(test.want.GetID())
		test.user.SetName(test.want.GetName())
		test.user.SetRefreshToken(test.want.GetRefreshToken())
		test.user.SetToken(test.want.GetToken())
		test.user.SetFavorites(test.want.GetFavorites())
		test.user.SetActive(test.want.GetActive())
		test.user.SetAdmin(test.want.GetAdmin())
		test.user.SetDashboards(test.want.GetDashboards())

		if test.user.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.user.GetID(), test.want.GetID())
		}

		if test.user.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.user.GetName(), test.want.GetName())
		}

		if test.user.GetRefreshToken() != test.want.GetRefreshToken() {
			t.Errorf("SetRefreshToken is %v, want %v", test.user.GetRefreshToken(), test.want.GetRefreshToken())
		}

		if test.user.GetToken() != test.want.GetToken() {
			t.Errorf("SetToken is %v, want %v", test.user.GetToken(), test.want.GetToken())
		}

		if !reflect.DeepEqual(test.user.GetFavorites(), test.want.GetFavorites()) {
			t.Errorf("SetFavorites is %v, want %v", test.user.GetFavorites(), test.want.GetFavorites())
		}

		if test.user.GetActive() != test.want.GetActive() {
			t.Errorf("SetActive is %v, want %v", test.user.GetActive(), test.want.GetActive())
		}

		if test.user.GetAdmin() != test.want.GetAdmin() {
			t.Errorf("SetAdmin is %v, want %v", test.user.GetAdmin(), test.want.GetAdmin())
		}

		if !reflect.DeepEqual(test.user.GetDashboards(), test.want.GetDashboards()) {
			t.Errorf("SetDashboards is %v, want %v", test.user.GetDashboards(), test.want.GetDashboards())
		}
	}
}

func TestTypes_User_String(t *testing.T) {
	// setup types
	u := testUser()

	want := fmt.Sprintf(`{
  Active: %t,
  Admin: %t,
  Dashboards: %s,
  Favorites: %s,
  ID: %d,
  Name: %s,
  Token: %s,
}`,
		u.GetActive(),
		u.GetAdmin(),
		u.GetDashboards(),
		u.GetFavorites(),
		u.GetID(),
		u.GetName(),
		u.GetToken(),
	)

	// run test
	got := u.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *User {
	u := new(User)

	u.SetID(1)
	u.SetName("octocat")
	u.SetToken("superSecretToken")
	u.SetFavorites([]string{"github/octocat"})
	u.SetActive(true)
	u.SetAdmin(false)
	u.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7", "ba657dab-bc6e-421f-9188-86272bd0069a"})

	return u
}
