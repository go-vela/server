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
)

func TestSecret_Engine_GetSecretForTeam(t *testing.T) {
	// setup types
	_secret := testutils.APISecret()
	_secret.SetID(1)
	_secret.SetOrg("foo")
	_secret.SetOrgSCMID(1)
	_secret.SetTeam("bar")
	_secret.SetTeamSCMID(1)
	_secret.SetName("baz")
	_secret.SetValue("foob")
	_secret.SetType("shared")
	_secret.SetCreatedAt(1)
	_secret.SetCreatedBy("user")
	_secret.SetUpdatedAt(1)
	_secret.SetUpdatedBy("user2")
	_secret.SetAllowEvents(api.NewEventsFromMask(1))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "org_scm_id", "repo", "repo_scm_id", "team", "team_scm_id", "name", "value", "images", "allow_events", "allow_command", "allow_substitution", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(1, "shared", "foo", 1, "", nil, "bar", 1, "baz", "foob", nil, 1, false, false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND team = $3 AND name = $4 LIMIT $5`).
		WithArgs(constants.SecretShared, "foo", "bar", "baz", 1).WillReturnRows(_rows)

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
