// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExecutable_Engine_CreateBuildExecutable(t *testing.T) {
	// setup types
	_bExecutable := testBuildExecutable()
	_bExecutable.SetID(1)
	_bExecutable.SetBuildID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "build_executables"
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
			err := test.database.CreateBuildExecutable(context.TODO(), _bExecutable)

			if test.failure {
				if err == nil {
					t.Errorf("CreateBuildExecutable for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateBuildExecutable for %s returned err: %v", test.name, err)
			}
		})
	}
}
