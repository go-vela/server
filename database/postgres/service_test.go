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
)

func TestPostgres_Client_GetService(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectBuildService).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Service
	}{
		{
			failure: false,
			want:    _service,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetService(1, _build)

		if test.failure {
			if err == nil {
				t.Errorf("GetService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetService returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetService is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "services" ("build_id","repo_id","number","name","image","status","error","exit_code","created","started","finished","host","runtime","distribution","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING "id"`).
		WithArgs(1, 1, 1, "foo", "bar", "", "", nil, nil, nil, nil, "", "", "", 1).
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
		err := _database.CreateService(_service)

		if test.failure {
			if err == nil {
				t.Errorf("CreateService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateService returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "services" SET "build_id"=$1,"repo_id"=$2,"number"=$3,"name"=$4,"image"=$5,"status"=$6,"error"=$7,"exit_code"=$8,"created"=$9,"started"=$10,"finished"=$11,"host"=$12,"runtime"=$13,"distribution"=$14 WHERE "id" = $15`).
		WithArgs(1, 1, 1, "foo", "bar", "", "", nil, nil, nil, nil, "", "", "", 1).
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
		err := _database.UpdateService(_service)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateService returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteService(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(dml.DeleteService).WillReturnResult(sqlmock.NewResult(1, 1))

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
		err := _database.DeleteService(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteService should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteService returned err: %v", err)
		}
	}
}

// testService is a test helper function to create a
// library Service type with all fields set to their
// zero values.
func testService() *library.Service {
	i64 := int64(0)
	i := 0
	str := ""

	return &library.Service{
		ID:           &i64,
		BuildID:      &i64,
		RepoID:       &i64,
		Number:       &i,
		Name:         &str,
		Image:        &str,
		Status:       &str,
		Error:        &str,
		ExitCode:     &i,
		Created:      &i64,
		Started:      &i64,
		Finished:     &i64,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
