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

func TestPostgres_Client_GetBuildServiceCount(t *testing.T) {
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
	_mock.ExpectQuery(dml.SelectBuildServicesCount).WillReturnRows(_rows)

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
		got, err := _database.GetBuildServiceCount(_build)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildServiceCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildServiceCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildServiceCount is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetServiceImageCount(t *testing.T) {
	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"image", "count"}).AddRow("foo", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectServiceImagesCount).WillReturnRows(_rows)

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
		got, err := _database.GetServiceImageCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetServiceImageCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetServiceImageCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetServiceImageCount is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetServiceStatusCount(t *testing.T) {
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
	_mock.ExpectQuery(dml.SelectServiceStatusesCount).WillReturnRows(_rows)

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
		got, err := _database.GetServiceStatusCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetServiceStatusCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetServiceStatusCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetServiceStatusCount is %v, want %v", got, test.want)
		}
	}
}
