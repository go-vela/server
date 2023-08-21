// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestService_Engine_GetServiceForBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_service := testService()
	_service.SetID(1)
	_service.SetRepoID(1)
	_service.SetBuildID(1)
	_service.SetNumber(1)
	_service.SetName("foo")
	_service.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"}).
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "services" WHERE build_id = $1 AND number = $2 LIMIT 1`).WithArgs(1, 1).WillReturnRows(_rows)

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
		want     *library.Service
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
			got, err := test.database.GetServiceForBuild(_build, 1)

			if test.failure {
				if err == nil {
					t.Errorf("GetServiceForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetServiceForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetServiceForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
