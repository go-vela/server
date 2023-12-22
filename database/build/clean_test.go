// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuild_Engine_CleanBuilds(t *testing.T) {
	// setup types
	_buildOne := testBuild()
	_buildOne.SetID(1)
	_buildOne.SetRepoID(1)
	_buildOne.SetNumber(1)
	_buildOne.SetCreated(1)
	_buildOne.SetStatus("pending")

	_buildTwo := testBuild()
	_buildTwo.SetID(2)
	_buildTwo.SetRepoID(1)
	_buildTwo.SetNumber(2)
	_buildTwo.SetCreated(2)
	_buildTwo.SetStatus("running")

	// setup types
	_buildThree := testBuild()
	_buildThree.SetID(3)
	_buildThree.SetRepoID(1)
	_buildThree.SetNumber(3)
	_buildThree.SetCreated(1)
	_buildThree.SetStatus("success")

	_buildFour := testBuild()
	_buildFour.SetID(4)
	_buildFour.SetRepoID(1)
	_buildFour.SetNumber(4)
	_buildFour.SetCreated(5)
	_buildFour.SetStatus("running")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the name query
	_mock.ExpectExec(`UPDATE "builds" SET "status"=$1,"error"=$2,"finished"=$3,"deploy_number"=$4,"deploy_payload"=$5 WHERE created < $6 AND (status = 'running' OR status = 'pending')`).
		WithArgs("error", "msg", NowTimestamp{}, 0, AnyArgument{}, 3).
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
