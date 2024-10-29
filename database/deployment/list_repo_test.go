// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestDeployment_Engine_ListDeploymentsForRepo(t *testing.T) {
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
	_repoOne.SetInstallID(0)

	_repoTwo := testutils.APIRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOwner(_owner)
	_repoTwo.SetHash("bazey")
	_repoTwo.SetOrg("fooey")
	_repoTwo.SetName("barey")
	_repoTwo.SetFullName("fooey/barey")
	_repoTwo.SetVisibility("public")
	_repoTwo.SetAllowEvents(api.NewEventsFromMask(1))
	_repoTwo.SetPipelineType(constants.PipelineTypeYAML)
	_repoTwo.SetTopics([]string{})
	_repoTwo.SetInstallID(0)

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
	_deploymentTwo.SetRepo(_repoTwo)
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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "url", "commit", "ref", "task", "target", "description", "payload", "created_at", "created_by", "builds"}).
		AddRow(1, 1, 1, "https://github.com/github/octocat/deployments/1", "48afb5bdc41ad69bf22588491333f7cf71135163", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", 1, "octocat", "{1}")

	_buildRows := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_number", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "", "", "", "", 0, 0, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "deployments" WHERE repo_id = $1 ORDER BY number DESC LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_buildRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	sqlitePopulateTables(
		t,
		_sqlite,
		[]*api.Deployment{_deploymentOne, _deploymentTwo},
		[]*api.User{},
		[]*api.Repo{},
		[]*api.Build{_build},
	)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*api.Deployment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Deployment{_deploymentOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Deployment{_deploymentOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListDeploymentsForRepo(context.TODO(), _repoOne, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListDeploymentsForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListDeploymentsForRepo for %s returned err: %v", test.name, err)
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
