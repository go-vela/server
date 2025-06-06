// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestBuild_Engine_ListPendingAndRunningBuilds(t *testing.T) {
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

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("running")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repo)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("pending")
	_buildTwo.SetCreated(1)
	_buildTwo.SetDeployPayload(nil)

	_queueOne := new(api.QueueBuild)
	_queueOne.SetCreated(1)
	_queueOne.SetFullName("foo/bar")
	_queueOne.SetNumber(1)
	_queueOne.SetStatus("running")

	_queueTwo := new(api.QueueBuild)
	_queueTwo.SetCreated(1)
	_queueTwo.SetFullName("foo/bar")
	_queueTwo.SetNumber(2)
	_queueTwo.SetStatus("pending")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := sqlmock.NewRows([]string{"created", "full_name", "number", "status"}).AddRow(1, "foo/bar", 2, "pending").AddRow(1, "foo/bar", 1, "running")

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT builds.created, builds.number, builds.status, repos.full_name FROM "builds" INNER JOIN repos ON builds.repo_id = repos.id WHERE builds.created > $1 AND (builds.status = 'running' OR builds.status = 'pending')`).WithArgs("0").WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _buildOne)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildTwo)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create repo table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.QueueBuild
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.QueueBuild{_queueTwo, _queueOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.QueueBuild{_queueTwo, _queueOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListPendingAndRunningBuilds(context.TODO(), "0")

			if test.failure {
				if err == nil {
					t.Errorf("ListPendingAndRunningBuilds for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListPendingAndRunningBuilds for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListPendingAndRunningBuilds for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
