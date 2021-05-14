// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
)

func TestPostgres_Client_GetTypeSecretCount_Org(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("org")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectOrgSecretsCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		got, err := _database.GetTypeSecretCount("org", "foo", "*")

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

func TestPostgres_Client_GetTypeSecretCount_Repo(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectRepoSecretsCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		got, err := _database.GetTypeSecretCount("repo", "foo", "bar")

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

func TestPostgres_Client_GetTypeSecretCount_Shared(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.SelectSharedSecretsCount).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
		got, err := _database.GetTypeSecretCount("shared", "foo", "bar")

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
