// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/google/go-cmp/cmp"
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
	_workerOne.SetLastCheckedIn(newer)

	_workerTwo := testWorker()
	_workerTwo.SetID(2)
	_workerTwo.SetHostname("worker_1")
	_workerTwo.SetAddress("localhost")
	_workerTwo.SetActive(true)
	_workerTwo.SetLastCheckedIn(older)

	_workerThree := testWorker()
	_workerThree.SetID(3)
	_workerThree.SetHostname("worker_2")
	_workerThree.SetAddress("localhost")
	_workerThree.SetActive(false)
	_workerThree.SetLastCheckedIn(newer)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "status", "last_status_update_at", "running_build_ids", "last_build_started_at", "last_build_finished_at", "last_checked_in", "build_limit"}).
		AddRow(1, "worker_0", "localhost", nil, true, nil, 0, nil, 0, 0, newer, 0).
		AddRow(2, "worker_1", "localhost", nil, true, nil, 0, nil, 0, 0, older, 0).
		AddRow(3, "worker_2", "localhost", nil, false, nil, 0, nil, 0, 0, newer, 0)

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
		name     string
		database *engine
		want     []*library.Worker
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Worker{_workerOne, _workerTwo, _workerThree},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Worker{_workerOne, _workerTwo, _workerThree},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListWorkers(context.TODO(), "all", newer+1, 0)

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
