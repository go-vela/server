// SPDX-License-Identifier: Apache-2.0

package step

import (
	"database/sql/driver"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStep_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))

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
			logger:       logger,
			skipCreation: false,
			want: &engine{
				client: _postgres,
				config: &config{SkipCreation: false},
				logger: logger,
			},
		},
		{
			failure:      false,
			name:         "sqlite3",
			client:       _sqlite,
			logger:       logger,
			skipCreation: false,
			want: &engine{
				client: _sqlite,
				config: &config{SkipCreation: false},
				logger: logger,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := New(
				WithClient(test.client),
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
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new postgres step engine: %v", err)
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
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new sqlite step engine: %v", err)
	}

	return _engine
}

// testBuild is a test helper function to create a library
// Build type with all fields set to their zero values.
func testBuild() *library.Build {
	return &library.Build{
		ID:           new(int64),
		RepoID:       new(int64),
		PipelineID:   new(int64),
		Number:       new(int),
		Parent:       new(int),
		Event:        new(string),
		EventAction:  new(string),
		Status:       new(string),
		Error:        new(string),
		Enqueued:     new(int64),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Deploy:       new(string),
		Clone:        new(string),
		Source:       new(string),
		Title:        new(string),
		Message:      new(string),
		Commit:       new(string),
		Sender:       new(string),
		Author:       new(string),
		Email:        new(string),
		Link:         new(string),
		Branch:       new(string),
		Ref:          new(string),
		BaseRef:      new(string),
		HeadRef:      new(string),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
	}
}

// testStep is a test helper function to create a library
// Step type with all fields set to their zero values.
func testStep() *library.Step {
	return &library.Step{
		ID:           new(int64),
		BuildID:      new(int64),
		RepoID:       new(int64),
		Number:       new(int),
		Name:         new(string),
		Image:        new(string),
		Stage:        new(string),
		Status:       new(string),
		Error:        new(string),
		ExitCode:     new(int),
		Created:      new(int64),
		Started:      new(int64),
		Finished:     new(int64),
		Host:         new(string),
		Runtime:      new(string),
		Distribution: new(string),
	}
}

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime

// NowTimestamp is used to test whether timestamps get updated correctly to the current time with lenience.
type NowTimestamp struct{}

// Match satisfies sqlmock.Argument interface.
func (t NowTimestamp) Match(v driver.Value) bool {
	ts, ok := v.(int64)
	if !ok {
		return false
	}
	now := time.Now().Unix()

	return now-ts < 10
}
