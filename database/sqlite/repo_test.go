// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

func TestSqlite_Client_GetRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Repo
	}{
		{
			failure: false,
			want:    _repo,
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		err := _database.CreateRepo(test.want)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

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

func TestSqlite_Client_CreateRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

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

func TestSqlite_Client_UpdateRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		err = _database.UpdateRepo(_repo)

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

func TestSqlite_Client_DeleteRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

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
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		err = _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

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
	str := ""
	b := false

	return &library.Repo{
		ID:           &i64,
		UserID:       &i64,
		Hash:         &str,
		Org:          &str,
		Name:         &str,
		FullName:     &str,
		Link:         &str,
		Clone:        &str,
		Branch:       &str,
		Timeout:      &i64,
		Visibility:   &str,
		Private:      &b,
		Trusted:      &b,
		Active:       &b,
		AllowPull:    &b,
		AllowPush:    &b,
		AllowDeploy:  &b,
		AllowTag:     &b,
		AllowComment: &b,
	}
}
