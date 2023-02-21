// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/database/hook"
	"github.com/go-vela/server/database/pipeline"
	"github.com/go-vela/server/database/postgres/ddl"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/server/database/user"
	"github.com/go-vela/server/database/worker"
	"github.com/go-vela/types/library"
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
	_mock.ExpectExec(ddl.CreateLogTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateServiceTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateStepTable).WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the index queries
	_mock.ExpectExec(ddl.CreateBuildRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateBuildSourceIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogBuildIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgRepo).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgTeam).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrg).WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the hook queries
	_mock.ExpectExec(hook.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(hook.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the pipeline queries
	_mock.ExpectExec(pipeline.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(pipeline.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the repo queries
	_mock.ExpectExec(repo.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(repo.CreateOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the user queries
	_mock.ExpectExec(user.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(user.CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the worker queries
	_mock.ExpectExec(worker.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(worker.CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
	_mock.ExpectExec(ddl.CreateLogTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateServiceTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateStepTable).WillReturnResult(sqlmock.NewResult(1, 1))

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
	_mock.ExpectExec(ddl.CreateBuildSourceIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateLogBuildIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgRepo).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrgTeam).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(ddl.CreateSecretTypeOrg).WillReturnResult(sqlmock.NewResult(1, 1))

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

func TestPostgres_createServices(t *testing.T) {
	// setup types
	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the hook queries
	_mock.ExpectExec(hook.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(hook.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the pipeline queries
	_mock.ExpectExec(pipeline.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(pipeline.CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the repo queries
	_mock.ExpectExec(repo.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(repo.CreateOrgNameIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the user queries
	_mock.ExpectExec(user.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(user.CreateUserRefreshIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	// ensure the mock expects the worker queries
	_mock.ExpectExec(worker.CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(worker.CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
		PipelineType: &str,
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
