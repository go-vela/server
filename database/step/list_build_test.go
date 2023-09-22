// SPDX-License-Identifier: Apache-2.0

package step

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestStep_Engine_ListStepsForBuild(t *testing.T) {
	// setup types
	_build := testBuild()
	_build.SetID(1)
	_build.SetRepoID(1)
	_build.SetNumber(1)

	_stepOne := testStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")

	_stepTwo := testStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("foo")
	_stepTwo.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "steps" WHERE build_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"}).
		AddRow(2, 1, 1, 2, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps" WHERE build_id = $1 ORDER BY id DESC LIMIT 10`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateStep(_stepOne)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	_, err = _sqlite.CreateStep(_stepTwo)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Step
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Step{_stepTwo, _stepOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Step{_stepTwo, _stepOne},
		},
	}

	filters := map[string]interface{}{}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListStepsForBuild(_build, filters, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListStepsForBuild for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListStepsForBuild for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListStepsForBuild for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
