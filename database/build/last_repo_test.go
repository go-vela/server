// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestBuild_Engine_LastBuildForRepo(t *testing.T) {
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

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(_repo)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)
	_build.SetBranch("main")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_build)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND branch = $2 ORDER BY number DESC LIMIT $3`).WithArgs(1, "main", 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _build)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _build,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _build,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.LastBuildForRepo(context.TODO(), _repo, "main")

			if test.failure {
				if err == nil {
					t.Errorf("LastBuildForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("LastBuildForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("LastBuildForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
