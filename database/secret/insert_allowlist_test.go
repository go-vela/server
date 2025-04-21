// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestSecret_Engine_InsertAllowlist(t *testing.T) {
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
	_secret.SetAllowEvents(api.NewEventsFromMask(1))
	_secret.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the repo secrets query
	_mock.ExpectQuery(`INSERT INTO "secret_repo_allowlist"
("secret_id","repo")
VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING RETURNING "id"`).
		WithArgs(1, "github/octocat", 1, "github/octokitty").
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		secret   *api.Secret
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			secret:   _secret,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			secret:   _secret,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := InsertAllowlist(context.TODO(), test.database.client, test.secret)

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
