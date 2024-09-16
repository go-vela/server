// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestHook_Engine_GetHookForRepo(t *testing.T) {
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
	_build.SetDeployNumber(0)
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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "event_action", "branch", "error", "status", "link", "webhook_id"}).
		AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 1)

	_buildRows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_number", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	_repoRows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "")

	_userRows := sqlmock.NewRows(
		[]string{"id", "name", "token", "hash", "active", "admin"}).
		AddRow(1, "foo", "bar", "baz", false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks" WHERE repo_id = $1 AND number = $2 LIMIT $3`).WithArgs(1, 1, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" = $1`).WithArgs(1).WillReturnRows(_buildRows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Hook{_hook},
		[]*api.User{_owner},
		[]*api.Repo{_repo},
		[]*api.Build{_build},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
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
			got, err := test.database.GetHookForRepo(context.TODO(), _repo, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetHookForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetHookForRepo for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetHookForRepo for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
