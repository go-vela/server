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

	// create the secret table
	err = _database.Sqlite.Exec(ddl.CreateSecretTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableSecret, err)
	}
}

func TestSqlite_Client_GetTypeSecretCount_Org(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("bar")
	_secretOne.SetType("org")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("bar")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secrets in the database
		err := _database.CreateSecret(_secretOne)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.CreateSecret(_secretTwo)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		got, err := _database.GetTypeSecretCount("org", "foo", "*", []string{})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretCount_Repo(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secrets in the database
		err := _database.CreateSecret(_secretOne)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.CreateSecret(_secretTwo)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		got, err := _database.GetTypeSecretCount("repo", "foo", "bar", []string{})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretCount_Shared(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secrets in the database
		err := _database.CreateSecret(_secretOne)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.CreateSecret(_secretTwo)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		got, err := _database.GetTypeSecretCount("shared", "foo", "bar", []string{})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretCount is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretCount_Shared_Wildcard(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bared")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secrets in the database
		err := _database.CreateSecret(_secretOne)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.CreateSecret(_secretTwo)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		got, err := _database.GetTypeSecretCount("shared", "foo", "*", []string{"bar", "bared"})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretCount should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretCount returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretCount is %v, want %v", got, test.want)
		}
	}
}
