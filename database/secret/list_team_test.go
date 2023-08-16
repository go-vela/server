// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestSecret_Engine_ListSecretsForTeam(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the name count query
	_mock.ExpectQuery(`SELECT count(*) FROM "secrets" WHERE type = $1 AND org = $2 AND team = $3`).
		WithArgs(constants.SecretShared, "foo", "bar").WillReturnRows(_rows)

	// create expected name query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(2, "shared", "foo", "", "bar", "foob", "baz", nil, nil, false, 1, "user", 1, "user2").
		AddRow(1, "shared", "foo", "", "bar", "baz", "foob", nil, nil, false, 1, "user", 1, "user2")

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND team = $3 ORDER BY id DESC LIMIT 10`).
		WithArgs(constants.SecretShared, "foo", "bar").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(_secretOne)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(_secretTwo)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Secret
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Secret{_secretTwo, _secretOne},
		},
		{
			failure:  false,
			name:     "sqlite",
			database: _sqlite,
			want:     []*library.Secret{_secretTwo, _secretOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListSecretsForTeam("foo", "bar", filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSecretsForTeam for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecretsForTeam for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSecretsForTeam for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}

func TestSecret_Engine_ListSecretsForTeams(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name count query result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the name count query
	_mock.ExpectQuery(`SELECT count(*) FROM "secrets" WHERE type = $1 AND org = $2 AND LOWER(team) IN ($3,$4)`).
		WithArgs(constants.SecretShared, "foo", "foo", "bar").WillReturnRows(_rows)

	// create expected name query result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(2, "shared", "foo", "", "bar", "foob", "baz", nil, nil, false, 1, "user", 1, "user2").
		AddRow(1, "shared", "foo", "", "bar", "baz", "foob", nil, nil, false, 1, "user", 1, "user2")

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND LOWER(team) IN ($3,$4) ORDER BY id DESC LIMIT 10`).
		WithArgs(constants.SecretShared, "foo", "foo", "bar").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(_secretOne)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(_secretTwo)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Secret
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Secret{_secretTwo, _secretOne},
		},
		{
			failure:  false,
			name:     "sqlite",
			database: _sqlite,
			want:     []*library.Secret{_secretTwo, _secretOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListSecretsForTeams("foo", []string{"foo", "bar"}, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSecretsForTeams for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecretsForTeams for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSecretsForTeams for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
