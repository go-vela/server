// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package init

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInit_Engine_CreateInit(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "inits"
("repo_id","build_id","number","reporter","name","mimetype","id")
VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`).
		WithArgs(1, 1, 1, "Foobar Runtime", "foobar", "text/plain", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			err := test.database.CreateInit(_init)

			if test.failure {
				if err == nil {
					t.Errorf("CreateInit for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateInit for %s returned err: %v", test.name, err)
			}
		})
	}
}
