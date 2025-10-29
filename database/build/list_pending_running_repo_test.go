// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestBuild_Engine_ListPendingAndRunningBuildsForRepo(t *testing.T) {
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
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))
	_repoOne.SetTopics([]string{})

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("bar")
	_repoTwo.SetOrg("foo")
	_repoTwo.SetName("baz")
	_repoTwo.SetFullName("foo/baz")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetPipelineType("yaml")
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_repoTwo.SetTopics([]string{})

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repoOne)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("running")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repoOne)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("pending")
	_buildTwo.SetCreated(1)
	_buildTwo.SetDeployPayload(nil)

	_buildThree := testutils.APIBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepo(_repoTwo)
	_buildThree.SetNumber(1)
	_buildThree.SetStatus("pending")
	_buildThree.SetCreated(1)
	_buildThree.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildTwo), *types.BuildFromAPI(_buildOne)})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE repo_id = $1 AND (status = 'running' OR status = 'pending' OR status = 'pending approval')`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := createTestBuild(t.Context(), _sqlite, _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = createTestBuild(t.Context(), _sqlite, _buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = createTestBuild(t.Context(), _sqlite, _buildThree)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Build
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Build{_buildTwo, _buildOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Build{_buildOne, _buildTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListPendingAndRunningBuildsForRepo(context.TODO(), _repoOne)

			if test.failure {
				if err == nil {
					t.Errorf("ListPendingAndRunningBuildsForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListPendingAndRunningBuildsForRepo for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListPendingAndRunningBuildsForRepo for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
