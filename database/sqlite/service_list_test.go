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
	"github.com/go-vela/types/library"
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

func TestSqlite_Client_GetServiceList(t *testing.T) {
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
		want    []*library.Service
	}{
		{
			failure: false,
			want:    []*library.Service{_serviceOne, _serviceTwo},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		for _, service := range test.want {
			// create the service in the database
			err := _database.CreateService(service)
			if err != nil {
				t.Errorf("unable to create test service: %v", err)
			}
		}

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

func TestSqlite_Client_GetBuildServiceList(t *testing.T) {
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
		want    []*library.Service
	}{
		{
			failure: false,
			want:    []*library.Service{_serviceTwo, _serviceOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		for _, service := range test.want {
			// create the service in the database
			err := _database.CreateService(service)
			if err != nil {
				t.Errorf("unable to create test service: %v", err)
			}
		}

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
