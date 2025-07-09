// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestLog_Engine_ListLogsForBuild(t *testing.T) {
	// setup types
	_service := testutils.APILog()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)
	_service.SetData([]byte{})
	_service.SetCreatedAt(1)

	_step := testutils.APILog()
	_step.SetID(2)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)
	_step.SetData([]byte{})
	_step.SetCreatedAt(1)

	_repo := testutils.APIRepo()
	_repo.SetID(1)

	_build := testutils.APIBuild()
	_build.SetID(1)
	_build.SetID(1)
	_build.SetRepo(_repo)
	_build.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.LogFromAPI(_service), *types.LogFromAPI(_step)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "logs" WHERE build_id = $1 ORDER BY service_id ASC NULLS LAST,step_id ASC LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(context.TODO(), _service)
	if err != nil {
		t.Errorf("unable to create test service log for sqlite: %v", err)
	}

	err = _sqlite.CreateLog(context.TODO(), _step)
	if err != nil {
		t.Errorf("unable to create test step log for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     []*api.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Log{_service, _step},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Log{_service, _step},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListLogsForBuild(context.TODO(), _build, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListLogsForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListLogsForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListLogsForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
