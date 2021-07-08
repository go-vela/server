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

	// create the build table
	err = _database.Sqlite.Exec(ddl.CreateBuildTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableBuild, err)
	}
}

func TestSqlite_Client_GetBuildList(t *testing.T) {
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
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		for _, build := range test.want {
			// create the build in the database
			err := _database.CreateBuild(build)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, err := _database.GetBuildList()

		if test.failure {
			if err == nil {
				t.Errorf("GetBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgBuildList(t *testing.T) {
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
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildTwo, _buildOne},
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

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		for _, build := range test.want {
			// create the build in the database
			err := _database.CreateBuild(build)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, _, err := _database.GetOrgBuildList("foo", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetOrgBuildListByEvent(t *testing.T) {
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
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildTwo, _buildOne},
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

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		for _, build := range test.want {
			// create the build in the database
			err := _database.CreateBuild(build)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, _, err := _database.GetOrgBuildListByEvent("foo", "push", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetOrgBuildListByEvent should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetOrgBuildListByEvent returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetOrgBuildListByEvent is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetRepoBuildList(t *testing.T) {
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
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildTwo, _buildOne},
		},
	}

	filters := map[string]string{}

	// run tests
	for _, test := range tests {
		// defer cleanup of the repos table
		defer _database.Sqlite.Exec("delete from repos;")

		// create the repo in the database
		err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		for _, build := range test.want {
			// create the build in the database
			err := _database.CreateBuild(build)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, _, err := _database.GetRepoBuildList(_repo, filters, 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoBuildList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoBuildList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoBuildList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetRepoBuildListByEvent(t *testing.T) {
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
		want    []*library.Build
	}{
		{
			failure: false,
			want:    []*library.Build{_buildTwo, _buildOne},
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

		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		for _, build := range test.want {
			// create the build in the database
			err := _database.CreateBuild(build)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, _, err := _database.GetRepoBuildListByEvent(_repo, "push", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoBuildListByEvent should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoBuildListByEvent returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoBuildListByEvent is %v, want %v", got, test.want)
		}
	}
}
