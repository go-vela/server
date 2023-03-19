// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInitStep_Engine_UpdateInitStep(t *testing.T) {
	// setup types
	_initStep := testInitStep()
	_initStep.SetID(1)
	_initStep.SetRepoID(1)
	_initStep.SetBuildID(1)
	_initStep.SetNumber(1)
	_initStep.SetReporter("Foobar Runtime")
	_initStep.SetName("foobar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "initsteps"
SET "repo_id"=$1,"build_id"=$2,"number"=$3,"reporter"=$4,"name"=$5,"mimetype"=$6
WHERE "id" = $7`).
		WithArgs(1, 1, 1, "Foobar Runtime", "foobar", "text/plain", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInitStep(_initStep)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
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
			err = test.database.UpdateInitStep(_initStep)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateInitStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateInitStep for %s returned err: %v", test.name, err)
			}
		})
	}
}
