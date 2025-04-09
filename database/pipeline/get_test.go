// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestPipeline_Engine_GetPipeline(t *testing.T) {
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

	_pipeline := testutils.APIPipeline()
	_pipeline.SetID(1)
	_pipeline.SetRepo(_repo)
	_pipeline.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_pipeline.SetRef("refs/heads/main")
	_pipeline.SetType("yaml")
	_pipeline.SetVersion("1")
	_pipeline.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	dbPipeline := types.PipelineFromAPI(_pipeline)

	err := dbPipeline.Compress(0)
	if err != nil {
		t.Errorf("unable to compress pipeline: %v", err)
	}

	_rows := testutils.CreateMockRows([]any{*dbPipeline})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "pipelines" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Pipeline{_pipeline},
		[]*api.User{_owner},
		[]*api.Repo{_repo},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Pipeline
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _pipeline,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _pipeline,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetPipeline(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetPipeline for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetPipeline for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetPipeline for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
