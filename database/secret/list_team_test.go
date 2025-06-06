// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSecret_Engine_ListSecretsForTeam(t *testing.T) {
	// setup types
	_secretOne := testutils.APISecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")
	_secretOne.SetAllowEvents(api.NewEventsFromMask(1))
	_secretOne.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})

	_secretTwo := testutils.APISecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
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

	_allowlistRows := testutils.CreateMockRows([]any{*types.SecretAllowlistFromAPI(_secretOne, "github/octocat"), *types.SecretAllowlistFromAPI(_secretOne, "github/octokitty")})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND team = $3 ORDER BY id DESC LIMIT $4`).
		WithArgs(constants.SecretShared, "foo", "bar", 10).WillReturnRows(_rows)

	_mock.ExpectQuery(`SELECT * FROM "secret_repo_allowlists" WHERE secret_id IN ($1,$2)`).WithArgs(2, 1).WillReturnRows(_allowlistRows)

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
			got, err := test.database.ListSecretsForTeam(context.TODO(), "foo", "bar", filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSecretsForTeam for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecretsForTeam for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListSecretsForTeam mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSecret_Engine_ListSecretsForTeams(t *testing.T) {
	// setup types
	_secretOne := testutils.APISecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")
	_secretOne.SetAllowEvents(api.NewEventsFromMask(1))
	_secretOne.SetRepoAllowlist([]string{"github/octocat", "github/octokitty"})

	_secretTwo := testutils.APISecret()
	_secretTwo.SetID(2)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")
	_secretTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_secretTwo.SetRepoAllowlist([]string{"alpha/beta", "github/octocat"})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "allow_events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"}).
		AddRow(2, "shared", "foo", "", "bar", "foob", "baz", nil, 1, false, 1, "user", 1, "user2").
		AddRow(1, "shared", "foo", "", "bar", "baz", "foob", nil, 1, false, 1, "user", 1, "user2")

	_allowlistRows := testutils.CreateMockRows(
		[]any{
			*types.SecretAllowlistFromAPI(_secretOne, "github/octocat"),
			*types.SecretAllowlistFromAPI(_secretOne, "github/octokitty"),
			*types.SecretAllowlistFromAPI(_secretTwo, "alpha/beta"),
			*types.SecretAllowlistFromAPI(_secretTwo, "github/octocat"),
		})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "secrets" WHERE type = $1 AND org = $2 AND LOWER(team) IN ($3,$4) ORDER BY id DESC LIMIT $5`).
		WithArgs(constants.SecretShared, "foo", "foo", "bar", 10).WillReturnRows(_rows)

	_mock.ExpectQuery(`SELECT * FROM "secret_repo_allowlists" WHERE secret_id IN ($1,$2)`).WithArgs(2, 1).WillReturnRows(_allowlistRows)

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
			got, err := test.database.ListSecretsForTeams(context.TODO(), "foo", []string{"foo", "bar"}, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSecretsForTeams for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSecretsForTeams for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListSecretsForTeams mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
