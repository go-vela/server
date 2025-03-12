// SPDX-License-Identifier: Apache-2.0

package deployment

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestDeployment_Engine_UpdateDeployment(t *testing.T) {
	// setup types
	_repoOne := testutils.APIRepo()
	_repoOne.SetID(1)
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")

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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "deployments"
SET "number"=$1,"repo_id"=$2,"url"=$3,"commit"=$4,"ref"=$5,"task"=$6,"target"=$7,"description"=$8,"payload"=$9,"created_at"=$10,"created_by"=$11,"builds"=$12
WHERE "id" = $13`).
		WithArgs(1, 1, "https://github.com/github/octocat/deployments/1", "48afb5bdc41ad69bf22588491333f7cf71135163", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", 1, "octocat", "{}", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateDeployment(context.TODO(), _deploymentOne)
	if err != nil {
		t.Errorf("unable to create test deployment for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.UpdateDeployment(context.TODO(), _deploymentOne)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateDeployment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateDeployment for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _deploymentOne) {
				t.Errorf("UpdateDeployment for %s is %v, want %v", test.name, got, _deploymentOne)
			}
		})
	}
}
