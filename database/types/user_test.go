// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"strconv"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_User_Decrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"
	encrypted := testUser()

	err := encrypted.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt user: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		key     string
		user    User
	}{
		{
			failure: false,
			key:     key,
			user:    *encrypted,
		},
		{
			failure: true,
			key:     "",
			user:    *encrypted,
		},
		{
			failure: true,
			key:     key,
			user:    *testUser(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.user.Decrypt(test.key)

		if test.failure {
			if err == nil {
				t.Errorf("Decrypt should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Decrypt returned err: %v", err)
		}
	}
}

func TestTypes_User_Encrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"

	// setup tests
	tests := []struct {
		failure bool
		key     string
		user    *User
	}{
		{
			failure: false,
			key:     key,
			user:    testUser(),
		},
		{
			failure: true,
			key:     "",
			user:    testUser(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.user.Encrypt(test.key)

		if test.failure {
			if err == nil {
				t.Errorf("Encrypt should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Encrypt returned err: %v", err)
		}
	}
}

func TestTypes_User_Nullify(t *testing.T) {
	// setup types
	var u *User

	want := &User{
		ID:           sql.NullInt64{Int64: 0, Valid: false},
		Name:         sql.NullString{String: "", Valid: false},
		RefreshToken: sql.NullString{String: "", Valid: false},
		Token:        sql.NullString{String: "", Valid: false},
		Active:       sql.NullBool{Bool: false, Valid: false},
		Admin:        sql.NullBool{Bool: false, Valid: false},
	}

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
			want: nil,
		},
		{
			user: new(User),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.user.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_User_ToAPI(t *testing.T) {
	// setup types
	want := new(api.User)

	want.SetID(1)
	want.SetName("octocat")
	want.SetRefreshToken("superSecretRefreshToken")
	want.SetToken("superSecretToken")
	want.SetFavorites([]string{"github/octocat"})
	want.SetActive(true)
	want.SetAdmin(false)
	want.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"})

	// run test
	got := testUser().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_User_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		user    *User
	}{
		{
			failure: false,
			user:    testUser(),
		},
		{ // no name set for user
			failure: true,
			user: &User{
				ID:    sql.NullInt64{Int64: 1, Valid: true},
				Token: sql.NullString{String: "superSecretToken", Valid: true},
			},
		},
		{ // no token set for user
			failure: true,
			user: &User{
				ID:   sql.NullInt64{Int64: 1, Valid: true},
				Name: sql.NullString{String: "octocat", Valid: true},
			},
		},
		{ // invalid name set for user
			failure: true,
			user: &User{
				ID:           sql.NullInt64{Int64: 1, Valid: true},
				Name:         sql.NullString{String: "!@#$%^&*()", Valid: true},
				RefreshToken: sql.NullString{String: "superSecretRefreshToken", Valid: true},
				Token:        sql.NullString{String: "superSecretToken", Valid: true},
			},
		},
		{ // invalid favorites set for user
			failure: true,
			user: &User{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Name:      sql.NullString{String: "octocat", Valid: true},
				Token:     sql.NullString{String: "superSecretToken", Valid: true},
				Favorites: exceededField(500),
			},
		},
		{ // invalid dashboards set for user
			failure: true,
			user: &User{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				Token:      sql.NullString{String: "superSecretToken", Valid: true},
				Dashboards: exceededField(11),
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.user.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestFromAPI(t *testing.T) {
	// setup types
	u := new(api.User)

	u.SetID(1)
	u.SetName("octocat")
	u.SetRefreshToken("superSecretRefreshToken")
	u.SetToken("superSecretToken")
	u.SetFavorites([]string{"github/octocat"})
	u.SetActive(true)
	u.SetAdmin(false)
	u.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"})

	want := testUser()

	// run test
	got := UserFromAPI(u)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromAPI is %v, want %v", got, want)
	}
}

// testUser is a test helper function to create a User
// type with all fields set to a fake value.
func testUser() *User {
	return &User{
		ID:           sql.NullInt64{Int64: 1, Valid: true},
		Name:         sql.NullString{String: "octocat", Valid: true},
		RefreshToken: sql.NullString{String: "superSecretRefreshToken", Valid: true},
		Token:        sql.NullString{String: "superSecretToken", Valid: true},
		Favorites:    []string{"github/octocat"},
		Active:       sql.NullBool{Bool: true, Valid: true},
		Admin:        sql.NullBool{Bool: false, Valid: true},
		Dashboards:   []string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"},
	}
}

// exceededField returns a list of strings that exceed the maximum size of a field.
func exceededField(indexes int) []string {
	// initialize empty favorites
	values := []string{}

	// add enough strings to exceed the character limit
	for i := 0; i < indexes; i++ {
		// construct field
		// use i to adhere to unique favorites
		field := "github/octocat-" + strconv.Itoa(i)

		values = append(values, field)
	}

	return values
}
