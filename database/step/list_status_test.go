// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestStep_Engine_ListStepStatusCount(t *testing.T) {
	// setup types
	_stepOne := testutils.APIStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")

	_stepTwo := testutils.APIStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("foo")
	_stepTwo.SetImage("bar")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"status", "count"}).
		AddRow("pending", 0).
		AddRow("failure", 0).
		AddRow("killed", 0).
		AddRow("running", 0).
		AddRow("success", 0)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT "status", count(status) as count FROM "steps" GROUP BY "status"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateStep(ctx, _stepOne)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	_, err = _sqlite.CreateStep(ctx, _stepTwo)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     map[string]float64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want: map[string]float64{
				"pending": 0,
				"failure": 0,
				"killed":  0,
				"running": 0,
				"success": 0,
			},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want: map[string]float64{
				"pending": 0,
				"failure": 0,
				"killed":  0,
				"running": 0,
				"success": 0,
			},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListStepStatusCount(ctx)

			if test.failure {
				if err == nil {
					t.Errorf("ListStepStatusCount for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListStepStatusCount for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListStepStatusCount for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
