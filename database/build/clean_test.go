// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestBuild_Engine_CleanBuilds(t *testing.T) {
	// setup types
	errMsg := "msg"
	errStatus := "error"

	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_owner := testutils.APIUser()
	_owner.SetID(1)

	_repo.SetOwner(_owner)

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(1)
	_buildOne.SetDeployPayload(nil)
	_buildOne.SetStatus("pending")

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repo)
	_buildTwo.SetNumber(2)
	_buildTwo.SetCreated(2)
	_buildTwo.SetDeployPayload(nil)
	_buildTwo.SetStatus("running")

	// setup types
	_buildThree := testutils.APIBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepo(_repo)
	_buildThree.SetNumber(3)
	_buildThree.SetCreated(1)
	_buildThree.SetDeployPayload(nil)
	_buildThree.SetStatus("success")

	_buildFour := testutils.APIBuild()
	_buildFour.SetID(4)
	_buildFour.SetRepo(_repo)
	_buildFour.SetNumber(4)
	_buildFour.SetCreated(5)
	_buildFour.SetDeployPayload(nil)
	_buildFour.SetStatus("running")

	_buildFive := testutils.APIBuild()
	_buildFive.SetID(5)
	_buildFive.SetRepo(_repo)
	_buildFive.SetNumber(5)
	_buildFive.SetCreated(1)
	_buildFive.SetDeployPayload(nil)
	_buildFive.SetStatus("pending approval")

	// Postgres returns get updated
	_wantBuildOnePG := *_buildOne
	_wantBuildOnePG.Repo = testutils.APIRepo()
	_wantBuildOnePG.Repo.Owner = testutils.APIUser().Crop()
	_wantBuildOnePG.Status = &errStatus
	_wantBuildOnePG.Error = &errMsg

	_wantBuildOneSQ := *_buildOne
	_wantBuildOneSQ.Repo = testutils.APIRepo()
	_wantBuildOneSQ.Repo.Owner = testutils.APIUser().Crop()

	_wantBuildTwoPG := *_buildTwo
	_wantBuildTwoPG.Repo = testutils.APIRepo()
	_wantBuildTwoPG.Repo.Owner = testutils.APIUser().Crop()
	_wantBuildTwoPG.Status = &errStatus
	_wantBuildTwoPG.Error = &errMsg

	// Sqlite returns do not return updated builds
	_wantBuildTwoSQ := *_buildTwo
	_wantBuildTwoSQ.Repo = testutils.APIRepo()
	_wantBuildTwoSQ.Repo.Owner = testutils.APIUser().Crop()

	_wantBuildFivePG := *_buildFive
	_wantBuildFivePG.Repo = testutils.APIRepo()
	_wantBuildFivePG.Repo.Owner = testutils.APIUser().Crop()
	_wantBuildFivePG.Status = &errStatus
	_wantBuildFivePG.Error = &errMsg

	_wantBuildFiveSQ := *_buildFive
	_wantBuildFiveSQ.Repo = testutils.APIRepo()
	_wantBuildFiveSQ.Repo.Owner = testutils.APIUser().Crop()

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_rowsPendingRunning := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_number", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(1, 1, nil, 1, 0, "", "", "error", "msg", 0, 1, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0).
		AddRow(2, 1, nil, 2, 0, "", "", "error", "msg", 0, 2, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	// ensure the mock expects the name query
	_mock.ExpectQuery(`UPDATE "builds" SET "status"=$1,"error"=$2,"finished"=$3,"deploy_payload"=$4 WHERE created < $5 AND (status = $6 OR status = $7) RETURNING *`).
		WithArgs("error", "msg", NowTimestamp{}, AnyArgument{}, 3, "pending", "running").
		WillReturnRows(_rowsPendingRunning)

	_rowsPendingApproval := sqlmock.NewRows(
		[]string{"id", "repo_id", "pipeline_id", "number", "parent", "event", "event_action", "status", "error", "enqueued", "created", "started", "finished", "deploy", "deploy_number", "deploy_payload", "clone", "source", "title", "message", "commit", "sender", "author", "email", "link", "branch", "ref", "base_ref", "head_ref", "host", "runtime", "distribution", "timestamp"}).
		AddRow(5, 1, nil, 5, 0, "", "", "error", "msg", 0, 1, 0, 0, "", 0, nil, "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", 0)

	_mock.ExpectQuery(`UPDATE "builds" SET "status"=$1,"error"=$2,"finished"=$3,"deploy_payload"=$4 WHERE created < $5 AND status = $6 RETURNING *`).
		WithArgs("error", "msg", NowTimestamp{}, AnyArgument{}, 3, "pending approval").
		WillReturnRows(_rowsPendingApproval)

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

	_, err = _sqlite.CreateBuild(context.TODO(), _buildFour)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	_, err = _sqlite.CreateBuild(context.TODO(), _buildFive)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure    bool
		name       string
		database   *engine
		statuses   []string
		wantBuilds []*api.Build
		wantCount  int64
	}{
		{
			failure:    false,
			name:       "postgres",
			database:   _postgres,
			statuses:   []string{"pending", "running"},
			wantBuilds: []*api.Build{&_wantBuildOnePG, &_wantBuildTwoPG},
			wantCount:  2,
		},
		{
			failure:    false,
			name:       "postgres with pending approval",
			database:   _postgres,
			statuses:   []string{"pending approval"},
			wantBuilds: []*api.Build{&_wantBuildFivePG},
			wantCount:  1,
		},
		{
			failure:    false,
			name:       "sqlite3",
			database:   _sqlite,
			statuses:   []string{"pending", "running"},
			wantBuilds: []*api.Build{&_wantBuildOneSQ, &_wantBuildTwoSQ},
			wantCount:  2,
		},
		{
			failure:    false,
			name:       "sqlite3 with pending approval",
			database:   _sqlite,
			statuses:   []string{"pending approval"},
			wantBuilds: []*api.Build{&_wantBuildFiveSQ},
			wantCount:  1,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotBuilds, gotCount, err := test.database.CleanBuilds(context.TODO(), "msg", test.statuses, 3)

			if test.failure {
				if err == nil {
					t.Errorf("CleanBuilds for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanBuilds for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.wantBuilds, gotBuilds); diff != "" {
				t.Errorf("CleanBuilds for %s is a mismatch (-want +got):\n%s", test.name, diff)
			}

			if test.wantCount != gotCount {
				t.Errorf("CleanBuilds Count for %s is %d, want %d", test.name, gotCount, test.wantCount)
			}
		})
	}
}
