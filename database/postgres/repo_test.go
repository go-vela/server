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

func TestPostgres_Client_GetRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetPreviousName("")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectRepo, "foo", "bar").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type", "previous_name"},
	).AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, 0, "public", false, false, false, false, false, false, false, false, "yaml", "")

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Repo
	}{
		{
			failure: false,
			want:    _repo,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetRepo("foo", "bar")

		if test.failure {
			if err == nil {
				t.Errorf("GetRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepo returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepo is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetPreviousName("oldName")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "repos" ("user_id","hash","org","name","full_name","link","clone","branch","build_limit","timeout","counter","visibility","private","trusted","active","allow_pull","allow_push","allow_deploy","allow_tag","allow_comment","pipeline_type","previous_name","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23) RETURNING "id"`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, false, false, false, false, false, "yaml", "oldName", 1).
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
		err := _database.CreateRepo(_repo)

		if test.failure {
			if err == nil {
				t.Errorf("CreateRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateRepo returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetPreviousName("oldName")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "repos" SET "user_id"=$1,"hash"=$2,"org"=$3,"name"=$4,"full_name"=$5,"link"=$6,"clone"=$7,"branch"=$8,"build_limit"=$9,"timeout"=$10,"counter"=$11,"visibility"=$12,"private"=$13,"trusted"=$14,"active"=$15,"allow_pull"=$16,"allow_push"=$17,"allow_deploy"=$18,"allow_tag"=$19,"allow_comment"=$20,"pipeline_type"=$21,"previous_name"=$22 WHERE "id" = $23`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, false, false, false, false, false, "yaml", "oldName", 1).
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
		err := _database.UpdateRepo(_repo)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateRepo returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteRepo(t *testing.T) {
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
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteRepo, 1).Statement

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
		err := _database.DeleteRepo(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteRepo should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteRepo returned err: %v", err)
		}
	}
}

// testRepo is a test helper function to create a
// library Repo type with all fields set to their
// zero values.
func testRepo() *library.Repo {
	i64 := int64(0)
	i := 0
	str := ""
	b := false

	return &library.Repo{
		ID:           &i64,
		PipelineType: &str,
		UserID:       &i64,
		Hash:         &str,
		Org:          &str,
		Name:         &str,
		FullName:     &str,
		Link:         &str,
		Clone:        &str,
		Branch:       &str,
		BuildLimit:   &i64,
		Timeout:      &i64,
		Counter:      &i,
		Visibility:   &str,
		Private:      &b,
		Trusted:      &b,
		Active:       &b,
		AllowPull:    &b,
		AllowPush:    &b,
		AllowDeploy:  &b,
		AllowTag:     &b,
		AllowComment: &b,
		PreviousName: &str,
	}
}
