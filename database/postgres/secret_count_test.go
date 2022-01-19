// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"

	"gorm.io/gorm"
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

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectOrgSecretsCount, "foo").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

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

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectRepoSecretsCount, "foo", "bar").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

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

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectSharedSecretsCount, "foo", "bar").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

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

func TestPostgres_Client_GetTypeSecretCount_Shared_Wildcard(t *testing.T) {
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
	_secretTwo.SetTeam("bared")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}
	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT count(*) FROM \"secrets\" WHERE (type = 'shared' AND org = $1) AND LOWER(team) IN ($2,$3)").WillReturnRows(_rows)

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
