// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestStep_Engine_DeleteStep(t *testing.T) {
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
	_mock.ExpectExec(`DELETE FROM "steps" WHERE "steps"."id" = $1`).
		WithArgs(1).
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
			err = test.database.DeleteStep(ctx, _step)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteStep for %s returned err: %v", test.name, err)
			}
		})
	}
}
