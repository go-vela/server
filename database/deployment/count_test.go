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
	"github.com/go-vela/types/raw"
)

func TestDeployment_Engine_CountDeployments(t *testing.T) {
	buildOne := new(library.Build)
	buildOne.SetID(1)
	buildOne.SetRepoID(1)
	buildOne.SetPipelineID(1)
	buildOne.SetNumber(1)
	buildOne.SetParent(1)
	buildOne.SetEvent("push")
	buildOne.SetEventAction("")
	buildOne.SetStatus("running")
	buildOne.SetError("")
	buildOne.SetEnqueued(1563474077)
	buildOne.SetCreated(1563474076)
	buildOne.SetStarted(1563474078)
	buildOne.SetFinished(1563474079)
	buildOne.SetDeploy("")
	buildOne.SetDeployPayload(raw.StringSliceMap{"foo": "test1"})
	buildOne.SetClone("https://github.com/github/octocat.git")
	buildOne.SetSource("https://github.com/github/octocat/deployments/1")
	buildOne.SetTitle("push received from https://github.com/github/octocat")
	buildOne.SetMessage("First commit...")
	buildOne.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	buildOne.SetSender("OctoKitty")
	buildOne.SetAuthor("OctoKitty")
	buildOne.SetEmail("OctoKitty@github.com")
	buildOne.SetLink("https://example.company.com/github/octocat/1")
	buildOne.SetBranch("main")
	buildOne.SetRef("refs/heads/main")
	buildOne.SetBaseRef("")
	buildOne.SetHeadRef("changes")
	buildOne.SetHost("example.company.com")
	buildOne.SetRuntime("docker")
	buildOne.SetDistribution("linux")

	builds := []*library.Build{}
	builds = append(builds, buildOne)

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
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "deployments"`).WillReturnRows(_rows)

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
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountDeployments(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CountDeployments for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountDeployments for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountHooks for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
