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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "event_action", "branch", "error", "status", "link", "webhook_id"}).
		AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 1).
		AddRow(2, 1, 0, 2, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 1)

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
	_mock.ExpectQuery(`SELECT * FROM "hooks"`).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE "builds"."id" IN ($1,$2)`).WithArgs(1, 0).WillReturnRows(_buildRows)
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
