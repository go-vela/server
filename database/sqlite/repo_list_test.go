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
	_repoOne.SetLastUpdate(1)

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
	_repoTwo.SetLastUpdate(2)

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
	_repoOne.SetLastUpdate(1)

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
	_repoTwo.SetLastUpdate(2)

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

		got, err := _database.GetOrgRepoList("foo", filters, 1, 10)

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
	_repoOne.SetLastUpdate(1)

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
	_repoTwo.SetLastUpdate(2)

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

		got, err := _database.GetOrgRepoList("foo", filters, 1, 10)

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
	_repoOne.SetLastUpdate(1)

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
	_repoTwo.SetLastUpdate(2)

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
