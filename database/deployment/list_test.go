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

func TestDeployment_Engine_ListDeployments(t *testing.T) {
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
	_repoOne.SetAllowEvents(api.NewEventsFromMask(1))
	_repoOne.SetPipelineType(constants.PipelineTypeYAML)
	_repoOne.SetTopics([]string{})

	_repoBuild := new(api.Repo)
	_repoBuild.SetID(1)

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(_repoBuild)
	_build.SetNumber(1)
	_build.SetDeployNumber(0)
	_build.SetDeployPayload(nil)

	_deploymentOne := testutils.APIDeployment()
	_deploymentOne.SetID(1)
	_deploymentOne.SetRepo(_repoOne)
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
	_deploymentOne.SetBuilds([]*api.Build{_build})

	_deploymentTwo := testutils.APIDeployment()
	_deploymentTwo.SetID(2)
	_deploymentTwo.SetRepo(_repoOne)
	_deploymentTwo.SetNumber(2)
	_deploymentTwo.SetURL("https://github.com/github/octocat/deployments/2")
	_deploymentTwo.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135164")
	_deploymentTwo.SetRef("refs/heads/master")
	_deploymentTwo.SetTask("vela-deploy")
	_deploymentTwo.SetTarget("production")
	_deploymentTwo.SetDescription("Deployment request from Vela")
	_deploymentTwo.SetPayload(map[string]string{"foo": "test1"})
	_deploymentTwo.SetCreatedAt(1)
	_deploymentTwo.SetCreatedBy("octocat")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.DeploymentFromAPI(_deploymentOne), *types.DeploymentFromAPI(_deploymentTwo)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repoOne)})

	_userRows := testutils.CreateMockRows([]any{*types.UserFromAPI(_owner)})

	_buildRows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_build)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "deployments"`).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_buildRows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Deployment{_deploymentOne, _deploymentTwo},
		[]*api.User{_owner},
		[]*api.Repo{_repoOne},
		[]*api.Build{_build},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Deployment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Deployment{_deploymentOne, _deploymentTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Deployment{_deploymentOne, _deploymentTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListDeployments(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListDeployments for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListDeployments for %s returned err: %v", test.name, err)
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
				t.Errorf("GetDeploymentForRepo for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
