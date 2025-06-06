// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestWorker_Engine_GetWorkerForName(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)
	_worker.SetRunningBuilds(nil) // sqlmock cannot parse string array values

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.WorkerFromAPI(_worker)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers" WHERE hostname = $1 LIMIT $2`).WithArgs("worker_0", 1).WillReturnRows(_rows)

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
		database *Engine
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
