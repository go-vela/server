// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestSecret_Engine_DeleteSecret(t *testing.T) {
	// setup types
	_secretRepo := testutils.APISecret()
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

	_secretOrg := testutils.APISecret()
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

	_secretShared := testutils.APISecret()
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
	_secretShared.SetRepoAllowlist([]string{"github/octocat"})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectBegin()
	// ensure the mock expects the repo query
	_mock.ExpectExec(`DELETE FROM "secrets" WHERE "secrets"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`DELETE FROM "secret_repo_allowlists" WHERE secret_id = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 0))

	_mock.ExpectCommit()
	_mock.ExpectBegin()

	// ensure the mock expects the org query
	_mock.ExpectExec(`DELETE FROM "secrets" WHERE "secrets"."id" = $1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`DELETE FROM "secret_repo_allowlists" WHERE secret_id = $1`).
		WithArgs(2).
		WillReturnResult(sqlmock.NewResult(1, 0))

	_mock.ExpectCommit()
	_mock.ExpectBegin()

	// ensure the mock expects the shared query
	_mock.ExpectExec(`DELETE FROM "secrets" WHERE "secrets"."id" = $1`).
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectExec(`DELETE FROM "secret_repo_allowlists" WHERE secret_id = $1`).
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_mock.ExpectCommit()

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secretRepo)
	if err != nil {
		t.Errorf("unable to create test repo secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(context.TODO(), _secretOrg)
	if err != nil {
		t.Errorf("unable to create test org secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(context.TODO(), _secretShared)
	if err != nil {
		t.Errorf("unable to create test shared secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		secret   *api.Secret
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
			err = test.database.DeleteSecret(context.TODO(), test.secret)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteSecret for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteSecret for %s returned err: %v", test.name, err)
			}
		})
	}
}
