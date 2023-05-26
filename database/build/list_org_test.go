// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

func TestBuild_Engine_ListBuildsForOrg(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)
	_buildOne.SetEvent("push")

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(2)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)
	_buildTwo.SetEvent("push")

	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")
	_repoOne.SetTopics([]string{})

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("bar")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected count query without filters result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
	// ensure the mock expects the count query without filters
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1`).WithArgs("foo").WillReturnRows(_rows)
	// create expected query without filters result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 2, nil, 2, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)
	// ensure the mock expects the query without filters
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 ORDER BY created DESC,id LIMIT 10`).WithArgs("foo").WillReturnRows(_rows)

	// create expected count query with event filter result in mock
	_rows = sqlmock.NewRows([]string{"count"}).AddRow(2)
	// ensure the mock expects the count query with event filter
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND "event" = $2`).WithArgs("foo", "push").WillReturnRows(_rows)
	// create expected query with event filter result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 2, nil, 2, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)
	// ensure the mock expects the query with event filter
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND "event" = $2 ORDER BY created DESC,id LIMIT 10`).WithArgs("foo", "push").WillReturnRows(_rows)

	// create expected count query with visibility filter result in mock
	_rows = sqlmock.NewRows([]string{"count"}).AddRow(2)
	// ensure the mock expects the count query with visibility filter
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND "visibility" = $2`).WithArgs("foo", "public").WillReturnRows(_rows)
	// create expected query with visibility filter result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 2, nil, 2, 0, "push", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)
	// ensure the mock expects the query with visibility filter
	_mock.ExpectQuery(`SELECT builds.* FROM "builds" JOIN repos ON builds.repo_id = repos.id WHERE repos.org = $1 AND "visibility" = $2 ORDER BY created DESC,id LIMIT 10`).WithArgs("foo", "public").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuild(_buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.CreateBuild(_buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&database.Repo{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(database.RepoFromLibrary(_repoOne)).Error
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
		filters  map[string]interface{}
		want     []*library.Build
	}{
		{
			failure:  false,
			name:     "postgres without filters",
			database: _postgres,
			filters:  map[string]interface{}{},
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "postgres with event filter",
			database: _postgres,
			filters: map[string]interface{}{
				"event": "push",
			},
			want: []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "postgres with visibility filter",
			database: _postgres,
			filters: map[string]interface{}{
				"visibility": "public",
			},
			want: []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 without filters",
			database: _sqlite,
			filters:  map[string]interface{}{},
			want:     []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 with event filter",
			database: _sqlite,
			filters: map[string]interface{}{
				"event": "push",
			},
			want: []*library.Build{_buildOne, _buildTwo},
		},
		{
			failure:  false,
			name:     "sqlite3 with visibility filter",
			database: _sqlite,
			filters: map[string]interface{}{
				"visibility": "public",
			},
			want: []*library.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListBuildsForOrg("foo", test.filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListBuildsForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListBuildsForOrg for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListBuildsForOrg for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
