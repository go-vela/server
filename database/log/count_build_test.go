// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestLog_Engine_CountLogsForBuild(t *testing.T) {
	// setup types
	_service := testutils.APILog()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)
	_service.SetCreatedAt(1)

	_step := testutils.APILog()
	_step.SetID(2)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)
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
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "logs" WHERE build_id = $1`).WithArgs(1).WillReturnRows(_rows)

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
			got, err := test.database.CountLogsForBuild(context.TODO(), _build)

			if test.failure {
				if err == nil {
					t.Errorf("CountLogsForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountLogsForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountLogsForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
