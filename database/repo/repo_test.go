// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/library/actions"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

// testRepo is a test helper function to create a library
// Repo type with all fields set to their zero values.
func testRepo() *library.Repo {
	return &library.Repo{
		ID:           new(int64),
		UserID:       new(int64),
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
		AllowPull:    new(bool),
		AllowPush:    new(bool),
		AllowDeploy:  new(bool),
		AllowTag:     new(bool),
		AllowComment: new(bool),
		AllowEvents:  testEvents(),
	}
}

func testEvents() *library.Events {
	return &library.Events{
		Push: &actions.Push{
			Branch: new(bool),
			Tag:    new(bool),
		},
		PullRequest: &actions.Pull{
			Opened:      new(bool),
			Edited:      new(bool),
			Synchronize: new(bool),
		},
		Deployment: &actions.Deploy{
			Created: new(bool),
		},
		Comment: &actions.Comment{
			Created: new(bool),
			Edited:  new(bool),
		},
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
