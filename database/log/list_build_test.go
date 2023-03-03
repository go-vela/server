// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestLog_Engine_ListLogsForBuild(t *testing.T) {
	// setup types
	_initStep := testLog()
	_initStep.SetID(1)
	_initStep.SetRepoID(1)
	_initStep.SetBuildID(1)
	_initStep.SetInitStepID(1)
	_initStep.SetData([]byte{})

	_service := testLog()
	_service.SetID(2)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)
	_service.SetData([]byte{})

	_step := testLog()
	_step.SetID(3)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)
	_step.SetData([]byte{})

	_build := testBuild()
	_build.SetID(1)
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "logs" WHERE build_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "build_id", "repo_id", "service_id", "step_id", "init_step_id", "data"}).
		AddRow(1, 1, 1, 0, 0, 1, []byte{}).
		AddRow(2, 1, 1, 1, 0, 0, []byte{}).
		AddRow(3, 1, 1, 0, 1, 0, []byte{})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "logs" WHERE build_id = $1 ORDER BY step_id ASC LIMIT 10`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(_initStep)
	if err != nil {
		t.Errorf("unable to create test init step log for sqlite: %v", err)
	}

	err = _sqlite.CreateLog(_service)
	if err != nil {
		t.Errorf("unable to create test service log for sqlite: %v", err)
	}

	err = _sqlite.CreateLog(_step)
	if err != nil {
		t.Errorf("unable to create test step log for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Log{_initStep, _service, _step},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Log{_initStep, _service, _step},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListLogsForBuild(_build, 1, 10)

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
