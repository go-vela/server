// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestStep_Engine_UpdateStep(t *testing.T) {
	// setup types
	_step := testutils.APIStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "steps"
SET "build_id"=$1,"repo_id"=$2,"number"=$3,"name"=$4,"image"=$5,"stage"=$6,"status"=$7,"error"=$8,"exit_code"=$9,"created"=$10,"started"=$11,"finished"=$12,"host"=$13,"runtime"=$14,"distribution"=$15,"report_as"=$16
WHERE "id" = $17`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateStep(ctx, _step)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
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
			got, err := test.database.UpdateStep(ctx, _step)

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
