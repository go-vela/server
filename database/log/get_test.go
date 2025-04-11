// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestLog_Engine_GetLog(t *testing.T) {
	// setup types
	_log := testutils.APILog()
	_log.SetID(1)
	_log.SetRepoID(1)
	_log.SetBuildID(1)
	_log.SetStepID(1)
	_log.SetData([]byte{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.LogFromAPI(_log)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "logs" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(context.TODO(), _log)
	if err != nil {
		t.Errorf("unable to create test log for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _log,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _log,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetLog(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetLog for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetLog for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetLog for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
