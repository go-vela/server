// SPDX-License-Identifier: Apache-2.0

package step

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStep_Engine_ListStepImageCount(t *testing.T) {
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
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("foo")
	_stepTwo.SetImage("bar")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"image", "count"}).AddRow("bar", 2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT "image", count(image) as count FROM "steps" GROUP BY "image"`).WillReturnRows(_rows)

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
		want     map[string]float64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     map[string]float64{"bar": 2},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     map[string]float64{"bar": 2},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListStepImageCount()

			if test.failure {
				if err == nil {
					t.Errorf("ListStepImageCount for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListStepImageCount for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListStepImageCount for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
