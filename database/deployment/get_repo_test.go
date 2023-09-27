// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestDeployment_Engine_GetDeploymentForRepo(t *testing.T) {
	builds := new([]library.Build)

	// setup types
	_deploymentOne := testDeployment()
	_deploymentOne.SetID(1)
	_deploymentOne.SetRepoID(1)
	_deploymentOne.SetNumber(1)
	_deploymentOne.SetURL("https://github.com/github/octocat/deployments/1")
	_deploymentOne.SetUser("octocat")
	_deploymentOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	_deploymentOne.SetRef("refs/heads/master")
	_deploymentOne.SetTask("vela-deploy")
	_deploymentOne.SetTarget("production")
	_deploymentOne.SetDescription("Deployment request from Vela")
	_deploymentOne.SetPayload(map[string]string{"foo": "test1"})
	_deploymentOne.SetBuilds(builds)

	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetUserID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "url", "user", "commit", "ref", "task", "target", "description", "payload", "builds"}).
		AddRow(1, 1, 1, "https://github.com/github/octocat/deployments/1", "octocat", "48afb5bdc41ad69bf22588491333f7cf71135163", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", "{}")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "deployments" WHERE repo_id = $1 AND number = $2 LIMIT 1`).WithArgs(1, 1).WillReturnRows(_rows)

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
		want     *library.Deployment
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
			got, err := test.database.GetDeploymentForRepo(context.TODO(), _repo, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetDeploymentForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetDeploymentForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetDeploymentForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
