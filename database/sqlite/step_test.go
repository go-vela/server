// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetStep(t *testing.T) {
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
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		if test.want != nil {
			// create the step in the database
			err := _database.CreateStep(test.want)
			if err != nil {
				t.Errorf("unable to create test step: %v", err)
			}
		}

		got, err := _database.GetStep(1, _build)

		// cleanup the steps table
		_ = _database.Sqlite.Exec("DELETE FROM steps;")

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

func TestSqlite_Client_CreateStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

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
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

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

func TestSqlite_Client_UpdateStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

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
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

		// create the step in the database
		err := _database.CreateStep(_step)
		if err != nil {
			t.Errorf("unable to create test step: %v", err)
		}

		err = _database.UpdateStep(_step)

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

func TestSqlite_Client_DeleteStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

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
		// defer cleanup of the steps table
		defer _database.Sqlite.Exec("delete from steps;")

		// create the step in the database
		err := _database.CreateStep(_step)
		if err != nil {
			t.Errorf("unable to create test step: %v", err)
		}

		err = _database.DeleteStep(1)

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
