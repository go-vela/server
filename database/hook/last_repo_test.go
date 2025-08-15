// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestHook_Engine_LastHookForRepo(t *testing.T) {
	// setup types
	_owner := testutils.APIUser().Crop()
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
	_repo.SetAllowEvents(api.NewEventsFromMask(1))
	_repo.SetPipelineType(constants.PipelineTypeYAML)
	_repo.SetTopics([]string{})

	_repoBuild := new(api.Repo)
	_repoBuild.SetID(1)

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(_repoBuild)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
	_rows := testutils.CreateMockRows([]any{*types.HookFromAPI(_hook)})

	_buildRows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_build)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks" WHERE repo_id = $1 ORDER BY number DESC LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" = $1`).WithArgs(1).WillReturnRows(_buildRows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(t, _sqlite, []*api.Hook{_hook}, []*api.User{}, []*api.Repo{}, []*api.Build{_build})

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Hook
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _hook,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _hook,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.LastHookForRepo(context.TODO(), _repo)

			if test.failure {
				if err == nil {
					t.Errorf("LastHookForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("LastHookForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LastHookForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
