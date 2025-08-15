// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestSecret_Engine_PruneAllowlist(t *testing.T) {
	// setup types
	_secret := testutils.APISecret()
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
	_secret.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the repo query
	_mock.ExpectExec(`DELETE FROM "secret_repo_allowlists" WHERE secret_id = $1 AND repo NOT IN ($2,$3)`).
		WithArgs(1, "github/octocat", "github/octokitty").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secret)
	if err != nil {
		t.Errorf("unable to create test repo secret for sqlite: %v", err)
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
			secret:   _secret,
		},
		{
			failure:  false,
			name:     "sqlite3 with repo",
			database: _sqlite,
			secret:   _secret,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = PruneAllowlist(context.TODO(), test.database.client, test.secret)

			if test.failure {
				if err == nil {
					t.Errorf("PruneAllowlist for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("PruneAllowlist for %s returned err: %v", test.name, err)
			}
		})
	}
}

func TestSecret_Engine_PruneAllowlistEmpty(t *testing.T) {
	// setup types
	_secret := testutils.APISecret()
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
	_secret.SetRepoAllowlist([]string{})

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the repo query
	_mock.ExpectExec(`DELETE FROM "secret_repo_allowlists" WHERE secret_id = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secret)
	if err != nil {
		t.Errorf("unable to create test repo secret for sqlite: %v", err)
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
			secret:   _secret,
		},
		{
			failure:  false,
			name:     "sqlite3 with repo",
			database: _sqlite,
			secret:   _secret,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = PruneAllowlist(context.TODO(), test.database.client, test.secret)

			if test.failure {
				if err == nil {
					t.Errorf("PruneAllowlist for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("PruneAllowlist for %s returned err: %v", test.name, err)
			}
		})
	}
}
