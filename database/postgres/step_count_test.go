// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
)

func TestPostgres_Client_GetBuildStepCount(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectBuildStepsCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuildStepCount(_build)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildStepCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildStepCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildStepCount is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetStepImageCount(t *testing.T) {
	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"image", "count"}).AddRow("foo", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectStepImagesCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    map[string]float64
	}{
		{
			failure: false,
			want:    map[string]float64{"foo": 0},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetStepImageCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetStepImageCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepImageCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStepImageCount is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetStepStatusCount(t *testing.T) {
	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("failure", 0).
		AddRow("killed", 0).
		AddRow("pending", 0).
		AddRow("running", 0).
		AddRow("success", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectStepStatusesCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    map[string]float64
	}{
		{
			failure: false,
			want: map[string]float64{
				"pending": 0,
				"failure": 0,
				"killed":  0,
				"running": 0,
				"success": 0,
			},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetStepStatusCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetStepStatusCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepStatusCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStepStatusCount is %v, want %v", got, test.want)
		}
	}
}
