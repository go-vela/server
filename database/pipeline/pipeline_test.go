// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"database/sql/driver"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/library"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/database"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPipeline_New(t *testing.T) {
	// setup types
	logger := logrus.NewEntry(logrus.StandardLogger())

	_sql, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new SQL mock: %v", err)
	}
	defer _sql.Close()

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
		name     string
		database *gorm.DB
		want     *engine
	}{
		{
			name:     "postgres",
			database: _postgres,
			want:     &engine{client: _postgres, logger: logger},
		},
		{
			name:     "sqlite",
			database: _sqlite,
			want:     &engine{client: _sqlite, logger: logger},
		},
	}

	// run tests
	for _, test := range tests {
		got := New(test.database, logger, 0)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("New for %s is %v, want %v", test.name, got, test.want)
		}
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

	return &engine{client: _postgres, logger: logrus.NewEntry(logrus.StandardLogger()), compressionLevel: 3}, _mock
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

	err = _sqlite.AutoMigrate(&database.Pipeline{})
	if err != nil {
		t.Errorf("unable to create pipeline schema for sqlite: %v", err)
	}

	return &engine{client: _sqlite, logger: logrus.NewEntry(logrus.StandardLogger()), compressionLevel: 3}
}

// testPipeline is a test helper function to create a
// library Pipeline type with all fields set to their
// zero values.
func testPipeline() *library.Pipeline {
	return &library.Pipeline{
		ID:        new(int64),
		RepoID:    new(int64),
		Number:    new(int),
		Flavor:    new(string),
		Platform:  new(string),
		Ref:       new(string),
		Type:      new(string),
		Version:   new(string),
		Services:  new(bool),
		Stages:    new(bool),
		Steps:     new(bool),
		Templates: new(bool),
		Data:      new([]byte),
	}
}

// This will be used with the github.com/DATA-DOG/go-sqlmock
// library to compare values that are otherwise not easily
// compared. These typically would be values generated before
// adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyArgument) Match(v driver.Value) bool {
	return true
}
