// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestHook_Engine_UpdateHook(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_build := testutils.APIBuild()
	_build.SetID(1)

	_hook := testutils.APIHook()
	_hook.SetID(1)
	_hook.SetRepo(_repo)
	_hook.SetBuild(_build)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "hooks"
SET "repo_id"=$1,"build_id"=$2,"number"=$3,"source_id"=$4,"created"=$5,"host"=$6,"event"=$7,"event_action"=$8,"branch"=$9,"error"=$10,"status"=$11,"link"=$12,"webhook_id"=$13
WHERE "id" = $14`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, nil, nil, nil, nil, nil, nil, nil, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(t, _sqlite, []*api.Hook{_hook}, nil, nil, nil)

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
			_, err := test.database.UpdateHook(context.TODO(), _hook)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateHook for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateHook for %s returned err: %v", test.name, err)
			}
		})
	}
}
