// SPDX-License-Identifier: Apache-2.0

package user

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
)

func TestUser_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
	_mock.ExpectExec(CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
		t.Errorf("unable to create new postgres user engine: %v", err)
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
		t.Errorf("unable to create new sqlite user engine: %v", err)
	}

	return _engine
}

// testAPIUser is a test helper function to create an API
// User type with all fields set to their zero values.
func testAPIUser() *api.User {
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

func TestUser_Decrypt(t *testing.T) {
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

func TestUser_Encrypt(t *testing.T) {
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

func TestUser_Nullify(t *testing.T) {
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

func TestUser_ToAPI(t *testing.T) {
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

func TestUser_Validate(t *testing.T) {
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
				Favorites: exceededField(),
			},
		},
		{ // invalid dashboards set for user
			failure: true,
			user: &User{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				Name:       sql.NullString{String: "octocat", Valid: true},
				Token:      sql.NullString{String: "superSecretToken", Valid: true},
				Dashboards: exceededField(),
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
	got := FromAPI(u)

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
func exceededField() []string {
	// initialize empty favorites
	values := []string{}

	// add enough strings to exceed the character limit
	for i := 0; i < 500; i++ {
		// construct field
		// use i to adhere to unique favorites
		field := "github/octocat-" + strconv.Itoa(i)

		values = append(values, field)
	}

	return values
}
