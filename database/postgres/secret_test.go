// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetSecret_Org(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("*")
	_secret.SetName("bar")
	_secret.SetValue("baz")
	_secret.SetType("org")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectOrgSecret, "foo", "bar").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "org", "foo", "*", "", "bar", "baz", "{}", "{}", false)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,

			want: _secret,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetSecret("org", "foo", "*", "bar")

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

func TestPostgres_Client_GetSecret_Repo(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectRepoSecret, "foo", "bar", "baz").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "repo", "foo", "bar", "", "baz", "foob", "{}", "{}", false)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,

			want: _secret,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetSecret("repo", "foo", "bar", "baz")

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

func TestPostgres_Client_GetSecret_Shared(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetTeam("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("shared")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectSharedSecret, "foo", "bar", "baz").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "shared", "foo", "", "bar", "baz", "foob", "{}", "{}", false)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Secret
	}{
		{
			failure: false,

			want: _secret,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetSecret("shared", "foo", "bar", "baz")

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

func TestPostgres_Client_CreateSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "secrets" ("org","repo","team","name","value","type","images","events","allow_command","id") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "id"`).
		WithArgs("foo", "bar", "", "baz", AnyArgument{}, "repo", "{}", "{}", false, 1).
		WillReturnRows(_rows)

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

func TestPostgres_Client_UpdateSecret(t *testing.T) {
	// setup types
	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "secrets" SET "org"=$1,"repo"=$2,"team"=$3,"name"=$4,"value"=$5,"type"=$6,"images"=$7,"events"=$8,"allow_command"=$9 WHERE "id" = $10`).
		WithArgs("foo", "bar", "", "baz", AnyArgument{}, "repo", "{}", "{}", false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
		err := _database.UpdateSecret(_secret)

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

func TestPostgres_Client_DeleteSecret(t *testing.T) {
	// setup types

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteSecret, 1).Statement

	// ensure the mock expects the query
	_mock.ExpectExec(_query.SQL.String()).WillReturnResult(sqlmock.NewResult(1, 1))

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
		err := _database.DeleteSecret(1)

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
	arr := []string{}
	booL := false

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
	}
}
