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

func TestHook_Engine_GetHookByWebhookID(t *testing.T) {
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
	_hook.SetWebhookID(123456)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.HookFromAPI(_hook)})

	_buildRows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_build)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks" WHERE webhook_id = $1 LIMIT $2`).WithArgs(123456, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" = $1`).WithArgs(1).WillReturnRows(_buildRows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(t, _sqlite, []*api.Hook{_hook}, []*api.User{_owner}, []*api.Repo{_repo}, []*api.Build{_build})

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
			got, err := test.database.GetHookByWebhookID(context.TODO(), 123456)

			if test.failure {
				if err == nil {
					t.Errorf("GetHookByWebhookID for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetHookByWebhookID for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetHookByWebhookID for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
