// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/types/library"
)

func TestBuild_Engine_GetBuildForRepo(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployNumber(0)
	_build.SetDeployPayload(nil)

	_repo := testRepo()
	_repo.SetID(1)
	_repo.GetOwner().SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_number", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND number = $2 LIMIT $3`).WithArgs(1, 1, 1).WillReturnRows(_rows)

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
		database *engine
		want     *library.Build
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
			got, err := test.database.GetBuildForRepo(context.TODO(), _repo, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetBuildForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetBuildForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetBuildForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
