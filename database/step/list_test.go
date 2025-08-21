// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestStep_Engine_ListSteps(t *testing.T) {
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
	_stepTwo.SetBuildID(2)
	_stepTwo.SetNumber(1)
	_stepTwo.SetName("bar")
	_stepTwo.SetImage("foo")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := testutils.CreateMockRows([]any{*types.StepFromAPI(_stepOne), *types.StepFromAPI(_stepTwo)})

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "steps"`).WillReturnRows(_rows)

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
		want     []*api.Step
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Step{_stepOne, _stepTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Step{_stepOne, _stepTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSteps(ctx)

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
