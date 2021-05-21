// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetBuildLogs(t *testing.T) {
	// setup types
	_logOne := testLog()
	_logOne.SetID(1)
	_logOne.SetStepID(1)
	_logOne.SetBuildID(1)
	_logOne.SetRepoID(1)
	_logOne.SetData([]byte{})

	_logTwo := testLog()
	_logTwo.SetID(2)
	_logTwo.SetServiceID(1)
	_logTwo.SetBuildID(1)
	_logTwo.SetRepoID(1)
	_logTwo.SetData([]byte{})

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListBuildLogs, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "service_id", "step_id", "data"},
	).AddRow(1, 1, 1, 0, 1, []byte{}).AddRow(2, 1, 1, 1, 0, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Log
	}{
		{
			failure: false,
			want:    []*library.Log{_logOne, _logTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuildLogs(1)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildLogs should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildLogs returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildLogs is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetStepLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectStepLog, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "service_id", "step_id", "data"},
	).AddRow(1, 1, 1, 0, 1, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Log
	}{
		{
			failure: false,
			want:    _log,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetStepLog(1)

		if test.failure {
			if err == nil {
				t.Errorf("GetStepLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepLog returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStepLog is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetServiceLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetServiceID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectServiceLog, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "service_id", "step_id", "data"},
	).AddRow(1, 1, 1, 1, 0, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Log
	}{
		{
			failure: false,
			want:    _log,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetServiceLog(1)

		if test.failure {
			if err == nil {
				t.Errorf("GetServiceLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetServiceLog returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetServiceLog is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "logs" ("build_id","repo_id","service_id","step_id","data","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(1, 1, nil, 1, AnyArgument{}, 1).
		WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.CreateLog(_log)

		if test.failure {
			if err == nil {
				t.Errorf("CreateLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateLog returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "logs" SET "build_id"=$1,"repo_id"=$2,"service_id"=$3,"step_id"=$4,"data"=$5 WHERE "id" = $6`).
		WithArgs(1, 1, nil, 1, AnyArgument{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.UpdateLog(_log)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateLog returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteLog(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteLog, 1).Statement

	// ensure the mock expects the query
	_mock.ExpectExec(_query.SQL.String()).WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.DeleteLog(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteLog returned err: %v", err)
		}
	}
}

// testLog is a test helper function to create a
// library Log type with all fields set to their
// zero values.
func testLog() *library.Log {
	i64 := int64(0)
	b := []byte{}

	return &library.Log{
		ID:        &i64,
		BuildID:   &i64,
		RepoID:    &i64,
		ServiceID: &i64,
		StepID:    &i64,
		Data:      &b,
	}
}
