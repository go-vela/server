// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetHook(t *testing.T) {
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
	_hook.SetWebhookID(123456)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Hook
	}{
		{
			failure: false,
			want:    _hook,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the hook in the database
			err := _database.CreateHook(test.want)
			if err != nil {
				t.Errorf("unable to create test hook: %v", err)
			}
		}

		got, err := _database.GetHook(1, _repo)

		// cleanup the hooks table
		_ = _database.Sqlite.Exec("DELETE FROM hooks;")

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

func TestSqlite_Client_GetLastHook(t *testing.T) {
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
	_hook.SetWebhookID(123456)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}

	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Hook
	}{
		{
			failure: false,
			want:    _hook,
		},
		{
			failure: false,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the hook in the database
			err := _database.CreateHook(test.want)
			if err != nil {
				t.Errorf("unable to create test hook: %v", err)
			}
		}

		got, err := _database.GetLastHook(_repo)

		// cleanup the hooks table
		_ = _database.Sqlite.Exec("DELETE FROM hooks;")

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

func TestSqlite_Client_CreateHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(123456)

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
		// defer cleanup of the hooks table
		defer _database.Sqlite.Exec("delete from hooks;")

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

func TestSqlite_Client_UpdateHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(123456)

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
		// defer cleanup of the hooks table
		defer _database.Sqlite.Exec("delete from hooks;")

		// create the hook in the database
		err := _database.CreateHook(_hook)
		if err != nil {
			t.Errorf("unable to create test hook: %v", err)
		}

		err = _database.UpdateHook(_hook)

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

func TestSqlite_Client_DeleteHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(123456)

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
		// defer cleanup of the hooks table
		defer _database.Sqlite.Exec("delete from hooks;")

		// create the hook in the database
		err := _database.CreateHook(_hook)
		if err != nil {
			t.Errorf("unable to create test hook: %v", err)
		}

		err = _database.DeleteHook(1)

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
		ID:        &i64,
		RepoID:    &i64,
		BuildID:   &i64,
		Number:    &i,
		SourceID:  &str,
		Created:   &i64,
		Host:      &str,
		Event:     &str,
		Branch:    &str,
		Error:     &str,
		Status:    &str,
		Link:      &str,
		WebhookID: &i64,
	}
}
