// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

func TestSecret_Engine_GetSecretForRepo(t *testing.T) {
	// setup types
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")

	_secret := testSecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetRepo("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("repo")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy("user")
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(1, "repo", "foo", "bar", "", "baz", "foob", nil, nil, false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND repo = $3 AND name = $4 LIMIT 1`).
		WithArgs(constants.SecretRepo, "foo", "bar", "baz").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secret)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Secret
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _secret,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _secret,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetSecretForRepo(context.TODO(), "baz", _repo)

			if test.failure {
				if err == nil {
					t.Errorf("GetSecretForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetSecretForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetSecretForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
