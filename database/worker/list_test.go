// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
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
		[]string{"id", "hostname", "address", "routes", "active", "status", "last_status_update_at", "running_build_ids", "last_build_started_at", "last_build_finished_at", "last_checked_in", "build_limit"}).
		AddRow(1, "worker_0", "localhost", nil, true, nil, 0, nil, 0, 0, 0, 0).
		AddRow(2, "worker_1", "localhost", nil, true, nil, 0, nil, 0, 0, 0, 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateWorker(context.TODO(), _workerOne)
	if err != nil {
		t.Errorf("unable to create test worker for sqlite: %v", err)
	}

	_, err = _sqlite.CreateWorker(context.TODO(), _workerTwo)
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
			got, err := test.database.ListWorkers(context.TODO())

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
