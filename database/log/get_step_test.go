// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/types/library"
)

func TestLog_Engine_GetLogForStep(t *testing.T) {
	// setup types
	_log := testLog()
	_log.SetID(1)
	_log.SetRepoID(1)
	_log.SetBuildID(1)
	_log.SetStepID(1)
	_log.SetData([]byte{})

	_step := testStep()
	_step.SetID(1)
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "step_id", "step_id", "data"}).
		AddRow(1, 1, 1, 0, 1, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "logs" WHERE step_id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(context.TODO(), _log)
	if err != nil {
		t.Errorf("unable to create test log for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _log,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _log,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetLogForStep(context.TODO(), _step)

			if test.failure {
				if err == nil {
					t.Errorf("GetLogForStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetLogForStep for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetLogForStep for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
