// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
	"github.com/google/go-cmp/cmp"
)

func TestBuild_Engine_ListBuildsForDashboardRepo(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)
	_buildOne.SetCreated(1)
	_buildOne.SetEvent("push")
	_buildOne.SetBranch("main")

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)
	_buildTwo.SetCreated(2)
	_buildTwo.SetEvent("pull_request")
	_buildTwo.SetBranch("main")

	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected query result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "approved_at", "approved_by", "timestamp"}).
		AddRow(2, 1, nil, 2, 0, "pull_request", "", "", "", 0, 2, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "main", "", "", "", "", "", "", 0, "", 0).
		AddRow(1, 1, nil, 1, 0, "push", "", "", "", 0, 1, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "main", "", "", "", "", "", "", 0, "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND branch IN ($2) AND event IN ($3,$4) ORDER BY number DESC LIMIT 5`).WithArgs(1, "main", "push", "pull_request").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Build{_buildTwo, _buildOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Build{_buildTwo, _buildOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListBuildsForDashboardRepo(context.TODO(), _repo, []string{"main"}, []string{"push", "pull_request"})

			if test.failure {
				if err == nil {
					t.Errorf("ListBuildsForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuildsForRepo for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("GetDashboard mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
