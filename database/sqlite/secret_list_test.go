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

	// create the secret table
	err = _database.Sqlite.Exec(ddl.CreateSecretTable).Error
	if err != nil {
		log.Fatalf("unable to create %s table: %v", constants.TableSecret, err)
	}
}

func TestSqlite_Client_GetSecretList(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		for _, secret := range test.want {
			// create the secret in the database
			err := _database.CreateSecret(secret)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetSecretList()

		if test.failure {
			if err == nil {
				t.Errorf("GetSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretList_Org(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("bar")
	_secretOne.SetType("org")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("bar")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretTwo, _secretOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		for _, secret := range test.want {
			// create the secret in the database
			err := _database.CreateSecret(secret)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetTypeSecretList("org", "foo", "*", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretList_Repo(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretTwo, _secretOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		for _, secret := range test.want {
			// create the secret in the database
			err := _database.CreateSecret(secret)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetTypeSecretList("repo", "foo", "bar", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetTypeSecretList_Shared(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")

	_secretTwo := testSecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretTwo, _secretOne},
		},
	}

	// run tests
	for _, test := range tests {
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		for _, secret := range test.want {
			// create the secret in the database
			err := _database.CreateSecret(secret)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetTypeSecretList("shared", "foo", "bar", 1, 10)

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}
