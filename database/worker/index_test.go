// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWorker_Engine_CreateWorkerIndexes(t *testing.T) {
	// setup types
	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(CreateHostnameAddressIndex).WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
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
			err := test.database.CreateWorkerIndexes(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CreateWorkerIndexes for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateWorkerIndexes for %s returned err: %v", test.name, err)
			}
		})
	}
}
