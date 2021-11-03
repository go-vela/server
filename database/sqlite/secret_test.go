// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
)

func TestSqlite_Client_GetSecret_Org(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("*")
	_secret.SetName("bar")
	_secret.SetValue("baz")
	_secret.SetType("org")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)
	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,
			want:    _secret,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the secret in the database
			err := _database.CreateSecret(test.want)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetSecret("org", "foo", "*", "bar")

		// cleanup the secrets table
		_ = _database.Sqlite.Exec("DELETE FROM secrets;")

		if test.failure {
			if err == nil {
				t.Errorf("GetSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetSecret returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetSecret is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetSecret_Repo(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,
			want:    _secret,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the secret in the database
			err := _database.CreateSecret(test.want)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetSecret("repo", "foo", "bar", "baz")

		// cleanup the secrets table
		_ = _database.Sqlite.Exec("DELETE FROM secrets;")

		if test.failure {
			if err == nil {
				t.Errorf("GetSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetSecret returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetSecret is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_GetSecret_Shared(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetTeam("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("shared")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)

	// setup the test database client
	_database, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new sqlite test database: %v", err)
	}
	defer func() { _sql, _ := _database.Sqlite.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,
			want:    _secret,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		if test.want != nil {
			// create the secret in the database
			err := _database.CreateSecret(test.want)
			if err != nil {
				t.Errorf("unable to create test secret: %v", err)
			}
		}

		got, err := _database.GetSecret("shared", "foo", "bar", "baz")

		// cleanup the secrets table
		_ = _database.Sqlite.Exec("DELETE FROM secrets;")

		if test.failure {
			if err == nil {
				t.Errorf("GetSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetSecret returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetSecret is %v, want %v", got, test.want)
		}
	}
}

func TestSqlite_Client_CreateSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		err := _database.CreateSecret(_secret)

		if test.failure {
			if err == nil {
				t.Errorf("CreateSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateSecret returned err: %v", err)
		}
	}
}

func TestSqlite_Client_UpdateSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secret in the database
		err := _database.CreateSecret(_secret)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.UpdateSecret(_secret)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateSecret returned err: %v", err)
		}
	}
}

func TestSqlite_Client_DeleteSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy(1)
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy(1)
	_secret.SetLastBuildID(1)

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
		// defer cleanup of the secrets table
		defer _database.Sqlite.Exec("delete from secrets;")

		// create the secret in the database
		err := _database.CreateSecret(_secret)
		if err != nil {
			t.Errorf("unable to create test secret: %v", err)
		}

		err = _database.DeleteSecret(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteSecret should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteSecret returned err: %v", err)
		}
	}
}

// testSecret is a test helper function to create a
// library Secret type with all fields set to their
// zero values.
func testSecret() *library.Secret {
	i64 := int64(0)
	str := ""
	booL := false
	var arr []string

	return &library.Secret{
		ID:           &i64,
		Org:          &str,
		Repo:         &str,
		Team:         &str,
		Name:         &str,
		Value:        &str,
		Type:         &str,
		Images:       &arr,
		Events:       &arr,
		AllowCommand: &booL,
		CreatedAt:    &i64,
		CreatedBy:    &i64,
		UpdatedAt:    &i64,
		UpdatedBy:    &i64,
		LastBuildID:  &i64,
	}
}
