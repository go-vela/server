// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package worker

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWorker_Engine_CreateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "workers"
("hostname","address","routes","active","status","last_status_update_at","running_build_ids","last_build_finished_at","last_checked_in","build_limit","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`).
		WithArgs("worker_0", "localhost", nil, true, nil, nil, nil, nil, nil, nil, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			err := test.database.CreateWorker(_worker)

			if test.failure {
				if err == nil {
					t.Errorf("CreateWorker for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateWorker for %s returned err: %v", test.name, err)
			}
		})
	}
}
