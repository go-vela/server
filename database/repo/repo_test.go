// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/types/constants"
)

func TestRepo_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))

	_config := &gorm.Config{SkipDefaultTransaction: true}

	_postgres, err := gorm.Open(postgres.New(postgres.Config{Conn: _sql}), _config)
	if err != nil {
		t.Errorf("unable to create new postgres database: %v", err)
	}

	_sqlite, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), _config)
	if err != nil {
		t.Errorf("unable to create new sqlite database: %v", err)
	}

	defer func() { _sql, _ := _sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure      bool
		name         string
		client       *gorm.DB
		key          string
		logger       *logrus.Entry
		skipCreation bool
		want         *engine
	}{
		{
			failure:      false,
			name:         "postgres",
			client:       _postgres,
			key:          "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			logger:       logger,
			skipCreation: false,
			want: &engine{
				client: _postgres,
				config: &config{EncryptionKey: "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW", SkipCreation: false},
				logger: logger,
			},
		},
		{
			failure:      false,
			name:         "sqlite3",
			client:       _sqlite,
			key:          "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			logger:       logger,
			skipCreation: false,
			want: &engine{
				client: _sqlite,
				config: &config{EncryptionKey: "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW", SkipCreation: false},
				logger: logger,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(
				WithClient(test.client),
				WithEncryptionKey(test.key),
				WithLogger(test.logger),
				WithSkipCreation(test.skipCreation),
			)

			if test.failure {
				if err == nil {
					t.Errorf("New for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("New for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("New for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

// testPostgres is a helper function to create a Postgres engine for testing.
func testPostgres(t *testing.T) (*engine, sqlmock.Sqlmock) {
	// create the new mock sql database
	//
	// https://pkg.go.dev/github.com/DATA-DOG/go-sqlmock#New
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))

	// create the new mock Postgres database client
	//
	// https://pkg.go.dev/gorm.io/gorm#Open
	_postgres, err := gorm.Open(
		postgres.New(postgres.Config{Conn: _sql}),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new postgres database: %v", err)
	}

	_engine, err := New(
		WithClient(_postgres),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new postgres repo engine: %v", err)
	}

	return _engine, _mock
}

// testSqlite is a helper function to create a Sqlite engine for testing.
func testSqlite(t *testing.T) *engine {
	_sqlite, err := gorm.Open(
		sqlite.Open("file::memory:?cache=shared"),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		t.Errorf("unable to create new sqlite database: %v", err)
	}

	_engine, err := New(
		WithClient(_sqlite),
		WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new sqlite repo engine: %v", err)
	}

	return _engine
}

// testAPIRepo is a test helper function to create an API
// Repo type with all fields set to their zero values.
func testAPIRepo() *api.Repo {
	return &api.Repo{
		ID:           new(int64),
		Owner:        testOwner(),
		BuildLimit:   new(int64),
		Timeout:      new(int64),
		Counter:      new(int),
		PipelineType: new(string),
		Hash:         new(string),
		Org:          new(string),
		Name:         new(string),
		FullName:     new(string),
		Link:         new(string),
		Clone:        new(string),
		Branch:       new(string),
		Visibility:   new(string),
		PreviousName: new(string),
		ApproveBuild: new(string),
		Private:      new(bool),
		Trusted:      new(bool),
		Active:       new(bool),
		AllowEvents:  testEvents(),
	}
}

func testEvents() *api.Events {
	return &api.Events{
		Push: &actions.Push{
			Branch:       new(bool),
			Tag:          new(bool),
			DeleteBranch: new(bool),
			DeleteTag:    new(bool),
		},
		PullRequest: &actions.Pull{
			Opened:      new(bool),
			Edited:      new(bool),
			Synchronize: new(bool),
			Reopened:    new(bool),
			Labeled:     new(bool),
			Unlabeled:   new(bool),
		},
		Deployment: &actions.Deploy{
			Created: new(bool),
		},
		Comment: &actions.Comment{
			Created: new(bool),
			Edited:  new(bool),
		},
		Schedule: &actions.Schedule{
			Run: new(bool),
		},
	}
}

func testOwner() *api.User {
	return &api.User{
		ID:           new(int64),
		Name:         new(string),
		RefreshToken: new(string),
		Token:        new(string),
		Favorites:    new([]string),
		Active:       new(bool),
		Admin:        new(bool),
		Dashboards:   new([]string),
	}
}

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyArgument) Match(_ driver.Value) bool {
	return true
}

func TestRepo_Decrypt(t *testing.T) {
	// setup types
	key := "C639A572E14D5075C526FDDD43E4ECF6"
	encrypted := testRepo()

	err := encrypted.Encrypt(key)
	if err != nil {
		t.Errorf("unable to encrypt repo: %v", err)
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

func TestRepo_Encrypt(t *testing.T) {
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

func TestRepo_Nullify(t *testing.T) {
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

func TestRepo_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Repo)
	e := api.NewEventsFromMask(1)
	owner := testOwner()

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

	// run test
	got := testRepo().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestRepo_Validate(t *testing.T) {
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

func TestRepo_FromAPI(t *testing.T) {
	// setup types
	r := new(api.Repo)
	owner := testOwner()
	owner.SetID(1)

	r.SetID(1)
	r.SetOwner(owner)
	r.SetHash("superSecretHash")
	r.SetOrg("github")
	r.SetName("octocat")
	r.SetFullName("github/octocat")
	r.SetLink("https://github.com/github/octocat")
	r.SetClone("https://github.com/github/octocat.git")
	r.SetBranch("main")
	r.SetTopics([]string{"cloud", "security"})
	r.SetBuildLimit(10)
	r.SetTimeout(30)
	r.SetCounter(0)
	r.SetVisibility("public")
	r.SetPrivate(false)
	r.SetTrusted(false)
	r.SetActive(true)
	r.SetAllowEvents(api.NewEventsFromMask(1))
	r.SetPipelineType("yaml")
	r.SetPreviousName("oldName")
	r.SetApproveBuild(constants.ApproveNever)

	want := testRepo()

	// run test
	got := FromAPI(r)

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FromAPI() mismatch (-want +got):\n%s", diff)
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
	}
}
