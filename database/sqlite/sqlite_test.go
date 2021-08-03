// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"testing"
	"time"
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

func TestPostgres_createIndexes(t *testing.T) {
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
