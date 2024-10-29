// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestTypes_Repo_Decrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"
	encrypted := testRepo()

	err := encrypted.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt repo: %v", err)
	}

	err = encrypted.Owner.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt user: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		key     string
		repo    Repo
	}{
		{
			failure: false,
			key:     key,
			repo:    *encrypted,
		},
		{
			failure: true,
			key:     "",
			repo:    *encrypted,
		},
		{
			failure: true,
			key:     key,
			repo:    *testRepo(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.repo.Decrypt(test.key)

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

func TestTypes_Repo_Encrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"

	// setup tests
	tests := []struct {
		failure bool
		key     string
		repo    *Repo
	}{
		{
			failure: false,
			key:     key,
			repo:    testRepo(),
		},
		{
			failure: true,
			key:     "",
			repo:    testRepo(),
		},
	}

	// run tests
	for _, test := range tests {
		err := test.repo.Encrypt(test.key)

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

func TestTypes_Repo_Nullify(t *testing.T) {
	// setup types
	var r *Repo

	want := &Repo{
		ID:           sql.NullInt64{Int64: 0, Valid: false},
		UserID:       sql.NullInt64{Int64: 0, Valid: false},
		Hash:         sql.NullString{String: "", Valid: false},
		Org:          sql.NullString{String: "", Valid: false},
		Name:         sql.NullString{String: "", Valid: false},
		FullName:     sql.NullString{String: "", Valid: false},
		Link:         sql.NullString{String: "", Valid: false},
		Clone:        sql.NullString{String: "", Valid: false},
		Branch:       sql.NullString{String: "", Valid: false},
		Timeout:      sql.NullInt64{Int64: 0, Valid: false},
		AllowEvents:  sql.NullInt64{Int64: 0, Valid: false},
		Visibility:   sql.NullString{String: "", Valid: false},
		PipelineType: sql.NullString{String: "", Valid: false},
		ApproveBuild: sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		repo *Repo
		want *Repo
	}{
		{
			repo: testRepo(),
			want: testRepo(),
		},
		{
			repo: r,
			want: nil,
		},
		{
			repo: new(Repo),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.repo.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_Repo_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Repo)
	e := api.NewEventsFromMask(1)

	owner := testutils.APIUser().Crop()
	owner.SetID(1)
	owner.SetName("octocat")
	owner.SetActive(true)
	owner.SetToken("superSecretToken")
	owner.SetRefreshToken("superSecretRefreshToken")

	want.SetID(1)
	want.SetOwner(owner)
	want.SetHash("superSecretHash")
	want.SetOrg("github")
	want.SetName("octocat")
	want.SetFullName("github/octocat")
	want.SetLink("https://github.com/github/octocat")
	want.SetClone("https://github.com/github/octocat.git")
	want.SetBranch("main")
	want.SetTopics([]string{"cloud", "security"})
	want.SetBuildLimit(10)
	want.SetTimeout(30)
	want.SetCounter(0)
	want.SetVisibility("public")
	want.SetPrivate(false)
	want.SetTrusted(false)
	want.SetActive(true)
	want.SetAllowEvents(e)
	want.SetPipelineType("yaml")
	want.SetPreviousName("oldName")
	want.SetApproveBuild(constants.ApproveNever)
	want.SetInstallID(0)

	// run test
	got := testRepo().ToAPI()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ToAPI() mismatch (-want +got):\n%s", diff)
	}
}

func TestTypes_Repo_Validate(t *testing.T) {
	// setup types
	topics := []string{}
	longTopic := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for len(topics) < 21 {
		topics = append(topics, longTopic)
	}

	// setup tests
	tests := []struct {
		failure bool
		repo    *Repo
	}{
		{
			failure: false,
			repo:    testRepo(),
		},
		{ // no user_id set for repo
			failure: true,
			repo: &Repo{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				Hash:       sql.NullString{String: "superSecretHash", Valid: true},
				Org:        sql.NullString{String: "github", Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				FullName:   sql.NullString{String: "github/octocat", Valid: true},
				Visibility: sql.NullString{String: "public", Valid: true},
			},
		},
		{ // no hash set for repo
			failure: true,
			repo: &Repo{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				UserID:     sql.NullInt64{Int64: 1, Valid: true},
				Org:        sql.NullString{String: "github", Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				FullName:   sql.NullString{String: "github/octocat", Valid: true},
				Visibility: sql.NullString{String: "public", Valid: true},
			},
		},
		{ // no org set for repo
			failure: true,
			repo: &Repo{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				UserID:     sql.NullInt64{Int64: 1, Valid: true},
				Hash:       sql.NullString{String: "superSecretHash", Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				FullName:   sql.NullString{String: "github/octocat", Valid: true},
				Visibility: sql.NullString{String: "public", Valid: true},
			},
		},
		{ // no name set for repo
			failure: true,
			repo: &Repo{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				UserID:     sql.NullInt64{Int64: 1, Valid: true},
				Hash:       sql.NullString{String: "superSecretHash", Valid: true},
				Org:        sql.NullString{String: "github", Valid: true},
				FullName:   sql.NullString{String: "github/octocat", Valid: true},
				Visibility: sql.NullString{String: "public", Valid: true},
			},
		},
		{ // no full_name set for repo
			failure: true,
			repo: &Repo{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				UserID:     sql.NullInt64{Int64: 1, Valid: true},
				Hash:       sql.NullString{String: "superSecretHash", Valid: true},
				Org:        sql.NullString{String: "github", Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				Visibility: sql.NullString{String: "public", Valid: true},
			},
		},
		{ // no visibility set for repo
			failure: true,
			repo: &Repo{
				ID:       sql.NullInt64{Int64: 1, Valid: true},
				UserID:   sql.NullInt64{Int64: 1, Valid: true},
				Hash:     sql.NullString{String: "superSecretHash", Valid: true},
				Org:      sql.NullString{String: "github", Valid: true},
				Name:     sql.NullString{String: "octocat", Valid: true},
				FullName: sql.NullString{String: "github/octocat", Valid: true},
			},
		},
		{ // topics exceed max size
			failure: true,
			repo: &Repo{
				ID:       sql.NullInt64{Int64: 1, Valid: true},
				UserID:   sql.NullInt64{Int64: 1, Valid: true},
				Hash:     sql.NullString{String: "superSecretHash", Valid: true},
				Org:      sql.NullString{String: "github", Valid: true},
				Name:     sql.NullString{String: "octocat", Valid: true},
				FullName: sql.NullString{String: "github/octocat", Valid: true},
				Topics:   topics,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.repo.Validate()

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

func TestTypes_RepoFromAPI(t *testing.T) {
	// setup types
	repo := new(api.Repo)
	owner := testutils.APIUser()
	owner.SetID(1)

	repo.SetID(1)
	repo.SetOwner(owner)
	repo.SetHash("superSecretHash")
	repo.SetOrg("github")
	repo.SetName("octocat")
	repo.SetFullName("github/octocat")
	repo.SetLink("https://github.com/github/octocat")
	repo.SetClone("https://github.com/github/octocat.git")
	repo.SetBranch("main")
	repo.SetTopics([]string{"cloud", "security"})
	repo.SetBuildLimit(10)
	repo.SetTimeout(30)
	repo.SetCounter(0)
	repo.SetVisibility("public")
	repo.SetPrivate(false)
	repo.SetTrusted(false)
	repo.SetActive(true)
	repo.SetAllowEvents(api.NewEventsFromMask(1))
	repo.SetPipelineType("yaml")
	repo.SetPreviousName("oldName")
	repo.SetApproveBuild(constants.ApproveNever)
	repo.SetInstallID(0)

	want := testRepo()
	want.Owner = User{}

	// run test
	got := RepoFromAPI(repo)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("RepoFromAPI() mismatch (-want +got):\n%s", diff)
	}
}

// testRepo is a test helper function to create a Repo
// type with all fields set to a fake value.
func testRepo() *Repo {
	return &Repo{
		ID:           sql.NullInt64{Int64: 1, Valid: true},
		UserID:       sql.NullInt64{Int64: 1, Valid: true},
		Hash:         sql.NullString{String: "superSecretHash", Valid: true},
		Org:          sql.NullString{String: "github", Valid: true},
		Name:         sql.NullString{String: "octocat", Valid: true},
		FullName:     sql.NullString{String: "github/octocat", Valid: true},
		Link:         sql.NullString{String: "https://github.com/github/octocat", Valid: true},
		Clone:        sql.NullString{String: "https://github.com/github/octocat.git", Valid: true},
		Branch:       sql.NullString{String: "main", Valid: true},
		Topics:       []string{"cloud", "security"},
		BuildLimit:   sql.NullInt64{Int64: 10, Valid: true},
		Timeout:      sql.NullInt64{Int64: 30, Valid: true},
		Counter:      sql.NullInt32{Int32: 0, Valid: true},
		Visibility:   sql.NullString{String: "public", Valid: true},
		Private:      sql.NullBool{Bool: false, Valid: true},
		Trusted:      sql.NullBool{Bool: false, Valid: true},
		Active:       sql.NullBool{Bool: true, Valid: true},
		AllowEvents:  sql.NullInt64{Int64: 1, Valid: true},
		PipelineType: sql.NullString{String: "yaml", Valid: true},
		PreviousName: sql.NullString{String: "oldName", Valid: true},
		ApproveBuild: sql.NullString{String: constants.ApproveNever, Valid: true},
		InstallID:    sql.NullInt64{Int64: 0, Valid: true},

		Owner: *testUser(),
	}
}
