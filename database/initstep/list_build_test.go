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

func TestInitStep_Engine_ListInitStepsForBuild(t *testing.T) {
	// setup types
	_initStepOne := testInitStep()
	_initStepOne.SetID(1)
	_initStepOne.SetRepoID(1)
	_initStepOne.SetBuildID(1)
	_initStepOne.SetNumber(1)
	_initStepOne.SetReporter("Foobar Runtime")
	_initStepOne.SetName("foobar")
	_initStepOne.SetMimetype("text/plain")

	_initStepTwo := testInitStep()
	_initStepTwo.SetID(2)
	_initStepTwo.SetRepoID(1)
	_initStepTwo.SetBuildID(2)
	_initStepTwo.SetNumber(2)
	_initStepTwo.SetReporter("Foobar Runtime")
	_initStepTwo.SetName("foobar")
	_initStepTwo.SetMimetype("text/plain")

	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "initsteps" WHERE build_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "reporter", "name", "mimetype"}).
		AddRow(2, 1, 2, 2, "Foobar Runtime", "foobar", "text/plain").
		AddRow(1, 1, 1, 1, "Foobar Runtime", "foobar", "text/plain")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "initsteps" WHERE build_id = $1 ORDER BY id DESC LIMIT 10`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateInitStep(_initStepOne)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
	}

	err = _sqlite.CreateInitStep(_initStepTwo)
	if err != nil {
		t.Errorf("unable to create test init step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.InitStep
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.InitStep{_initStepTwo, _initStepOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.InitStep{_initStepTwo, _initStepOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListInitStepsForBuild(_build, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListInitStepsForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListInitStepsForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListInitStepsForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
