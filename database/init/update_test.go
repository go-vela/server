// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInit_Engine_UpdateInit(t *testing.T) {
	// setup types
	_init := testInit()
	_init.SetID(1)
	_init.SetRepoID(1)
	_init.SetBuildID(1)
	_init.SetNumber(1)
	_init.SetReporter("Foobar Runtime")
	_init.SetName("foobar")
	_init.SetMimetype("text/plain")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "inits"
SET "repo_id"=$1,"build_id"=$2,"number"=$3,"reporter"=$4,"name"=$5,"mimetype"=$6
WHERE "id" = $7`).
		WithArgs(1, 1, 1, "Foobar Runtime", "foobar", "text/plain", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInit(_init)
	if err != nil {
		t.Errorf("unable to create test init for sqlite: %v", err)
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
			err = test.database.UpdateInit(_init)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateInit for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateInit for %s returned err: %v", test.name, err)
			}
		})
	}
}
