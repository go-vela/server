// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSecret_Engine_ListSecrets(t *testing.T) {
	// setup types
	_secretOne := testutils.APISecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")
	_secretOne.SetAllowEvents(api.NewEventsFromMask(1))
	_secretOne.SetRepoAllowlist([]string{})

	_secretTwo := testutils.APISecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")
	_secretTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_secretTwo.SetRepoAllowlist([]string{})

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.SecretFromAPI(_secretOne), *types.SecretFromAPI(_secretTwo)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "secrets"`).WillReturnRows(_rows)

	_mock.ExpectQuery(`SELECT * FROM "secret_repo_allowlists" WHERE secret_id IN ($1,$2)`).WithArgs(1, 2).WillReturnRows(sqlmock.NewRows([]string{}))

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
			want:     []*api.Secret{_secretOne, _secretTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSecrets(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListSecrets for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecrets for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSecrets for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
