// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWorker_Engine_DeleteWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "workers" WHERE "workers"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.database.DeleteWorker(context.TODO(), _worker)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteWorker for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteWorker for %s returned err: %v", test.name, err)
			}
		})
	}
}
