// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestDeployment_Engine_GetDeployment(t *testing.T) {
	// setup types
	_owner := testutils.APIUser().Crop()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetOwner(_owner)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetAllowEvents(api.NewEventsFromMask(1))
	_repo.SetPipelineType(constants.PipelineTypeYAML)
	_repo.SetTopics([]string{})

	_deploymentOne := testutils.APIDeployment()
	_deploymentOne.SetID(1)
	_deploymentOne.SetRepo(_repo)
	_deploymentOne.SetNumber(1)
	_deploymentOne.SetURL("https://github.com/github/octocat/deployments/1")
	_deploymentOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_deploymentOne.SetRef("refs/heads/master")
	_deploymentOne.SetTask("vela-deploy")
	_deploymentOne.SetTarget("production")
	_deploymentOne.SetDescription("Deployment request from Vela")
	_deploymentOne.SetPayload(map[string]string{"foo": "test1"})
	_deploymentOne.SetCreatedAt(1)
	_deploymentOne.SetCreatedBy("octocat")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.DeploymentFromAPI(_deploymentOne)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "deployments" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Deployment{_deploymentOne},
		[]*api.User{_owner},
		[]*api.Repo{_repo},
		[]*api.Build{},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Deployment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _deploymentOne,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _deploymentOne,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetDeployment(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetDeployment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetDeployment for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(
				test.want,
				got,
				cmp.Options{
					cmp.FilterPath(func(p cmp.Path) bool {
						return p.String() == "Builds.Repo.Owner" // ignore this nested struct due to ToAPI
					}, cmp.Ignore()),
				},
			); diff != "" {
				t.Errorf("GetDeployment for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
