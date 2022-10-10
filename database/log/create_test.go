// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package log

import (
	"github.com/go-vela/types/library"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLog_Engine_CreateLog(t *testing.T) {
	// setup types
	_service := testLog()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetServiceID(1)

	_step := testLog()
	_step.SetID(2)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetStepID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the service query
	_mock.ExpectQuery(`INSERT INTO "logs"
("build_id","repo_id","service_id","step_id","data","id")
VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(1, 1, 1, nil, AnyArgument{}, 1).
		WillReturnRows(_rows)

	// ensure the mock expects the step query
	_mock.ExpectQuery(`INSERT INTO "logs"
("build_id","repo_id","service_id","step_id","data","id")
VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`).
		WithArgs(1, 1, nil, 1, AnyArgument{}, 2).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
				err := test.database.CreateLog(log)

				if test.failure {
					if err == nil {
						t.Errorf("CreateLog for %s should have returned err", test.name)
					}

					return
				}

				if err != nil {
					t.Errorf("CreateLog for %s returned err: %v", test.name, err)
				}
			}
		})
	}
}
