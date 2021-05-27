// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetBuildLogs(t *testing.T) {
	// setup types
	_logOne := testLog()
	_logOne.SetID(1)
	_logOne.SetStepID(1)
	_logOne.SetBuildID(1)
	_logOne.SetRepoID(1)
	_logOne.SetData([]byte{})

	_logTwo := testLog()
	_logTwo.SetID(2)
	_logTwo.SetServiceID(1)
	_logTwo.SetBuildID(1)
	_logTwo.SetRepoID(1)
	_logTwo.SetData([]byte{})

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Log
	}{
		{
			failure: false,
			want:    []*library.Log{_logTwo, _logOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the logs table
		defer _database.Sqlite.Exec("delete from logs;")

		for _, log := range test.want {
			// create the log in the database
			err := _database.CreateLog(log)
			if err != nil {
				t.Errorf("unable to create test log: %v", err)
			}
		}

		got, err := _database.GetBuildLogs(1)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildLogs should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildLogs returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildLogs is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetStepLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Log
	}{
		{
			failure: false,
			want:    _log,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the log in the database
			err := _database.CreateLog(test.want)
			if err != nil {
				t.Errorf("unable to create test log: %v", err)
			}
		}

		got, err := _database.GetStepLog(1)

		// cleanup the logs table
		_ = _database.Sqlite.Exec("DELETE FROM logs;")

		if test.failure {
			if err == nil {
				t.Errorf("GetStepLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetStepLog returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetStepLog is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetServiceLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetServiceID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Log
	}{
		{
			failure: false,
			want:    _log,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the log in the database
			err := _database.CreateLog(test.want)
			if err != nil {
				t.Errorf("unable to create test log: %v", err)
			}
		}

		got, err := _database.GetServiceLog(1)

		// cleanup the logs table
		_ = _database.Sqlite.Exec("DELETE FROM logs;")

		if test.failure {
			if err == nil {
				t.Errorf("GetServiceLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetServiceLog returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetServiceLog is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_CreateLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

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
		// defer cleanup of the logs table
		defer _database.Sqlite.Exec("delete from logs;")

		err := _database.CreateLog(_log)

		if test.failure {
			if err == nil {
				t.Errorf("CreateLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateLog returned err: %v", err)
		}
	}
}

func TestSqlite_Client_UpdateLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

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
		// defer cleanup of the logs table
		defer _database.Sqlite.Exec("delete from logs;")

		// create the log in the database
		err := _database.CreateLog(_log)
		if err != nil {
			t.Errorf("unable to create test log: %v", err)
		}

		err = _database.UpdateLog(_log)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateLog returned err: %v", err)
		}
	}
}

func TestSqlite_Client_DeleteLog(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetStepID(1)
	_log.SetBuildID(1)
	_log.SetRepoID(1)
	_log.SetData([]byte{})

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
		// defer cleanup of the logs table
		defer _database.Sqlite.Exec("delete from logs;")

		// create the log in the database
		err := _database.CreateLog(_log)
		if err != nil {
			t.Errorf("unable to create test log: %v", err)
		}

		err = _database.DeleteLog(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteLog should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteLog returned err: %v", err)
		}
	}
}

// testLog is a test helper function to create a
// library Log type with all fields set to their
// zero values.
func testLog() *library.Log {
	i64 := int64(0)
	b := []byte{}

	return &library.Log{
		ID:        &i64,
		BuildID:   &i64,
		RepoID:    &i64,
		ServiceID: &i64,
		StepID:    &i64,
		Data:      &b,
	}
}
