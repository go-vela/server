// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package deployment

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestDeployment_Engine_CreateDeployment(t *testing.T) {
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

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "deployments"
("number","repo_id","url","commit","ref","task","target","description","payload","created_at","created_by","builds","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING "id"`).
		WithArgs(1, 1, "https://github.com/github/octocat/deployments/1", "48afb5bdc41ad69bf22588491333f7cf71135163", "refs/heads/master", "vela-deploy", "production", "Deployment request from Vela", "{\"foo\":\"test1\"}", 1, "octocat", "{}", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			_, err := test.database.CreateDeployment(context.TODO(), _deploymentOne)

			if test.failure {
				if err == nil {
					t.Errorf("CreateDeployment for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Create for %s returned err: %v", test.name, err)
			}
		})
	}
}
