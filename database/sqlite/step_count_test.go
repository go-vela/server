// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

	// create the step table
	err = _database.Sqlite.Exec(ddl.CreateStepTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableStep, err)
	}
}

func TestSqlite_Client_GetBuildStepCount(t *testing.T) {
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
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("bar")
	_stepTwo.SetImage("foo")

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
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

		// create the steps in the database
		err := _database.CreateStep(_stepOne)
		if err != nil {
			t.Errorf("unable to create test step: %v", err)
		}

		err = _database.CreateStep(_stepTwo)
		if err != nil {
			t.Errorf("unable to create test step: %v", err)
		}

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

func TestSqlite_Client_GetStepImageCount(t *testing.T) {
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

func TestSqlite_Client_GetStepStatusCount(t *testing.T) {
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
