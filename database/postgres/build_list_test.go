// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetBuildList(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListBuilds).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 1, 2, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuildList()

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetOrgBuildList(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT count(*) FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1").WillReturnRows(_rows)

	// create expected return in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 1, 2, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT builds.* FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1 ORDER BY created DESC,id LIMIT 10").WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne, _buildTwo},
		},
	}
	filters := map[string]string{}
	// run tests
	for _, test := range tests {
		got, _, err := _database.GetOrgBuildList("foo", filters, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetOrgBuildList_NonAdmin(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT count(*) FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1 WHERE \"visibility\" = $2").WillReturnRows(_rows)

	// create expected return in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT builds.* FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1 WHERE \"visibility\" = $2 ORDER BY created DESC,id LIMIT 10").WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne},
		},
	}
	filters := map[string]string{}
	filters["visibility"] = "public"
	// run tests
	for _, test := range tests {
		got, _, err := _database.GetOrgBuildList("foo", filters, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetOrgBuildListByEvent(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT count(*) FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1 WHERE \"event\" = $2").WillReturnRows(_rows)

	// create expected return in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 1, 2, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT builds.* FROM \"builds\" JOIN repos ON builds.repo_id = repos.id and repos.org = $1 WHERE \"event\" = $2 ORDER BY created DESC,id LIMIT 10").WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne, _buildTwo},
		},
	}
	filters := map[string]string{}
	filters["event"] = "push"
	// run tests
	for _, test := range tests {
		got, _, err := _database.GetOrgBuildList("foo", filters, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildListByEvent should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildListByEvent returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildListByEvent is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetRepoBuildList(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "builds" WHERE repo_id = $1`).WillReturnRows(_rows)

	// create expected return in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 1, 2, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 ORDER BY number DESC LIMIT 10`).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne, _buildTwo},
		},
	}

	filters := map[string]string{}

	// run tests
	for _, test := range tests {
		got, _, err := _database.GetRepoBuildList(_repo, filters, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoBuildList is %v, want %v", got, test.want)
		}
	}
}
