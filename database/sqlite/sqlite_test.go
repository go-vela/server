// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"testing"
	"time"

	"github.com/go-vela/types/library"
)

func TestSqlite_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: false,
			address: ":memory:",
			want:    ":memory:",
		},
		{
			failure: true,
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(
			WithAddress(test.address),
			WithCompressionLevel(3),
			WithConnectionLife(10*time.Second),
			WithConnectionIdle(5),
			WithConnectionOpen(20),
			WithEncryptionKey("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW"),
			WithSkipCreation(false),
		)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}

func TestSqlite_setupDatabase(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup the skip test database client
	_skipDatabase, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new skip sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _skipDatabase.Sqlite.DB(); _sql.Close() }()

	err = WithSkipCreation(true)(_skipDatabase)
	if err != nil {
		t.Errorf("unable to set SkipCreation for sqlite test database: %v", err)
	}

	tests := []struct {
		failure  bool
		database *client
	}{
		{
			failure:  false,
			database: _database,
		},
		{
			failure:  false,
			database: _skipDatabase,
		},
	}

	// run tests
	for _, test := range tests {
		err := setupDatabase(test.database)

		if test.failure {
			if err == nil {
				t.Errorf("setupDatabase should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("setupDatabase returned err: %v", err)
		}
	}
}

func TestSqlite_createTables(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := createTables(_database)

		if test.failure {
			if err == nil {
				t.Errorf("createTables should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("createTables returned err: %v", err)
		}
	}
}

func TestSqlite_createIndexes(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := createIndexes(_database)

		if test.failure {
			if err == nil {
				t.Errorf("createIndexes should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("createIndexes returned err: %v", err)
		}
	}
}

func TestSqlite_createServices(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := createServices(_database)

		if test.failure {
			if err == nil {
				t.Errorf("createServices should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("createServices returned err: %v", err)
		}
	}
}

// testRepo is a test helper function to create a
// library Repo type with all fields set to their
// zero values.
func testRepo() *library.Repo {
	i64 := int64(0)
	i := 0
	str := ""
	b := false

	return &library.Repo{
		ID:           &i64,
		UserID:       &i64,
		Hash:         &str,
		Org:          &str,
		Name:         &str,
		FullName:     &str,
		Link:         &str,
		Clone:        &str,
		Branch:       &str,
		BuildLimit:   &i64,
		Timeout:      &i64,
		Counter:      &i,
		Visibility:   &str,
		Private:      &b,
		Trusted:      &b,
		Active:       &b,
		AllowPull:    &b,
		AllowPush:    &b,
		AllowDeploy:  &b,
		AllowTag:     &b,
		AllowComment: &b,
		PreviousName: &str,
	}
}
