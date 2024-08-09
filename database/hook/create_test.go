// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestHook_Engine_CreateHook(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "hooks"
("repo_id","build_id","number","source_id","created","host","event","event_action","branch","error","status","link","webhook_id","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, nil, nil, nil, nil, nil, nil, nil, 1, 1).
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
			_, err := test.database.CreateHook(context.TODO(), _hook)

			if test.failure {
				if err == nil {
					t.Errorf("CreateHook for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateHook for %s returned err: %v", test.name, err)
			}
		})
	}
}
