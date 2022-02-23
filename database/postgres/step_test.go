// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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

func TestPostgres_Client_GetStep(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectBuildStep, 1, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Step
	}{
		{
			failure: false,
			want:    _step,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetStep(1, _build)

		if test.failure {
			if err == nil {
				t.Errorf("GetStep should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStep returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStep is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "steps" ("build_id","repo_id","number","name","image","stage","status","error","exit_code","created","started","finished","host","runtime","distribution","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING "id"`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
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
		err := _database.CreateStep(_step)

		if test.failure {
			if err == nil {
				t.Errorf("CreateStep should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateStep returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "steps" SET "build_id"=$1,"repo_id"=$2,"number"=$3,"name"=$4,"image"=$5,"stage"=$6,"status"=$7,"error"=$8,"exit_code"=$9,"created"=$10,"started"=$11,"finished"=$12,"host"=$13,"runtime"=$14,"distribution"=$15 WHERE "id" = $16`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
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
		err := _database.UpdateStep(_step)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateStep should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateStep returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteStep(t *testing.T) {
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
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteStep, 1).Statement

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
		err := _database.DeleteStep(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteStep should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteStep returned err: %v", err)
		}
	}
}

// testStep is a test helper function to create a
// library Step type with all fields set to their
// zero values.
func testStep() *library.Step {
	i64 := int64(0)
	i := 0
	str := ""

	return &library.Step{
		ID:           &i64,
		BuildID:      &i64,
		RepoID:       &i64,
		Number:       &i,
		Name:         &str,
		Image:        &str,
		Stage:        &str,
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
