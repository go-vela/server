// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package step

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestStep_Engine_GetStep(t *testing.T) {
	// setup types
	_step := testStep()
	_step.SetID(1)
	_step.SetRepoID(1)
	_step.SetBuildID(1)
	_step.SetNumber(1)
	_step.SetName("foo")
	_step.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"}).
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateStep(_step)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.Step
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _step,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _step,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetStep(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetStep for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetStep for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
