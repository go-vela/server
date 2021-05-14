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
)

func TestPostgres_Client_GetSecretList(t *testing.T) {
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
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "repo", "foo", "bar", "", "baz", "foob", nil, nil, false).
		AddRow(1, "repo", "foo", "bar", "", "foob", "baz", nil, nil, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListSecrets).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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

func TestPostgres_Client_GetTypeSecretList_Org(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("bar")
	_secretOne.SetType("org")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("bar")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")

	// create the new fake sql database
	_sql, _mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("unable to create new sql mock database: %v", err)
	}
	defer _sql.Close()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "org", "foo", "*", "", "baz", "bar", nil, nil, false).
		AddRow(1, "org", "foo", "*", "", "bar", "baz", nil, nil, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListOrgSecrets).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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

func TestPostgres_Client_GetTypeSecretList_Repo(t *testing.T) {
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
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "repo", "foo", "bar", "", "baz", "foob", nil, nil, false).
		AddRow(1, "repo", "foo", "bar", "", "foob", "baz", nil, nil, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListRepoSecrets).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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

func TestPostgres_Client_GetTypeSecretList_Shared(t *testing.T) {
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
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command"},
	).AddRow(1, "shared", "foo", "", "bar", "baz", "foob", nil, nil, false).
		AddRow(1, "shared", "foo", "", "bar", "foob", "baz", nil, nil, false)

	// ensure the mock expects the query
	_mock.ExpectQuery(dml.ListSharedSecrets).WillReturnRows(_rows)

	// setup the database client
	_database, err := NewTest(_sql)
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

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
