// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/types/constants"
)

func init() {
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		log.Fatalf("unable to create new sqlite test database: %v", err)
	}

	// create the service table
	err = _database.Sqlite.Exec(ddl.CreateServiceTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableService, err)
	}
}

func TestSqlite_Client_GetBuildServiceCount(t *testing.T) {
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
	_serviceTwo.SetNumber(2)
	_serviceTwo.SetName("bar")
	_serviceTwo.SetImage("foo")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		// create the services in the database
		_, err := _database.CreateService(_serviceOne)
		if err != nil {
			t.Errorf("unable to create test service: %v", err)
		}

		_, err = _database.CreateService(_serviceTwo)
		if err != nil {
			t.Errorf("unable to create test service: %v", err)
		}

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

func TestSqlite_Client_GetServiceImageCount(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    map[string]float64
	}{
		{
			failure: false,
			want:    map[string]float64{},
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

func TestSqlite_Client_GetServiceStatusCount(t *testing.T) {
	// setup types
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
