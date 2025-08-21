// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSecret_Engine_GetSecretForTeam(t *testing.T) {
	// setup types
	_secret := testutils.APISecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetTeam("bar")
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("shared")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy("user")
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy("user2")
	_secret.SetAllowEvents(api.NewEventsFromMask(1))
	_secret.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.SecretFromAPI(_secret)})

	_allowlistRows := testutils.CreateMockRows([]any{*types.SecretAllowlistFromAPI(_secret, "github/octocat"), *types.SecretAllowlistFromAPI(_secret, "github/octokitty")})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND team = $3 AND name = $4 LIMIT $5`).
		WithArgs(constants.SecretShared, "foo", "bar", "baz", 1).WillReturnRows(_rows)

	_mock.ExpectQuery(`SELECT * FROM "secret_repo_allowlists" WHERE secret_id = $1`).WithArgs(1).WillReturnRows(_allowlistRows)

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
		database *Engine
		want     *api.Secret
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
			got, err := test.database.GetSecretForTeam(context.TODO(), "foo", "bar", "baz")

			if test.failure {
				if err == nil {
					t.Errorf("GetSecretForTeam for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetSecretForTeam for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetSecretForTeam for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
