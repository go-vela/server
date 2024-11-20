// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestBuild_Engine_CleanBuilds(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(api.NewEventsFromMask(1))

	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("foo")
	_owner.SetToken("bar")

	_repo.SetOwner(_owner)

	_buildOne := testutils.APIBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepo(_repo)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(1)
	_buildOne.SetStatus("pending")

	_buildTwo := testutils.APIBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepo(_repo)
	_buildTwo.SetNumber(2)
	_buildTwo.SetCreated(2)
	_buildTwo.SetStatus("running")

	// setup types
	_buildThree := testutils.APIBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepo(_repo)
	_buildThree.SetNumber(3)
	_buildThree.SetCreated(1)
	_buildThree.SetStatus("success")

	_buildFour := testutils.APIBuild()
	_buildFour.SetID(4)
	_buildFour.SetRepo(_repo)
	_buildFour.SetNumber(4)
	_buildFour.SetCreated(5)
	_buildFour.SetStatus("running")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the name query
	_mock.ExpectExec(`UPDATE "builds" SET "error"=$1,"finished"=$2,"status"=$3 WHERE created < $4 AND (status = 'running' OR status = 'pending')`).
		WithArgs("msg", NowTimestamp{}, "error", 3).
		WillReturnResult(sqlmock.NewResult(1, 2))

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
			got, err := test.database.CleanBuilds(context.TODO(), "msg", 3)

			if test.failure {
				if err == nil {
					t.Errorf("CleanBuilds for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanBuilds for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CleanBuilds for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
