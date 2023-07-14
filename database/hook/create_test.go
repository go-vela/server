// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHook_Engine_CreateHook(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "hooks"
("repo_id","build_id","number","source_id","created","host","event","event_action","branch","error","status","link","webhook_id","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, nil, nil, nil, nil, nil, nil, nil, 1, 1).
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
			_, err := test.database.CreateHook(_hook)

			if test.failure {
				if err == nil {
					t.Errorf("CreateHook for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateHook for %s returned err: %v", test.name, err)
			}
		})
	}
}
