// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetWorkerList(t *testing.T) {
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
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListWorkers).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "last_checked_in", "build_limit"},
	).AddRow(1, "worker_0", "localhost", "{}", true, 0, 0).
		AddRow(2, "worker_1", "localhost", "{}", true, 0, 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Worker
	}{
		{
			failure: false,
			want:    []*library.Worker{_workerOne, _workerTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetWorkerList()

		if test.failure {
			if err == nil {
				t.Errorf("GetWorkerList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorkerList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorkerList is %v, want %v", got, test.want)
		}
	}
}
