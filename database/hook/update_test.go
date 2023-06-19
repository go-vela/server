// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHook_Engine_UpdateHook(t *testing.T) {
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

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "hooks"
SET "repo_id"=$1,"build_id"=$2,"number"=$3,"source_id"=$4,"created"=$5,"host"=$6,"event"=$7,"event_action"=$8,"branch"=$9,"error"=$10,"status"=$11,"link"=$12,"webhook_id"=$13
WHERE "id" = $14`).
		WithArgs(1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", nil, nil, nil, nil, nil, nil, nil, nil, 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateHook(_hook)
	if err != nil {
		t.Errorf("unable to create test hook for sqlite: %v", err)
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
			_, err = test.database.UpdateHook(_hook)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateHook for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateHook for %s returned err: %v", test.name, err)
			}
		})
	}
}
