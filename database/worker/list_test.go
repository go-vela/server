// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestWorker_Engine_ListWorkers(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "workers"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "last_checked_in", "build_limit"}).
		AddRow(1, "worker_0", "localhost", nil, true, 0, 0).
		AddRow(2, "worker_1", "localhost", nil, true, 0, 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateWorker(_workerOne)
	if err != nil {
		t.Errorf("unable to create test worker for sqlite: %v", err)
	}

	err = _sqlite.CreateWorker(_workerTwo)
	if err != nil {
		t.Errorf("unable to create test worker for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Worker
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Worker{_workerOne, _workerTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Worker{_workerOne, _workerTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListWorkers()

			if test.failure {
				if err == nil {
					t.Errorf("ListWorkers for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListWorkers for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListWorkers for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
