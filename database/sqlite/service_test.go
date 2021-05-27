// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetService(t *testing.T) {
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
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		// create the service in the database
		err := _database.CreateService(test.want)
		if err != nil {
			t.Errorf("unable to create test service: %v", err)
		}

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

func TestSqlite_Client_CreateService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

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

func TestSqlite_Client_UpdateService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		// create the service in the database
		err := _database.CreateService(_service)
		if err != nil {
			t.Errorf("unable to create test service: %v", err)
		}

		err = _database.UpdateService(_service)

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

func TestSqlite_Client_DeleteService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the services table
		defer _database.Sqlite.Exec("delete from services;")

		// create the service in the database
		err := _database.CreateService(_service)
		if err != nil {
			t.Errorf("unable to create test service: %v", err)
		}

		err = _database.DeleteService(1)

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
