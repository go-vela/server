// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuild_Engine_CreateBuildIndexes(t *testing.T) {
	// setup types
	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(CreateCreatedIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateEventIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))
	_mock.ExpectExec(CreateStatusIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.CreateBuildIndexes(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CreateBuildIndexes for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateBuildIndexes for %s returned err: %v", test.name, err)
			}
		})
	}
}
