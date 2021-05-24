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
	"github.com/go-vela/types/library"
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

func TestSqlite_Client_GetStepList(t *testing.T) {
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
		want    []*library.Step
	}{
		{
			failure: false,
			want:    []*library.Step{_stepOne, _stepTwo},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

		for _, step := range test.want {
			// create the step in the database
			err := _database.CreateStep(step)
			if err != nil {
				t.Errorf("unable to create test step: %v", err)
			}
		}

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

func TestSqlite_Client_GetBuildStepList(t *testing.T) {
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
		want    []*library.Step
	}{
		{
			failure: false,
			want:    []*library.Step{_stepTwo, _stepOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

		for _, step := range test.want {
			// create the step in the database
			err := _database.CreateStep(step)
			if err != nil {
				t.Errorf("unable to create test step: %v", err)
			}
		}

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
