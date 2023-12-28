// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestWorker_Engine_GetWorkerForName(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "status", "last_status_update_at", "running_build_ids", "last_build_started_at", "last_build_finished_at", "last_checked_in", "build_limit"}).
		AddRow(1, "worker_0", "localhost", nil, true, nil, 0, nil, 0, 0, 0, 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers" WHERE hostname = $1 LIMIT 1`).WithArgs("worker_0").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateWorker(context.TODO(), _worker)
	if err != nil {
		t.Errorf("unable to create test worker for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Worker
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _worker,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _worker,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetWorkerForHostname(context.TODO(), "worker_0")

			if test.failure {
				if err == nil {
					t.Errorf("GetWorkerForHostname for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetWorkerForHostname for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetWorkerForHostname for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
