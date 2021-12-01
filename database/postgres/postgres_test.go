// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"database/sql/driver"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/database/postgres/ddl"
)

func TestPostgres_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: true,
			address: "postgres://foo:bar@localhost:5432/vela",
			want:    "postgres://foo:bar@localhost:5432/vela",
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

func TestPostgres_setupDatabase(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the ping
	_mock.ExpectPing()

	// ensure the mock expects the table queries
	_mock.ExpectExec(ddl.CreateBuildTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateHookTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateRepoTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateServiceTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateStepTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateUserTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateWorkerTable).WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the index queries
	_mock.ExpectExec(ddl.CreateBuildRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateHookRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogBuildIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateRepoOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgRepo).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgTeam).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrg).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateWorkerHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

	// setup the skip test database client
	_skipDatabase, _skipMock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new skip postgres test database: %v", err)
	}
	defer func() { _sql, _ := _skipDatabase.Postgres.DB(); _sql.Close() }()

	err = WithSkipCreation(true)(_skipDatabase)
	if err != nil {
		t.Errorf("unable to set SkipCreation for postgres test database: %v", err)
	}

	// ensure the mock expects the ping
	_skipMock.ExpectPing()

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

func TestPostgres_createTables(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the table queries
	_mock.ExpectExec(ddl.CreateBuildTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateHookTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateRepoTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateServiceTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateStepTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateUserTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateWorkerTable).WillReturnResult(sqlmock.NewResult(1, 1))

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
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the index queries
	_mock.ExpectExec(ddl.CreateBuildRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateHookRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogBuildIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateRepoOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgRepo).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgTeam).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrg).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateWorkerHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
