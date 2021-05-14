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

func TestPostgres_Client_GetHookList(t *testing.T) {
	// setup types
	_hookOne := testHook()
	_hookOne.SetID(1)
	_hookOne.SetRepoID(1)
	_hookOne.SetBuildID(1)
	_hookOne.SetNumber(1)
	_hookOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	_hookTwo := testHook()
	_hookTwo.SetID(2)
	_hookTwo.SetRepoID(1)
	_hookTwo.SetBuildID(2)
	_hookTwo.SetNumber(2)
	_hookTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "branch", "error", "status", "link"},
	).AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "").
		AddRow(2, 1, 2, 2, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListHooks).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Hook
	}{
		{
			failure: false,
			want:    []*library.Hook{_hookOne, _hookTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetHookList()

		if test.failure {
			if err == nil {
				t.Errorf("GetHookList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetHookList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetHookList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetRepoHookList(t *testing.T) {
	// setup types
	_hookOne := testHook()
	_hookOne.SetID(1)
	_hookOne.SetRepoID(1)
	_hookOne.SetBuildID(1)
	_hookOne.SetNumber(1)
	_hookOne.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	_hookTwo := testHook()
	_hookTwo.SetID(2)
	_hookTwo.SetRepoID(1)
	_hookTwo.SetBuildID(2)
	_hookTwo.SetNumber(2)
	_hookTwo.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")

	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "branch", "error", "status", "link"},
	).AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "").
		AddRow(2, 1, 2, 2, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListRepoHooks).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Hook
	}{
		{
			failure: false,
			want:    []*library.Hook{_hookOne, _hookTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetRepoHookList(_repo, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoHookList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoHookList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoHookList is %v, want %v", got, test.want)
		}
	}
}
