// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestService_Engine_UpdateService(t *testing.T) {
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
	_mock.ExpectExec(`UPDATE "services" SET "build_id"=$1,"repo_id"=$2,"number"=$3,"name"=$4,"image"=$5,"status"=$6,"error"=$7,"exit_code"=$8,"created"=$9,"started"=$10,"finished"=$11,"host"=$12,"runtime"=$13,"distribution"=$14 WHERE "id" = $15`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

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
			got, err := test.database.UpdateService(context.TODO(), _service)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateService for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateService for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _service) {
				t.Errorf("UpdateService for %s returned %s, want %s", test.name, got, _service)
			}
		})
	}
}
