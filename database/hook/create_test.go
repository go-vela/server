// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestHook_Engine_CreateHook(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetOwner(_owner)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

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

	_counterRow := sqlmock.NewRows([]string{"hook_counter"}).AddRow(1)

	_mock.ExpectBegin()

	_mock.ExpectQuery(`UPDATE repos SET hook_counter = hook_counter + 1 WHERE id = $1 RETURNING hook_counter`).WithArgs(1).WillReturnRows(_counterRow)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "hooks"
("repo_id","build_id","number","source_id","created","host","event","event_action","branch","error","status","link","webhook_id","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, nil, nil, nil, nil, nil, nil, nil, 1, 1).
		WillReturnRows(_rows)

	_mock.ExpectCommit()

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Hook{},
		[]*api.User{_owner},
		[]*api.Repo{_repo},
		[]*api.Build{},
	)

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

// createTestHook is a helper function to create hook objects without the repo update transaction.
// used only for unit tests.
func createTestHook(_ context.Context, e *Engine, h *api.Hook) error {
	err := e.client.Table(constants.TableHook).Create(types.HookFromAPI(h)).Error
	if err != nil {
		return fmt.Errorf("unable to create test hook: %w", err)
	}

	return nil
}
