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

func TestPostgres_Client_GetRepoList(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")

	_repoTwo := testRepo()
	_repoTwo.SetID(1)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListRepos).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type"},
	).AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml").
		AddRow(1, 1, "baz", "bar", "foo", "bar/foo", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoOne, _repoTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetRepoList()

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetOrgRepoList(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")

	_repoTwo := testRepo()
	_repoTwo.SetID(1)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListOrgRepos, "foo", []string{""}, 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type"},
	).AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml").
		AddRow(1, 1, "baz", "foo", "baz", "foo/baz", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoOne, _repoTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetOrgRepoList("foo", []string{""}, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgRepoList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgRepoList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgRepoList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetOrgPrivateRepoList(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("private")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("private")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListPrivateOrgRepos, "foo").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment"},
	).AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, "private", false, false, false, false, false, false, false, false).
		AddRow(2, 1, "baz", "foo", "baz", "foo/baz", "", "", "", 0, 0, "private", false, false, false, false, false, false, false, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoOne, _repoTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetOrgPrivateRepoList("foo")

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgRepoList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgRepoList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgRepoList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetUserRepoList(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")

	_repoTwo := testRepo()
	_repoTwo.SetID(1)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")

	_user := new(library.User)
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListUserRepos, 1, 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_pull", "allow_push", "allow_deploy", "allow_tag", "allow_comment", "pipeline_type"},
	).AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml").
		AddRow(1, 1, "baz", "bar", "foo", "bar/foo", "", "", "", 0, 0, "public", false, false, false, false, false, false, false, false, "yaml")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoOne, _repoTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetUserRepoList(_user, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetUserRepoList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUserRepoList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUserRepoList is %v, want %v", got, test.want)
		}
	}
}
