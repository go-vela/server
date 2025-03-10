// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	api "github.com/go-vela/server/api/types"
)

var (
	currentTime = time.Now()
	tsCreate    = currentTime.UTC().Unix()
	tsUpdate    = currentTime.Add(time.Hour * 1).UTC().Unix()
)

func TestDatabase_Secret_Decrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"
	encrypted := testSecret()

	err := encrypted.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt secret: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		key     string
		secret  Secret
	}{
		{
			failure: false,
			key:     key,
			secret:  *encrypted,
		},
		{
			failure: true,
			key:     "",
			secret:  *encrypted,
		},
		{
			failure: true,
			key:     key,
			secret:  *testSecret(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.secret.Decrypt(test.key)

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

func TestDatabase_Secret_Encrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"

	// setup tests
	tests := []struct {
		failure bool
		key     string
		secret  *Secret
	}{
		{
			failure: false,
			key:     key,
			secret:  testSecret(),
		},
		{
			failure: true,
			key:     "",
			secret:  testSecret(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.secret.Encrypt(test.key)

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

func TestDatabase_Secret_Nullify(t *testing.T) {
	// setup types
	var s *Secret

	want := &Secret{
		ID:          sql.NullInt64{Int64: 0, Valid: false},
		Org:         sql.NullString{String: "", Valid: false},
		Repo:        sql.NullString{String: "", Valid: false},
		Team:        sql.NullString{String: "", Valid: false},
		Name:        sql.NullString{String: "", Valid: false},
		Value:       sql.NullString{String: "", Valid: false},
		Type:        sql.NullString{String: "", Valid: false},
		AllowEvents: sql.NullInt64{Int64: 0, Valid: false},
		CreatedAt:   sql.NullInt64{Int64: 0, Valid: false},
		CreatedBy:   sql.NullString{String: "", Valid: false},
		UpdatedAt:   sql.NullInt64{Int64: 0, Valid: false},
		UpdatedBy:   sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		secret *Secret
		want   *Secret
	}{
		{
			secret: testSecret(),
			want:   testSecret(),
		},
		{
			secret: s,
			want:   nil,
		},
		{
			secret: new(Secret),
			want:   want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.secret.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_Secret_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Secret)

	want.SetID(1)
	want.SetOrg("github")
	want.SetOrgSCMID(1)
	want.SetRepo("octocat")
	want.SetRepoSCMID(1)
	want.SetTeam("octokitties")
	want.SetTeamSCMID(1)
	want.SetName("foo")
	want.SetValue("bar")
	want.SetType("repo")
	want.SetImages([]string{"alpine"})
	want.SetAllowEvents(api.NewEventsFromMask(1))
	want.SetAllowCommand(true)
	want.SetAllowSubstitution(true)
	want.SetCreatedAt(tsCreate)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(tsUpdate)
	want.SetUpdatedBy("octocat2")

	// run test
	got := testSecret().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_Secret_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		failure bool
		secret  *Secret
	}{
		{
			failure: false,
			secret:  testSecret(),
		},
		{ // no name set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				Value:     sql.NullString{String: "bar", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no org set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Value:     sql.NullString{String: "bar", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no repo set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Value:     sql.NullString{String: "bar", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no team set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				TeamSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Value:     sql.NullString{String: "bar", Valid: true},
				Type:      sql.NullString{String: "shared", Valid: true},
			},
		},
		{ // no type set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Value:     sql.NullString{String: "bar", Valid: true},
			},
		},
		{ // no value set for secret
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no value set for org scm id
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				RepoSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				TeamSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no value set for repo scm id
			failure: true,
			secret: &Secret{
				ID:        sql.NullInt64{Int64: 1, Valid: true},
				Org:       sql.NullString{String: "github", Valid: true},
				OrgSCMID:  sql.NullInt64{Int64: 1, Valid: true},
				Repo:      sql.NullString{String: "octocat", Valid: true},
				Team:      sql.NullString{String: "octokitties", Valid: true},
				TeamSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Name:      sql.NullString{String: "foo", Valid: true},
				Type:      sql.NullString{String: "repo", Valid: true},
			},
		},
		{ // no value set for team scm id on shared secret
			failure: true,
			secret: &Secret{
				ID:       sql.NullInt64{Int64: 1, Valid: true},
				Org:      sql.NullString{String: "github", Valid: true},
				OrgSCMID: sql.NullInt64{Int64: 1, Valid: true},
				Repo:     sql.NullString{String: "octocat", Valid: true},
				Team:     sql.NullString{String: "octokitties", Valid: true},
				Name:     sql.NullString{String: "foo", Valid: true},
				Type:     sql.NullString{String: "shared", Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.secret.Validate()

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

func TestDatabase_SecretFromAPI(t *testing.T) {
	// setup types
	s := new(api.Secret)

	s.SetID(1)
	s.SetOrg("github")
	s.SetOrgSCMID(1)
	s.SetRepo("octocat")
	s.SetRepoSCMID(1)
	s.SetTeam("octokitties")
	s.SetTeamSCMID(1)
	s.SetName("foo")
	s.SetValue("bar")
	s.SetType("repo")
	s.SetImages([]string{"alpine"})
	s.SetAllowEvents(api.NewEventsFromMask(1))
	s.SetAllowCommand(true)
	s.SetAllowSubstitution(true)
	s.SetCreatedAt(tsCreate)
	s.SetCreatedBy("octocat")
	s.SetUpdatedAt(tsUpdate)
	s.SetUpdatedBy("octocat2")

	want := testSecret()

	// run test
	got := SecretFromAPI(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SecretFromAPI is %v, want %v", got, want)
	}
}

// testSecret is a test helper function to create a Secret
// type with all fields set to a fake value.
func testSecret() *Secret {
	return &Secret{
		ID:                sql.NullInt64{Int64: 1, Valid: true},
		Org:               sql.NullString{String: "github", Valid: true},
		OrgSCMID:          sql.NullInt64{Int64: 1, Valid: true},
		Repo:              sql.NullString{String: "octocat", Valid: true},
		RepoSCMID:         sql.NullInt64{Int64: 1, Valid: true},
		Team:              sql.NullString{String: "octokitties", Valid: true},
		TeamSCMID:         sql.NullInt64{Int64: 1, Valid: true},
		Name:              sql.NullString{String: "foo", Valid: true},
		Value:             sql.NullString{String: "bar", Valid: true},
		Type:              sql.NullString{String: "repo", Valid: true},
		Images:            []string{"alpine"},
		AllowEvents:       sql.NullInt64{Int64: 1, Valid: true},
		AllowCommand:      sql.NullBool{Bool: true, Valid: true},
		AllowSubstitution: sql.NullBool{Bool: true, Valid: true},
		CreatedAt:         sql.NullInt64{Int64: tsCreate, Valid: true},
		CreatedBy:         sql.NullString{String: "octocat", Valid: true},
		UpdatedAt:         sql.NullInt64{Int64: tsUpdate, Valid: true},
		UpdatedBy:         sql.NullString{String: "octocat2", Valid: true},
	}
}
