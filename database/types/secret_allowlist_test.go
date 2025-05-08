// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestDatabase_SecretAllowlist_Nullify(t *testing.T) {
	// setup types
	var s *SecretAllowlist

	want := &SecretAllowlist{
		ID:       sql.NullInt64{Int64: 0, Valid: false},
		SecretID: sql.NullInt64{Int64: 0, Valid: false},
		Repo:     sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		record *SecretAllowlist
		want   *SecretAllowlist
	}{
		{
			record: testSecretAllowlist(),
			want:   testSecretAllowlist(),
		},
		{
			record: s,
			want:   nil,
		},
		{
			record: new(SecretAllowlist),
			want:   want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.record.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_SecretAllowlist_Validate(t *testing.T) {
	// setup types
	tests := []struct {
		failure bool
		record  *SecretAllowlist
	}{
		{
			failure: false,
			record:  testSecretAllowlist(),
		},
		{ // no secret_id set
			failure: true,
			record: &SecretAllowlist{
				ID:   sql.NullInt64{Int64: 1, Valid: true},
				Repo: sql.NullString{String: "github/octocat", Valid: true},
			},
		},
		{ // no repo set
			failure: true,
			record: &SecretAllowlist{
				ID:       sql.NullInt64{Int64: 1, Valid: true},
				SecretID: sql.NullInt64{Int64: 1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.record.Validate()

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

func TestDatabase_SecretAllowlistFromAPI(t *testing.T) {
	// setup types
	s := new(api.Secret)

	s.SetID(1)
	s.SetOrg("github")
	s.SetRepo("octocat")
	s.SetTeam("octokitties")
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

	want := &SecretAllowlist{
		SecretID: sql.NullInt64{Int64: 1, Valid: true},
		Repo:     sql.NullString{String: "github/octocat", Valid: true},
	}

	// run test
	got := SecretAllowlistFromAPI(s, "github/octocat")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SecretAllowlistFromAPI is %v, want %v", got, want)
	}
}

// testSecretAllowlist is a test helper function to create a SecretAllowlist
// type with all fields set to a fake value.
func testSecretAllowlist() *SecretAllowlist {
	return &SecretAllowlist{
		ID:       sql.NullInt64{Int64: 1, Valid: true},
		SecretID: sql.NullInt64{Int64: 1, Valid: true},
		Repo:     sql.NullString{String: "github/octocat", Valid: true},
	}
}
