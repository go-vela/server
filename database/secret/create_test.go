// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSecret_Engine_CreateSecret(t *testing.T) {
	// setup types
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
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "secrets"
("org","repo","team","name","value","type","images","events","allow_command","created_at","created_by","updated_at","updated_by","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs("foo", "bar", nil, "baz", AnyArgument{}, "repo", nil, nil, false, 1, "user", 1, "user2", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.CreateSecret(_secret)

			if test.failure {
				if err == nil {
					t.Errorf("CreateSecret for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateSecret for %s returned err: %v", test.name, err)
			}
		})
	}
}
