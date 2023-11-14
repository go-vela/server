// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListPendingAndRunningBuildsForRepo(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("running")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("pending")
	_buildTwo.SetCreated(1)
	_buildTwo.SetDeployPayload(nil)

	_buildThree := testBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepoID(2)
	_buildThree.SetNumber(1)
	_buildThree.SetStatus("pending")
	_buildThree.SetCreated(1)
	_buildThree.SetDeployPayload(nil)

	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("bazzy")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(2, 1, nil, 2, 0, "", "", "pending", "", 0, 1, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(1, 1, nil, 1, 0, "", "", "running", "", 0, 1, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND (status = 'running' OR status = 'pending')`).WithArgs(1).WillReturnRows(_rows)

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

	_, err = _sqlite.CreateBuild(context.TODO(), _buildThree)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&database.Repo{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(database.RepoFromLibrary(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(database.RepoFromLibrary(_repoTwo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
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
			want:     []*library.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListPendingAndRunningBuildsForRepo(context.TODO(), _repo)

			if test.failure {
				if err == nil {
					t.Errorf("ListPendingAndRunningBuildsForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListPendingAndRunningBuildsForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListPendingAndRunningBuildsForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
