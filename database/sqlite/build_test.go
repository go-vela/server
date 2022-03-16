// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the build in the database
			err := _database.CreateBuild(test.want)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, err := _database.GetBuild(1, _repo)

		// cleanup the builds table
		_ = _database.Sqlite.Exec("DELETE FROM builds;")

		if test.failure {
			if err == nil {
				t.Errorf("GetBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetBuild is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetLastBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

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
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: false,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the build in the database
			err := _database.CreateBuild(test.want)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, err := _database.GetLastBuild(_repo)

		// cleanup the builds table
		_ = _database.Sqlite.Exec("DELETE FROM builds;")

		if test.failure {
			if err == nil {
				t.Errorf("GetLastBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastBuild returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastBuild is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetLastBuildByBranch(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetBranch("master")
	_build.SetDeployPayload(nil)

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
		want    *library.Build
	}{
		{
			failure: false,
			want:    _build,
		},
		{
			failure: false,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the build in the database
			err := _database.CreateBuild(test.want)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, err := _database.GetLastBuildByBranch(_repo, "master")

		// cleanup the builds table
		_ = _database.Sqlite.Exec("DELETE FROM builds;")

		if test.failure {
			if err == nil {
				t.Errorf("GetLastBuildByBranch should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetLastBuildByBranch returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetLastBuildByBranch is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetPendingAndRunningBuilds(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("running")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("pending")
	_buildTwo.SetCreated(1)
	_buildTwo.SetDeployPayload(nil)

	_queueOne := new(library.BuildQueue)
	_queueOne.SetCreated(1)
	_queueOne.SetFullName("foo/bar")
	_queueOne.SetNumber(1)
	_queueOne.SetStatus("running")

	_queueTwo := new(library.BuildQueue)
	_queueTwo.SetCreated(1)
	_queueTwo.SetFullName("foo/bar")
	_queueTwo.SetNumber(2)
	_queueTwo.SetStatus("pending")

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
		want    []*library.BuildQueue
	}{
		{
			failure: false,
			want:    []*library.BuildQueue{_queueOne, _queueTwo},
		},
		{
			failure: false,
			want:    []*library.BuildQueue{},
		},
	}

	// run tests
	for _, test := range tests {
		// create the repo in the database
		err := _database.CreateRepo(_repo)
		if err != nil {
			t.Errorf("unable to create test repo: %v", err)
		}

		if len(test.want) > 0 {
			// create the builds in the database
			err = _database.CreateBuild(_buildOne)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}

			err = _database.CreateBuild(_buildTwo)
			if err != nil {
				t.Errorf("unable to create test build: %v", err)
			}
		}

		got, err := _database.GetPendingAndRunningBuilds("0")

		// cleanup the repos table
		_ = _database.Sqlite.Exec("DELETE FROM repos;")
		// cleanup the builds table
		_ = _database.Sqlite.Exec("DELETE FROM builds;")

		if test.failure {
			if err == nil {
				t.Errorf("GetPendingAndRunningBuilds should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetPendingAndRunningBuilds returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetPendingAndRunningBuilds is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_CreateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

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
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		err := _database.CreateBuild(_build)

		if test.failure {
			if err == nil {
				t.Errorf("CreateBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateBuild returned err: %v", err)
		}
	}
}

func TestSqlite_Client_UpdateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

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
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the build in the database
		err = _database.CreateBuild(_build)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err := _database.UpdateBuild(_build)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateBuild returned err: %v", err)
		}
	}
}

func TestSqlite_Client_DeleteBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

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
		// defer cleanup of the builds table
		defer _database.Sqlite.Exec("delete from builds;")

		// create the build in the database
		err = _database.CreateBuild(_build)
		if err != nil {
			t.Errorf("unable to create test build: %v", err)
		}

		err = _database.DeleteBuild(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteBuild should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteBuild returned err: %v", err)
		}
	}
}

// testBuild is a test helper function to create a
// library Build type with all fields set to their
// zero values.
func testBuild() *library.Build {
	i64 := int64(0)
	i := 0
	str := ""

	return &library.Build{
		ID:           &i64,
		RepoID:       &i64,
		PipelineID:   &i64,
		Number:       &i,
		Parent:       &i,
		Event:        &str,
		Status:       &str,
		Error:        &str,
		Enqueued:     &i64,
		Created:      &i64,
		Started:      &i64,
		Finished:     &i64,
		Deploy:       &str,
		Clone:        &str,
		Source:       &str,
		Title:        &str,
		Message:      &str,
		Commit:       &str,
		Sender:       &str,
		Author:       &str,
		Email:        &str,
		Link:         &str,
		Branch:       &str,
		Ref:          &str,
		BaseRef:      &str,
		HeadRef:      &str,
		Host:         &str,
		Runtime:      &str,
		Distribution: &str,
	}
}
