// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestWorker_Engine_ListWorkers(t *testing.T) {
	older := time.Now().Unix() - 60
	newer := time.Now().Unix() - 30
	// setup types
	_workerOne := testWorker()
	_workerOne.SetID(1)
	_workerOne.SetHostname("worker_0")
	_workerOne.SetAddress("localhost")
	_workerOne.SetActive(true)
	_workerOne.SetRunningBuilds(nil)
	_workerOne.SetLastCheckedIn(newer)

	_workerTwo := testWorker()
	_workerTwo.SetID(2)
	_workerTwo.SetHostname("worker_1")
	_workerTwo.SetAddress("localhost")
	_workerTwo.SetActive(true)
	_workerTwo.SetLastCheckedIn(older)
	_workerTwo.SetRunningBuilds(nil)

	_workerThree := testWorker()
	_workerThree.SetID(3)
	_workerThree.SetHostname("worker_2")
	_workerThree.SetAddress("localhost")
	_workerThree.SetActive(false)
	_workerThree.SetLastCheckedIn(newer)
	_workerThree.SetRunningBuilds(nil)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.WorkerFromAPI(_workerOne), *types.WorkerFromAPI(_workerTwo), *types.WorkerFromAPI(_workerThree)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "workers" WHERE last_checked_in < $1 AND last_checked_in > $2`).WillReturnRows(_rows)

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

	_, err = _sqlite.CreateWorker(context.TODO(), _workerThree)
	if err != nil {
		t.Errorf("unable to create test worker for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		before   int64
		active   string
		name     string
		database *Engine
		want     []*api.Worker
	}{
		{
			failure:  false,
			before:   newer,
			active:   "all",
			name:     "sqlite3 before filter",
			database: _sqlite,
			want:     []*api.Worker{_workerTwo},
		},
		{
			failure:  false,
			before:   newer + 1,
			active:   "all",
			name:     "postgres catch all",
			database: _postgres,
			want:     []*api.Worker{_workerOne, _workerTwo, _workerThree},
		},
		{
			failure:  false,
			before:   newer + 1,
			active:   "all",
			name:     "sqlite3 catch all",
			database: _sqlite,
			want:     []*api.Worker{_workerOne, _workerTwo, _workerThree},
		},
		{
			failure:  false,
			before:   newer + 1,
			active:   "true",
			name:     "sqlite3 active filter",
			database: _sqlite,
			want:     []*api.Worker{_workerOne, _workerTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListWorkers(context.TODO(), test.active, test.before, 0)

			if test.failure {
				if err == nil {
					t.Errorf("ListWorkers for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListWorkers for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListWorkers() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
