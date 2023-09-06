// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestHook_Engine_GetHookByWebhookID(t *testing.T) {
	// setup types
	_hook := testHook()
	_hook.SetID(1)
	_hook.SetRepoID(1)
	_hook.SetBuildID(1)
	_hook.SetNumber(1)
	_hook.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	_hook.SetWebhookID(123456)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "source_id", "created", "host", "event", "event_action", "branch", "error", "status", "link", "webhook_id"},
	).AddRow(1, 1, 1, 1, "c8da1302-07d6-11ea-882f-4893bca275b8", 0, "", "", "", "", "", "", "", 123456)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "hooks" WHERE webhook_id = $1 LIMIT 1`).WithArgs(123456).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateHook(context.TODO(), _hook)
	if err != nil {
		t.Errorf("unable to create test hook for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Hook
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _hook,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _hook,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetHookByWebhookID(context.TODO(), 123456)

			if test.failure {
				if err == nil {
					t.Errorf("GetHookByWebhookID for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetHookByWebhookID for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetHookByWebhookID for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
