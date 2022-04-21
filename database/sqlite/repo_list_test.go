// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"log"
	"reflect"
	"testing"

	"github.com/go-vela/server/database/sqlite/ddl"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func init() {
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		log.Fatalf("unable to create new sqlite test database: %v", err)
	}

	// create the repo table
	err = _database.Sqlite.Exec(ddl.CreateRepoTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableRepo, err)
	}
}

func TestSqlite_Client_GetRepoList(t *testing.T) {
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
	_repoOne.SetPreviousName("")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetPreviousName("oldName")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range test.want {
			// create the repo in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

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

func TestSqlite_Client_GetOrgRepoList(t *testing.T) {
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
	_repoOne.SetPreviousName("oldName")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetPreviousName("")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
	filters := map[string]string{}
	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range test.want {
			// create the repo in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

		got, err := _database.GetOrgRepoList("foo", filters, 1, 10, "name")

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

func TestSqlite_Client_GetOrgRepoList_NonAdmin(t *testing.T) {
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
	_repoOne.SetPreviousName("")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("private")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetPreviousName("")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoOne},
		},
	}
	filters := map[string]string{}
	filters["visibility"] = "public"
	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range test.want {
			// create the repo in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

		got, err := _database.GetOrgRepoList("foo", filters, 1, 10, "name")

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

func TestSqlite_Client_GetOrgRepoList_LastUpdate(t *testing.T) {
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
	_repoOne.SetPreviousName("")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetPreviousName("oldName")

	_repoThree := testRepo()
	_repoThree.SetID(3)
	_repoThree.SetUserID(1)
	_repoThree.SetHash("baz")
	_repoThree.SetOrg("foo")
	_repoThree.SetName("bat")
	_repoThree.SetFullName("foo/bat")
	_repoThree.SetVisibility("public")
	_repoThree.SetPipelineType("yaml")
	_repoThree.SetPreviousName("")

	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetCreated(1)
	_buildOne.SetNumber(1)
	_buildOne.SetRepoID(2)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetCreated(2)
	_buildTwo.SetNumber(1)
	_buildTwo.SetRepoID(1)

	_buildThree := testBuild()
	_buildThree.SetID(3)
	_buildThree.SetCreated(3)
	_buildThree.SetNumber(1)
	_buildThree.SetRepoID(3)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	err = _database.CreateBuild(_buildOne)
	if err != nil {
		t.Errorf("unable to create build: %v", err)
	}

	err = _database.CreateBuild(_buildTwo)
	if err != nil {
		t.Errorf("unable to create build: %v", err)
	}

	err = _database.CreateBuild(_buildThree)
	if err != nil {
		t.Errorf("unable to create build: %v", err)
	}

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoThree, _repoOne, _repoTwo},
		},
	}

	filters := map[string]string{}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range test.want {
			// create the repo in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

		got, err := _database.GetOrgRepoList("foo", filters, 1, 10, "latest")

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

func TestSqlite_Client_GetUserRepoList(t *testing.T) {
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
	_repoOne.SetPreviousName("")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetPreviousName("")

	_user := new(library.User)
	_user.SetID(1)
	_user.SetName("foo")
	_user.SetToken("bar")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Repo
	}{
		{
			failure: false,
			want:    []*library.Repo{_repoTwo, _repoOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range test.want {
			// create the repo in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

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
