// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestExecutable_Engine_PopBuildExecutable(t *testing.T) {
	// setup types
	_bExecutable := testBuildExecutable()
	_bExecutable.SetID(1)
	_bExecutable.SetBuildID(1)
	_bExecutable.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "data"}).
		AddRow(1, 1, "+//18dbf7mF+v7ZPK3Wo5h2TD6v4Zg95sCMUJYO2tpwY37DEgTxW5xdyt3Tey9w=")

	// ensure the mock expects the query
	_mock.ExpectQuery(`DELETE FROM "build_executables" WHERE build_id = $1 RETURNING *`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuildExecutable(context.TODO(), _bExecutable)
	if err != nil {
		t.Errorf("unable to create test build executable for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.BuildExecutable
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _bExecutable,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _bExecutable,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.PopBuildExecutable(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("PopBuildExecutable for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("PopBuildExecutable for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("PopBuildExecutable for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
