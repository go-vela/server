// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
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
	dbExecutable := types.BuildExecutableFromAPI(_bExecutable)

	err := dbExecutable.Compress(0)
	if err != nil {
		t.Errorf("unable to compress build executable: %v", err)
	}

	err = dbExecutable.Encrypt("A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW")
	if err != nil {
		t.Errorf("unable to encrypt build executable: %v", err)
	}

	_rows := testutils.CreateMockRows([]any{*dbExecutable})

	// ensure the mock expects the query
	_mock.ExpectQuery(`DELETE FROM "build_executables" WHERE build_id = $1 RETURNING *`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err = _sqlite.CreateBuildExecutable(context.TODO(), _bExecutable)
	if err != nil {
		t.Errorf("unable to create test build executable for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *api.BuildExecutable
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
