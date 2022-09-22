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
)

func init() {
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		log.Fatalf("unable to create new sqlite test database: %v", err)
	}

	// create the build table
	err = _database.Sqlite.Exec(ddl.CreateBuildTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableBuild, err)
	}
}

func TestSqlite_Client_GetBuildCount(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

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
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the builds in the database
		err := _database.CreateBuild(_buildOne)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.CreateBuild(_buildTwo)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		got, err := _database.GetBuildCount()

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetBuildCountByStatus(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("running")
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("running")
	_buildTwo.SetDeployPayload(nil)

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
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the builds in the database
		err := _database.CreateBuild(_buildOne)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.CreateBuild(_buildTwo)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		got, err := _database.GetBuildCountByStatus("running")

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildCountByStatus should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildCountByStatus returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildCountByStatus is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgBuildCount(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

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
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		_, err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the builds in the database
		err = _database.CreateBuild(_buildOne)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.CreateBuild(_buildTwo)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		got, err := _database.GetOrgBuildCount("foo", filters)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgBuildCountByEvent(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetEvent("push")
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetEvent("push")
	_buildTwo.SetDeployPayload(nil)

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
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	filters := map[string]interface{}{
		"event": "push",
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		_, err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the builds in the database
		err = _database.CreateBuild(_buildOne)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.CreateBuild(_buildTwo)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		got, err := _database.GetOrgBuildCount("foo", filters)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildCountByEvent should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildCountByEvent returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildCountByEvent is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetRepoBuildCount(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetDeployPayload(nil)

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
		want    int64
	}{
		{
			failure: false,
			want:    2,
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		_, err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the builds in the database
		err = _database.CreateBuild(_buildOne)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.CreateBuild(_buildTwo)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		got, err := _database.GetRepoBuildCount(_repo, filters)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoBuildCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoBuildCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoBuildCount is %v, want %v", got, test.want)
		}
	}
}
