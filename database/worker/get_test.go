// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	api "github.com/go-vela/server/api/types"
)

func TestWorker_Engine_GetWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)
	_worker.SetRunningBuilds(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "last_checked_in", "build_limit"}).
		AddRow(1, "worker_0", "localhost", nil, true, 0, 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

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
		want     *api.Worker
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
			got, err := test.database.GetWorker(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetWorker for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetWorker for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetWorker for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
