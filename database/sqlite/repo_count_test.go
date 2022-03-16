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

func TestSqlite_Client_GetRepoCount(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repos in the database
		err := _database.CreateRepo(_repoOne)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		err = _database.CreateRepo(_repoTwo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		got, err := _database.GetRepoCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetUserRepoCount(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")

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
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repos in the database
		err := _database.CreateRepo(_repoOne)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		err = _database.CreateRepo(_repoTwo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		got, err := _database.GetUserRepoCount(_user)

		if test.failure {
			if err == nil {
				t.Errorf("GetUserRepoCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetUserRepoCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetUserRepoCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgRepoCount(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    int64
	}{
		{
			failure: false,
			want:    1,
		},
	}
	filters := map[string]string{}
	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repos in the database
		err := _database.CreateRepo(_repoOne)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		err = _database.CreateRepo(_repoTwo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		got, err := _database.GetOrgRepoCount("foo", filters)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgRepoCount_NonAdmin(t *testing.T) {
	// setup types
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetUserID(1)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetUserID(1)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("foo/foo")
	_repoTwo.SetVisibility("private")

	_repoThree := testRepo()
	_repoThree.SetID(3)
	_repoThree.SetUserID(1)
	_repoThree.SetHash("baz")
	_repoThree.SetOrg("bar")
	_repoThree.SetName("foo")
	_repoThree.SetFullName("bar/foo")
	_repoThree.SetVisibility("private")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    int64
	}{
		{
			failure: false,
			want:    1,
		},
	}
	filters := map[string]string{}
	filters["visibility"] = "public"
	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		for _, repo := range []*library.Repo{_repoOne, _repoTwo, _repoThree} {
			// create the repos in the database
			err := _database.CreateRepo(repo)
			if err != nil {
				t.Errorf("unable to create test repo: %v", err)
			}
		}

		got, err := _database.GetOrgRepoCount("foo", filters)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoCount is %v, want %v", got, test.want)
		}
	}
}
