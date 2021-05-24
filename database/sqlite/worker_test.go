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

	// create the worker table
	err = _database.Sqlite.Exec(ddl.CreateWorkerTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableWorker, err)
	}
}

func TestSqlite_Client_GetWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Worker
	}{
		{
			failure: false,
			want:    _worker,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		// create the worker in the database
		err := _database.CreateWorker(test.want)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		got, err := _database.GetWorker("worker_0")

		if test.failure {
			if err == nil {
				t.Errorf("GetWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorker returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorker is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetWorkerByAddress(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Worker
	}{
		{
			failure: false,
			want:    _worker,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		// create the worker in the database
		err := _database.CreateWorker(test.want)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		got, err := _database.GetWorkerByAddress("localhost")

		if test.failure {
			if err == nil {
				t.Errorf("GetWorkerByAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorkerByAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorkerByAddress is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_CreateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

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
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		err := _database.CreateWorker(_worker)

		if test.failure {
			if err == nil {
				t.Errorf("CreateWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateWorker returned err: %v", err)
		}
	}
}

func TestSqlite_Client_UpdateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

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
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		// create the worker in the database
		err := _database.CreateWorker(_worker)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		err = _database.UpdateWorker(_worker)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateWorker returned err: %v", err)
		}
	}
}

func TestSqlite_Client_DeleteWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

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
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		// create the worker in the database
		err := _database.CreateWorker(_worker)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		err = _database.DeleteWorker(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteWorker returned err: %v", err)
		}
	}
}

// testWorker is a test helper function to create a
// library Worker type with all fields set to their
// zero values.
func testWorker() *library.Worker {
	i64 := int64(0)
	str := ""
	b := false
	var arr []string

	return &library.Worker{
		ID:            &i64,
		Hostname:      &str,
		Address:       &str,
		Routes:        &arr,
		Active:        &b,
		LastCheckedIn: &i64,
		BuildLimit:    &i64,
	}
}
