// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestService_Engine_DeleteService(t *testing.T) {
	// setup types
	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "services" WHERE "services"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateService(_service)
	if err != nil {
		t.Errorf("unable to create test service for sqlite: %v", err)
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
			err = test.database.DeleteService(_service)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteService for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteService for %s returned err: %v", test.name, err)
			}
		})
	}
}
