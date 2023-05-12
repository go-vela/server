// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiled

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCompiled_Engine_CreateCompiled(t *testing.T) {
	// setup types
	_compiled := testCompiled()
	_compiled.SetID(1)
	_compiled.SetBuildID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "compiled"
("build_id","data","id")
VALUES ($1,$2,$3) RETURNING "id"`).
		WithArgs(1, AnyArgument{}, 1).
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
			err := test.database.CreateCompiled(_compiled)

			if test.failure {
				if err == nil {
					t.Errorf("CreateCompiled for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateCompiled for %s returned err: %v", test.name, err)
			}
		})
	}
}
