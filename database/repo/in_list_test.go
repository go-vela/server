// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestRepo_Engine_GetReposInList(t *testing.T) {
	// setup types
	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetOwner(_owner)
	_repoOne.SetHash("baz")
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")
	_repoOne.SetVisibility("public")
	_repoOne.SetPipelineType("yaml")
	_repoOne.SetTopics([]string{})
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("baz")
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne), *types.RepoFromAPI(_repoTwo)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE full_name IN ($1,$2)`).WithArgs("foo/bar", "bar/foo").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repoOne)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	_, err = _sqlite.CreateRepo(context.TODO(), _repoTwo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Repo
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Repo{_repoOne, _repoTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Repo{_repoTwo, _repoOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetReposInList(t.Context(), []string{"foo/bar", "bar/foo"})

			if test.failure {
				if err == nil {
					t.Errorf("GetReposInList for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetReposInList for %s returned err: %v", test.name, err)
			}

			// ignore hash and owner because we do not decrypt or preload for this func
			if diff := cmp.Diff(test.want, got, cmpopts.IgnoreFields(api.Repo{}, "Hash", "Owner")); diff != "" {
				t.Errorf("GetReposInList for %s is %s", test.name, diff)
			}
		})
	}
}
