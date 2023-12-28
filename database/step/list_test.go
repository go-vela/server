// SPDX-License-Identifier: Apache-2.0

package step

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestStep_Engine_ListSteps(t *testing.T) {
	// setup types
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
	_stepTwo.SetBuildID(2)
	_stepTwo.SetNumber(1)
	_stepTwo.SetName("bar")
	_stepTwo.SetImage("foo")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "steps"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "build_id", "number", "name", "image", "stage", "status", "error", "exit_code", "created", "started", "finished", "host", "runtime", "distribution"}).
		AddRow(1, 1, 1, 1, "foo", "bar", "", "", "", 0, 0, 0, 0, "", "", "").
		AddRow(2, 1, 2, 1, "bar", "foo", "", "", "", 0, 0, 0, 0, "", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps"`).WillReturnRows(_rows)

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
			want:     []*library.Step{_stepOne, _stepTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Step{_stepOne, _stepTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSteps()

			if test.failure {
				if err == nil {
					t.Errorf("ListSteps for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSteps for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSteps for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
