// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestHook_Engine_ListHooks(t *testing.T) {
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

	_hookOne := testutils.APIHook()
	_hookOne.SetID(1)
	_hookOne.SetRepo(_repo)
	_hookOne.SetBuild(_build)
	_hookOne.SetNumber(1)
	_hookOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hookOne.SetWebhookID(1)

	_hookTwo := testutils.APIHook()
	_hookTwo.SetID(2)
	_hookTwo.SetRepo(_repo)
	_hookTwo.SetNumber(2)
	_hookTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hookTwo.SetWebhookID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.HookFromAPI(_hookOne), *types.HookFromAPI(_hookTwo)})

	_buildRows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_build)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks"`).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" = $1`).WithArgs(1).WillReturnRows(_buildRows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Hook{_hookOne, _hookTwo},
		[]*api.User{_owner},
		[]*api.Repo{_repo},
		[]*api.Build{_build},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*api.Hook
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Hook{_hookOne, _hookTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Hook{_hookOne, _hookTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListHooks(context.TODO())

			// empty values of build are different in testing but still empty
			// TODO: fix complex types such as deploy payload, dashboards, favorites, etc for empty comps
			for i, gotHook := range got {
				if gotHook.GetBuild().GetID() == 0 {
					gotHook.SetBuild(test.want[i].GetBuild())
				}
			}

			if test.failure {
				if err == nil {
					t.Errorf("ListHooks for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListHooks for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListHooks for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
