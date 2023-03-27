// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestSecret_Engine_CreateSecret(t *testing.T) {
	// setup types
	_secretRepo := testSecret()
	_secretRepo.SetID(1)
	_secretRepo.SetOrg("foo")
	_secretRepo.SetRepo("bar")
	_secretRepo.SetName("baz")
	_secretRepo.SetValue("foob")
	_secretRepo.SetType("repo")
	_secretRepo.SetCreatedAt(1)
	_secretRepo.SetCreatedBy("user")
	_secretRepo.SetUpdatedAt(1)
	_secretRepo.SetUpdatedBy("user2")

	_secretOrg := testSecret()
	_secretOrg.SetID(2)
	_secretOrg.SetOrg("foo")
	_secretOrg.SetRepo("*")
	_secretOrg.SetName("bar")
	_secretOrg.SetValue("baz")
	_secretOrg.SetType("org")
	_secretOrg.SetCreatedAt(1)
	_secretOrg.SetCreatedBy("user")
	_secretOrg.SetUpdatedAt(1)
	_secretOrg.SetUpdatedBy("user2")

	_secretShared := testSecret()
	_secretShared.SetID(3)
	_secretShared.SetOrg("foo")
	_secretShared.SetTeam("bar")
	_secretShared.SetName("baz")
	_secretShared.SetValue("foob")
	_secretShared.SetType("shared")
	_secretShared.SetCreatedAt(1)
	_secretShared.SetCreatedBy("user")
	_secretShared.SetUpdatedAt(1)
	_secretShared.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the repo secrets query
	_mock.ExpectQuery(`INSERT INTO "secrets"
("org","repo","team","name","value","type","images","events","allow_command","created_at","created_by","updated_at","updated_by","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs("foo", "bar", nil, "baz", AnyArgument{}, "repo", nil, nil, false, 1, "user", 1, "user2", 1).
		WillReturnRows(_rows)

	// ensure the mock expects the org secrets query
	_mock.ExpectQuery(`INSERT INTO "secrets"
("org","repo","team","name","value","type","images","events","allow_command","created_at","created_by","updated_at","updated_by","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs("foo", "*", nil, "bar", AnyArgument{}, "org", nil, nil, false, 1, "user", 1, "user2", 2).
		WillReturnRows(_rows)

	// ensure the mock expects the shared secrets query
	_mock.ExpectQuery(`INSERT INTO "secrets"
("org","repo","team","name","value","type","images","events","allow_command","created_at","created_by","updated_at","updated_by","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs("foo", nil, "bar", "baz", AnyArgument{}, "shared", nil, nil, false, 1, "user", 1, "user2", 3).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		secret   *library.Secret
	}{
		{
			failure:  false,
			name:     "postgres with repo",
			database: _postgres,
			secret:   _secretRepo,
		},
		{
			failure:  false,
			name:     "postgres with org",
			database: _postgres,
			secret:   _secretOrg,
		},
		{
			failure:  false,
			name:     "postgres with shared",
			database: _postgres,
			secret:   _secretShared,
		},
		{
			failure:  false,
			name:     "sqlite3 with repo",
			database: _sqlite,
			secret:   _secretRepo,
		},
		{
			failure:  false,
			name:     "sqlite3 with org",
			database: _sqlite,
			secret:   _secretOrg,
		},
		{
			failure:  false,
			name:     "sqlite3 with shared",
			database: _sqlite,
			secret:   _secretShared,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.CreateSecret(test.secret)

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
