// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestRepo_Engine_UpdateRepo(t *testing.T) {
	// setup types
	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.GetOwner().SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetPreviousName("oldName")
	_repo.SetApproveBuild(constants.ApproveForkAlways)
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(api.NewEventsFromMask(1))
	_repo.SetApprovalTimeout(5)
	_repo.SetCustomProps(map[string]any{"foo": "bar"})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "repos"
SET "user_id"=$1,"hash"=$2,"org"=$3,"name"=$4,"full_name"=$5,"link"=$6,"clone"=$7,"branch"=$8,"topics"=$9,"build_limit"=$10,"timeout"=$11,"counter"=$12,"visibility"=$13,"private"=$14,"trusted"=$15,"active"=$16,"allow_events"=$17,"pipeline_type"=$18,"previous_name"=$19,"approve_build"=$20,"approval_timeout"=$21,"install_id"=$22,"custom_props"=$23
WHERE "id" = $24`).
		WithArgs(1, AnyArgument{}, "foo", "bar", "foo/bar", nil, nil, nil, AnyArgument{}, AnyArgument{}, AnyArgument{}, AnyArgument{}, "public", false, false, false, 1, "yaml", "oldName", constants.ApproveForkAlways, 5, 0, `{"foo":"bar"}`, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateRepo(context.TODO(), _repo)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
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
			got, err := test.database.UpdateRepo(context.TODO(), _repo)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _repo) {
				t.Errorf("UpdateRepo for %s returned %s, want %s", test.name, got, _repo)
			}
		})
	}
}
