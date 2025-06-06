// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPipeline_Engine_CreatePipelineIndexes(t *testing.T) {
	// setup types
	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(CreateRepoIDIndex).WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.CreatePipelineIndexes(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CreatePipelineIndexes for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreatePipelineIndexes for %s returned err: %v", test.name, err)
			}
		})
	}
}
