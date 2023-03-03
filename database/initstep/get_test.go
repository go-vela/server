// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package initstep

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestInitStep_Engine_GetInitStep(t *testing.T) {
	// setup types
	_initStep := testInitStep()
	_initStep.SetID(1)
	_initStep.SetRepoID(1)
	_initStep.SetBuildID(1)
	_initStep.SetNumber(1)
	_initStep.SetReporter("Foobar Runtime")
	_initStep.SetName("foobar")
	_initStep.SetMimetype("text/plain")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "reporter", "name", "mimetype"},
	).AddRow(1, 1, 1, 1, "Foobar Runtime", "foobar", "text/plain")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "initsteps" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInitStep(_initStep)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.InitStep
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _initStep,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _initStep,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetInitStep(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetInitStep for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetInitStep for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetInitStep for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
