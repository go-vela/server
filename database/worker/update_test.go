// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWorker_Engine_UpdateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "workers"
SET "hostname"=$1,"address"=$2,"routes"=$3,"active"=$4,"status"=$5,"last_status_update_at"=$6,"running_build_ids"=$7,"last_build_started_at"=$8,"last_build_finished_at"=$9,"last_checked_in"=$10,"build_limit"=$11
WHERE "id" = $12`).
		WithArgs("worker_0", "localhost", nil, true, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateWorker(_worker)
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
			err = test.database.UpdateWorker(_worker)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateWorker for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateWorker for %s returned err: %v", test.name, err)
			}
		})
	}
}
