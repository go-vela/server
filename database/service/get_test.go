// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestService_Engine_GetService(t *testing.T) {
	// setup types
	_service := testutils.APIService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.ServiceFromAPI(_service)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "services" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateService(context.TODO(), _service)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     *api.Service
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _service,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _service,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetService(context.TODO(), 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetService for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetService for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetService for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
