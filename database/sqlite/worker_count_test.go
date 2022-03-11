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

	// create the worker table
	err = _database.Sqlite.Exec(ddl.CreateWorkerTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableWorker, err)
	}
}

func TestSqlite_Client_GetWorkerCount(t *testing.T) {
	// setup types
	_workerOne := testWorker()
	_workerOne.SetID(1)
	_workerOne.SetHostname("worker_0")
	_workerOne.SetAddress("localhost")
	_workerOne.SetActive(true)

	_workerTwo := testWorker()
	_workerTwo.SetID(2)
	_workerTwo.SetHostname("worker_1")
	_workerTwo.SetAddress("localhost")
	_workerTwo.SetActive(true)

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
		// defer cleanup of the workers table
		defer _database.Sqlite.Exec("delete from workers;")

		// create the workers in the database
		err := _database.CreateWorker(_workerOne)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		err = _database.CreateWorker(_workerTwo)
		if err != nil {
			t.Errorf("unable to create test worker: %v", err)
		}

		got, err := _database.GetWorkerCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetWorkerCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorkerCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorkerCount is %v, want %v", got, test.want)
		}
	}
}
