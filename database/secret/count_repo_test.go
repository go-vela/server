// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestSecret_Engine_CountSecretsForRepo(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.GetOwner().SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")

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

	_secretTwo := testutils.APISecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("bar")
	_secretTwo.SetRepo("foo")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "secrets" WHERE type = $1 AND org = $2 AND repo = $3`).
		WithArgs(constants.SecretRepo, "foo", "bar").WillReturnRows(_rows)

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
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     1,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     1,
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountSecretsForRepo(context.TODO(), _repo, filters)

			if test.failure {
				if err == nil {
					t.Errorf("CountSecretsForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountSecretsForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountSecretsForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
