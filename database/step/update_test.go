// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStep_Engine_UpdateStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "steps"
SET "build_id"=$1,"repo_id"=$2,"number"=$3,"name"=$4,"image"=$5,"stage"=$6,"status"=$7,"error"=$8,"exit_code"=$9,"created"=$10,"started"=$11,"finished"=$12,"host"=$13,"runtime"=$14,"distribution"=$15
WHERE "id" = $16`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateStep(_step)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
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
			got, err := test.database.UpdateStep(_step)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateStep for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _step) {
				t.Errorf("UpdateStep for %s returned %s, want %s", test.name, got, _step)
			}
		})
	}
}
