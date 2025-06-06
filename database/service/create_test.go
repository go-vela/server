// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestService_Engine_CreateService(t *testing.T) {
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
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "services"
("build_id","repo_id","number","name","image","status","error","exit_code","created","started","finished","host","runtime","distribution","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING "id"`).
		WithArgs(1, 1, 1, "foo", "bar", nil, nil, nil, nil, nil, nil, nil, nil, nil, 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
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
			got, err := test.database.CreateService(context.TODO(), _service)

			if test.failure {
				if err == nil {
					t.Errorf("CreateService for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateService for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _service) {
				t.Errorf("CreateService for %s returned %s, want %s", test.name, got, _service)
			}
		})
	}
}
