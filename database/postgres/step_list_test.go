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

func TestPostgres_Client_GetStepList(t *testing.T) {
	// setup types
	_stepOne := testStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")

	_stepTwo := testStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(1)
	_stepTwo.SetName("bar")
	_stepTwo.SetImage("foo")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListSteps).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 1, 1, "bar", "foo", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Step
	}{
		{
			failure: false,
			want:    []*library.Step{_stepOne, _stepTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetStepList()

		if test.failure {
			if err == nil {
				t.Errorf("GetStepList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStepList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetBuildStepList(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_stepOne := testStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")

	_stepTwo := testStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(1)
	_stepTwo.SetName("bar")
	_stepTwo.SetImage("foo")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListBuildSteps, 1, 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 1, 1, "bar", "foo", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Step
	}{
		{
			failure: false,
			want:    []*library.Step{_stepOne, _stepTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuildStepList(_build, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildStepList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildStepList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildStepList is %v, want %v", got, test.want)
		}
	}
}
