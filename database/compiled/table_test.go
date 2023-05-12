// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCompiled_Engine_CreateCompiledTable(t *testing.T) {
	// setup types
	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))

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
			err := test.database.CreateCompiledTable(test.name)

			if test.failure {
				if err == nil {
					t.Errorf("CreateCompiledTable for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateCompiledTable for %s returned err: %v", test.name, err)
			}
		})
	}
}
