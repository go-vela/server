// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSecret_Engine_ListSecretsForOrg(t *testing.T) {
	// setup types
	_secretOne := testutils.APISecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("bar")
	_secretOne.SetType("org")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")
	_secretOne.SetAllowEvents(api.NewEventsFromMask(1))
	_secretOne.SetRepoAllowlist([]string{})

	_secretTwo := testutils.APISecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("bar")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")
	_secretTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_secretTwo.SetRepoAllowlist([]string{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := testutils.CreateMockRows([]any{*types.SecretFromAPI(_secretTwo), *types.SecretFromAPI(_secretOne)})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 ORDER BY id DESC LIMIT $3`).
		WithArgs(constants.SecretOrg, "foo", 10).WillReturnRows(_rows)

	_mock.ExpectQuery(`SELECT * FROM "secret_repo_allowlist" WHERE secret_id IN ($1,$2)`).WithArgs(2, 1).WillReturnRows(sqlmock.NewRows([]string{}))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSecret(context.TODO(), _secretOne)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSecret(context.TODO(), _secretTwo)
	if err != nil {
		t.Errorf("unable to create test secret for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Secret
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Secret{_secretTwo, _secretOne},
		},
		{
			failure:  false,
			name:     "sqlite",
			database: _sqlite,
			want:     []*api.Secret{_secretTwo, _secretOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSecretsForOrg(context.TODO(), "foo", filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSecretsForOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecretsForOrg for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSecretsForOrg for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
