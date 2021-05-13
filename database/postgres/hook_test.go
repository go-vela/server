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

func TestPostgres_Client_GetHook(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "branch", "error", "status", "link"},
	).AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectRepoHook).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Hook
	}{
		{
			failure: false,
			want:    _hook,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetHook(1, _repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetHook should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetHook returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetHook is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetLastHook(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "branch", "error", "status", "link"},
	).AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectLastRepoHook).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Hook
	}{
		{
			failure: false,
			want:    _hook,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetLastHook(_repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetLastHook should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastHook returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastHook is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "hooks" ("repo_id","build_id","number","source_id","created","host","event","branch","error","status","link","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, "", "", "", "", "", "", 1).
		WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		err := _database.CreateHook(_hook)

		if test.failure {
			if err == nil {
				t.Errorf("CreateHook should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateHook returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "hooks" SET "repo_id"=$1,"build_id"=$2,"number"=$3,"source_id"=$4,"created"=$5,"host"=$6,"event"=$7,"branch"=$8,"error"=$9,"status"=$10,"link"=$11 WHERE "id" = $12`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, "", "", "", "", "", "", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		err := _database.UpdateHook(_hook)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateHook should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateHook returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteHook(t *testing.T) {
	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// ensure the mock expects the query
	_mock.ExpectExec(dml.DeleteHook).WillReturnResult(sqlmock.NewResult(1, 1))

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		err := _database.DeleteHook(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteHook should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteHook returned err: %v", err)
		}
	}
}

// testHook is a test helper function to create a
// library Hook type with all fields set to their
// zero values.
func testHook() *library.Hook {
	i := 0
	i64 := int64(0)
	str := ""

	return &library.Hook{
		ID:       &i64,
		RepoID:   &i64,
		BuildID:  &i64,
		Number:   &i,
		SourceID: &str,
		Created:  &i64,
		Host:     &str,
		Event:    &str,
		Branch:   &str,
		Error:    &str,
		Status:   &str,
		Link:     &str,
	}
}
