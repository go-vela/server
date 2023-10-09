// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBuild_Engine_UpdateBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "builds"
SET "repo_id"=$1,"pipeline_id"=$2,"number"=$3,"parent"=$4,"event"=$5,"event_action"=$6,"status"=$7,"error"=$8,"enqueued"=$9,"created"=$10,"started"=$11,"finished"=$12,"deploy"=$13,"deploy_payload"=$14,"clone"=$15,"source"=$16,"title"=$17,"message"=$18,"commit"=$19,"sender"=$20,"author"=$21,"email"=$22,"link"=$23,"branch"=$24,"ref"=$25,"base_ref"=$26,"head_ref"=$27,"host"=$28,"runtime"=$29,"distribution"=$30
WHERE "id" = $31`).
		WithArgs(1, nil, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, AnyArgument{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateBuild(context.TODO(), _build)
	if err != nil {
		t.Errorf("unable to create test build for sqlite: %v", err)
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
			got, err := test.database.UpdateBuild(context.TODO(), _build)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _build) {
				t.Errorf("UpdateBuild for %s returned %s, want %s", test.name, got, _build)
			}
		})
	}
}
