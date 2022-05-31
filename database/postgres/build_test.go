// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
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

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectRepoBuild, 1, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: true,
			want:    nil,
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

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectLastRepoBuild, 1).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: false,
			want:    nil,
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

func TestPostgres_Client_GetLastCommitBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
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
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "48afb5bdc41ad69bf22588491333f7cf71135163", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND "commit" = $2 LIMIT 1`).
		WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163").WillReturnRows(_rows)
	// ensure the mock expects the query for test case 2
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND "commit" = $2 LIMIT 1`).
		WithArgs(1, "48afb5bdc41ad69bf22588491333f7cf71135163").WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetLastCommitBuild("48afb5bdc41ad69bf22588491333f7cf71135163", _repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetLastCommitBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastCommitBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastCommitBuild is %v, want %v", got, test.want)
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

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectLastRepoBuildByBranch, 1, "master").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"},
	).AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: false,
			want:    nil,
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

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectPendingAndRunningBuilds, "").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"created", "full_name", "number", "status"}).
		AddRow(0, "", 1, "").AddRow(0, "", 2, "")

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.BuildQueue
	}{
		{
			failure: false,
			want:    []*library.BuildQueue{_buildOne, _buildTwo},
		},
		{
			failure: true,
			want:    nil,
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
	_mock.ExpectQuery(`INSERT INTO "builds" ("repo_id","pipeline_id","number","parent","event","event_action","status","error","enqueued","created","started","finished","deploy","deploy_payload","clone","source","title","message","commit","sender","author","email","link","branch","ref","base_ref","head_ref","host","runtime","distribution","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31) RETURNING "id"`).
		WithArgs(1, nil, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, AnyArgument{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
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
	_mock.ExpectExec(`UPDATE "builds" SET "repo_id"=$1,"pipeline_id"=$2,"number"=$3,"parent"=$4,"event"=$5,"event_action"=$6,"status"=$7,"error"=$8,"enqueued"=$9,"created"=$10,"started"=$11,"finished"=$12,"deploy"=$13,"deploy_payload"=$14,"clone"=$15,"source"=$16,"title"=$17,"message"=$18,"commit"=$19,"sender"=$20,"author"=$21,"email"=$22,"link"=$23,"branch"=$24,"ref"=$25,"base_ref"=$26,"head_ref"=$27,"host"=$28,"runtime"=$29,"distribution"=$30 WHERE "id" = $31`).
		WithArgs(1, nil, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, AnyArgument{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
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

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteBuild, 1).Statement

	// ensure the mock expects the query
	_mock.ExpectExec(_query.SQL.String()).WillReturnResult(sqlmock.NewResult(1, 1))

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
		PipelineID:   &i64,
		Number:       &i,
		Parent:       &i,
		Event:        &str,
		EventAction:  &str,
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
