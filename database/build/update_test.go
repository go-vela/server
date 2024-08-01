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

func TestBuild_Engine_UpdateBuild(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
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
	_repo.SetAllowEvents(api.NewEventsFromMask(1))

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetRepo(_repo)
	_build.SetNumber(1)
	_build.SetDeployPayload(nil)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "builds"
SET "repo_id"=$1,"pipeline_id"=$2,"number"=$3,"parent"=$4,"event"=$5,"event_action"=$6,"status"=$7,"error"=$8,"enqueued"=$9,"created"=$10,"started"=$11,"finished"=$12,"deploy"=$13,"deploy_number"=$14,"deploy_payload"=$15,"clone"=$16,"source"=$17,"title"=$18,"message"=$19,"commit"=$20,"sender"=$21,"sender_scm_id"=$22,"author"=$23,"email"=$24,"link"=$25,"branch"=$26,"ref"=$27,"base_ref"=$28,"head_ref"=$29,"host"=$30,"runtime"=$31,"distribution"=$32,"approved_at"=$33,"approved_by"=$34
WHERE "id" = $35`).
		WithArgs(1, nil, 1, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, AnyArgument{}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
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
