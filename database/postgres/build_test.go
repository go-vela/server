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
)

func TestPostgres_Client_GetBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectRepoBuild).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetBuild(1, _repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuild is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetLastBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectLastRepoBuild).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetLastBuild(_repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetLastBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastBuild is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetLastBuildByBranch(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "parent", "event", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, 1, 0, "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectLastRepoBuildByBranch).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetLastBuildByBranch(_repo, "master")

		if test.failure {
			if err == nil {
				t.Errorf("GetLastBuildByBranch should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastBuildByBranch returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastBuildByBranch is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetPendingAndRunningBuilds(t *testing.T) {
	// setup types
	_buildOne := new(library.BuildQueue)
	_buildOne.SetCreated(0)
	_buildOne.SetFullName("")
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("")

	_buildTwo := new(library.BuildQueue)
	_buildTwo.SetCreated(0)
	_buildTwo.SetFullName("")
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"created", "full_name", "number", "status"}).
		AddRow(0, "", 1, "").
		AddRow(0, "", 2, "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectPendingAndRunningBuilds).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.BuildQueue
	}{
		{
			failure: false,
			want:    []*library.BuildQueue{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetPendingAndRunningBuilds("")

		if test.failure {
			if err == nil {
				t.Errorf("GetPendingAndRunningBuilds should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetPendingAndRunningBuilds returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetPendingAndRunningBuilds is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

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
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "builds" ("repo_id","number","parent","event","status","error","enqueued","created","started","finished","deploy","deploy_payload","clone","source","title","message","commit","sender","author","email","link","branch","ref","base_ref","head_ref","host","runtime","distribution","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29) RETURNING "id"`).
		WithArgs(1, 1, nil, "", "", "", nil, nil, nil, nil, "", AnyArgument{}, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 1).
		WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.CreateBuild(_build)

		if test.failure {
			if err == nil {
				t.Errorf("CreateBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateBuild returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

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

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "builds" SET "repo_id"=$1,"number"=$2,"parent"=$3,"event"=$4,"status"=$5,"error"=$6,"enqueued"=$7,"created"=$8,"started"=$9,"finished"=$10,"deploy"=$11,"deploy_payload"=$12,"clone"=$13,"source"=$14,"title"=$15,"message"=$16,"commit"=$17,"sender"=$18,"author"=$19,"email"=$20,"link"=$21,"branch"=$22,"ref"=$23,"base_ref"=$24,"head_ref"=$25,"host"=$26,"runtime"=$27,"distribution"=$28 WHERE "id" = $29`).
		WithArgs(1, 1, nil, "", "", "", nil, nil, nil, nil, "", AnyArgument{}, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.UpdateBuild(_build)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateBuild returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteBuild(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(dml.DeleteBuild).WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.DeleteBuild(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteBuild returned err: %v", err)
		}
	}
}

// testBuild is a test helper function to create a
// library Build type with all fields set to their
// zero values.
func testBuild() *library.Build {
	i64 := int64(0)
	i := 0
	str := ""

	return &library.Build{
		ID:           &i64,
		RepoID:       &i64,
		Number:       &i,
		Parent:       &i,
		Event:        &str,
		Status:       &str,
		Error:        &str,
		Enqueued:     &i64,
		Created:      &i64,
		Started:      &i64,
		Finished:     &i64,
		Deploy:       &str,
		Clone:        &str,
		Source:       &str,
		Title:        &str,
		Message:      &str,
		Commit:       &str,
		Sender:       &str,
		Author:       &str,
		Email:        &str,
		Link:         &str,
		Branch:       &str,
		Ref:          &str,
		BaseRef:      &str,
		HeadRef:      &str,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
