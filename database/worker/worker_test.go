// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWorker_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
	_mock.ExpectExec(CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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

	_interface, err := New(
		WithClient(_postgres),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new postgres worker engine: %v", err)
	}

	return _interface.(*engine), _mock
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

	_interface, err := New(
		WithClient(_sqlite),
		WithLogger(logrus.NewEntry(logrus.StandardLogger())),
		WithSkipCreation(false),
	)
	if err != nil {
		t.Errorf("unable to create new sqlite worker engine: %v", err)
	}

	return _interface.(*engine)
}

// testWorker is a test helper function to create a library
// Worker type with all fields set to their zero values.
func testWorker() *library.Worker {
	return &library.Worker{
		ID:                  new(int64),
		Hostname:            new(string),
		Address:             new(string),
		Routes:              new([]string),
		Active:              new(bool),
		Status:              new(string),
		LastStatusUpdateAt:  new(int64),
		RunningBuildIDs:     new([]string),
		LastBuildStartedAt:  new(int64),
		LastBuildFinishedAt: new(int64),
		LastCheckedIn:       new(int64),
		BuildLimit:          new(int64),
	}
}
