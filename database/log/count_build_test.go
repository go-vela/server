// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLog_Engine_CountLogsForBuild(t *testing.T) {
	// setup types
	_init := testLog()
	_init.SetID(1)
	_init.SetRepoID(1)
	_init.SetBuildID(1)
	_init.SetInitID(1)

	_service := testLog()
	_service.SetID(2)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)

	_step := testLog()
	_step.SetID(3)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)

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

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(_init)
	if err != nil {
		t.Errorf("unable to create test init log for sqlite: %v", err)
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
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     3,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     3,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountLogsForBuild(_build)

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
