// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
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

func TestDeployment_Engine_ListDeployments(t *testing.T) {
	builds := []*library.Build{}

	// setup types
	_deploymentOne := testDeployment()
	_deploymentOne.SetID(1)
	_deploymentOne.SetRepoID(1)
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
	_deploymentOne.SetBuilds(builds)

	_deploymentTwo := testDeployment()
	_deploymentTwo.SetID(2)
	_deploymentTwo.SetRepoID(2)
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
	_deploymentTwo.SetBuilds(builds)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "number", "url", "commit", "ref", "task", "target", "description", "payload", "created_at", "created_by", "builds"}).
		AddRow(1, 1, 1, "https://github.com/github/octocat/deployments/1", "48afb5bdc41ad69bf22588491333f7cf71135163", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", 1, "octocat", "{}").
		AddRow(2, 2, 2, "https://github.com/github/octocat/deployments/2", "48afb5bdc41ad69bf22588491333f7cf71135164", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", 1, "octocat", "{}")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "deployments"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateDeployment(context.TODO(), _deploymentOne)
	if err != nil {
		t.Errorf("unable to create test deployment for sqlite: %v", err)
	}

	_, err = _sqlite.CreateDeployment(context.TODO(), _deploymentTwo)
	if err != nil {
		t.Errorf("unable to create test deployment for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Deployment
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Deployment{_deploymentOne, _deploymentTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Deployment{_deploymentOne, _deploymentTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListDeployments(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListDeploymentss for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListDeployments for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListDeployments for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
