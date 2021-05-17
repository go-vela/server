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

func TestPostgres_Client_GetServiceList(t *testing.T) {
	// setup types
	_serviceOne := testService()
	_serviceOne.SetID(1)
	_serviceOne.SetRepoID(1)
	_serviceOne.SetBuildID(1)
	_serviceOne.SetNumber(1)
	_serviceOne.SetName("foo")
	_serviceOne.SetImage("bar")

	_serviceTwo := testService()
	_serviceTwo.SetID(2)
	_serviceTwo.SetRepoID(1)
	_serviceTwo.SetBuildID(1)
	_serviceTwo.SetNumber(1)
	_serviceTwo.SetName("bar")
	_serviceTwo.SetImage("foo")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 1, 1, "bar", "foo", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListServices).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Service
	}{
		{
			failure: false,
			want:    []*library.Service{_serviceOne, _serviceTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetServiceList()

		if test.failure {
			if err == nil {
				t.Errorf("GetServiceList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetServiceList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetServiceList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetBuildServiceList(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_serviceOne := testService()
	_serviceOne.SetID(1)
	_serviceOne.SetRepoID(1)
	_serviceOne.SetBuildID(1)
	_serviceOne.SetNumber(1)
	_serviceOne.SetName("foo")
	_serviceOne.SetImage("bar")

	_serviceTwo := testService()
	_serviceTwo.SetID(2)
	_serviceTwo.SetRepoID(1)
	_serviceTwo.SetBuildID(1)
	_serviceTwo.SetNumber(1)
	_serviceTwo.SetName("bar")
	_serviceTwo.SetImage("foo")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"},
	).AddRow(1, 1, 1, 1, "foo", "bar", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 1, 1, "bar", "foo", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListBuildServices).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Service
	}{
		{
			failure: false,
			want:    []*library.Service{_serviceOne, _serviceTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuildServiceList(_build, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildServiceList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildServiceList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildServiceList is %v, want %v", got, test.want)
		}
	}
}
