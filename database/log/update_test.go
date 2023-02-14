// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/library"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLog_Engine_UpdateLog(t *testing.T) {
	// setup types
	_service := testLog()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)
	_service.SetData([]byte{})

	_step := testLog()
	_step.SetID(2)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)
	_step.SetData([]byte{})

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the service query
	_mock.ExpectExec(`UPDATE "logs"
SET "build_id"=$1,"repo_id"=$2,"service_id"=$3,"step_id"=$4,"data"=$5
WHERE "id" = $6`).
		WithArgs(1, 1, 1, nil, AnyArgument{}, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// ensure the mock expects the step query
	_mock.ExpectExec(`UPDATE "logs"
SET "build_id"=$1,"repo_id"=$2,"service_id"=$3,"step_id"=$4,"data"=$5
WHERE "id" = $6`).
		WithArgs(1, 1, nil, 1, AnyArgument{}, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateLog(_service)
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
		logs     []*library.Log
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			logs:     []*library.Log{_service, _step},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			logs:     []*library.Log{_service, _step},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, log := range test.logs {
				err = test.database.UpdateLog(log)

				if test.failure {
					if err == nil {
						t.Errorf("UpdateLog for %s should have returned err", test.name)
					}

					return
				}

				if err != nil {
					t.Errorf("UpdateLog for %s returned err: %v", test.name, err)
				}
			}
		})
	}
}
