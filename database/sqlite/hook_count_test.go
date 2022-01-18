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

	// create the hook table
	err = _database.Sqlite.Exec(ddl.CreateHookTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableHook, err)
	}
}

func TestSqlite_Client_GetRepoHookCount(t *testing.T) {
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
		// defer cleanup of the hooks table
		defer _database.Sqlite.Exec("delete from hooks;")

		// create the hooks in the database
		err := _database.CreateHook(_hookOne)
		if err != nil {
			t.Errorf("unable to create test hook: %v", err)
		}

		err = _database.CreateHook(_hookTwo)
		if err != nil {
			t.Errorf("unable to create test hook: %v", err)
		}

		got, err := _database.GetRepoHookCount(_repo)

		if test.failure {
			if err == nil {
				t.Errorf("GetRepoHookCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetRepoHookCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetRepoHookCount is %v, want %v", got, test.want)
		}
	}
}
