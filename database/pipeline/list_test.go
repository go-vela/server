// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestPipeline_Engine_ListPipelines(t *testing.T) {
	// setup types
	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetOwner(_owner)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))
	_repoOne.SetPipelineType(constants.PipelineTypeYAML)
	_repoOne.SetTopics([]string{})
	_repoOne.SetInstallID(0)

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("bazey")
	_repoTwo.SetOrg("fooey")
	_repoTwo.SetName("barey")
	_repoTwo.SetFullName("fooey/barey")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_repoTwo.SetPipelineType(constants.PipelineTypeYAML)
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetInstallID(0)

	_pipelineOne := testutils.APIPipeline()
	_pipelineOne.SetID(1)
	_pipelineOne.SetRepo(_repoOne)
	_pipelineOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipelineOne.SetRef("refs/heads/main")
	_pipelineOne.SetType("yaml")
	_pipelineOne.SetVersion("1")
	_pipelineOne.SetData([]byte("foo"))

	_pipelineTwo := testutils.APIPipeline()
	_pipelineTwo.SetID(2)
	_pipelineTwo.SetRepo(_repoTwo)
	_pipelineTwo.SetCommit("a49aaf4afae6431a79239c95247a2b169fd9f067")
	_pipelineTwo.SetRef("refs/heads/main")
	_pipelineTwo.SetType("yaml")
	_pipelineTwo.SetVersion("1")
	_pipelineTwo.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "pipelines"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "commit", "flavor", "platform", "ref", "type", "version", "services", "stages", "steps", "templates", "data"}).
		AddRow(1, 1, "48afb5bdc41ad69bf22588491333f7cf71135163", "", "", "refs/heads/main", "yaml", "1", false, false, false, false, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69}).
		AddRow(2, 2, "a49aaf4afae6431a79239c95247a2b169fd9f067", "", "", "refs/heads/main", "yaml", "1", false, false, false, false, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69})

	_repoRows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "").
		AddRow(2, 1, "bazey", "fooey", "barey", "fooey/barey", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "")

	_userRows := sqlmock.NewRows(
		[]string{"id", "name", "token", "hash", "active", "admin"}).
		AddRow(1, "foo", "bar", "baz", false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "pipelines"`).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" IN ($1,$2)`).WithArgs(1, 2).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Pipeline{_pipelineOne, _pipelineTwo},
		[]*api.User{_owner},
		[]*api.Repo{_repoOne, _repoTwo},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*api.Pipeline
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Pipeline{_pipelineOne, _pipelineTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Pipeline{_pipelineOne, _pipelineTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListPipelines(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListPipelines for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListPipelines for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListPipelines for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
