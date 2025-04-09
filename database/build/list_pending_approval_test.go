// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestBuild_Engine_ListPendingApprovalBuilds(t *testing.T) {
	// setup types
	_repoOwner := new(api.User)
	_repoOwner.SetID(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetOwner(_repoOwner)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)
	_buildOne.SetNumber(1)
	_buildOne.SetStatus("pending approval")
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repo)
	_buildTwo.SetNumber(2)
	_buildTwo.SetStatus("pending approval")
	_buildTwo.SetCreated(3)
	_buildTwo.SetDeployPayload(nil)

	_buildThree := testutils.APIBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepo(_repo)
	_buildThree.SetNumber(3)
	_buildThree.SetStatus("pending approval")
	_buildThree.SetCreated(6)
	_buildThree.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected name query result in mock
	_rows := testutils.CreateMockRows([]any{*types.BuildFromAPI(_buildOne), *types.BuildFromAPI(_buildTwo)})

	_repoRows := testutils.CreateMockRows([]any{*types.RepoFromAPI(_repo)})

	// ensure the mock expects the name query
	_mock.ExpectQuery(`SELECT * FROM "builds" WHERE builds.created < $1 AND builds.status = 'pending approval'`).WithArgs("5").WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)

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

	_, err = _sqlite.CreateBuild(context.TODO(), _buildThree)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
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
			want:     []*api.Build{_buildOne, _buildTwo},
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
			got, err := test.database.ListPendingApprovalBuilds(context.TODO(), "5")

			if test.failure {
				if err == nil {
					t.Errorf("ListPendingApprovalBuilds for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListPendingApprovalBuilds for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("ListPendingApprovalBuilds for %s (-got +want): %s", test.name, diff)
			}
		})
	}
}
