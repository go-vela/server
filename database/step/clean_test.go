// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestStep_Engine_CleanStep(t *testing.T) {
	// setup types
	_stepOne := testutils.APIStep()
	_stepOne.SetID(1)
	_stepOne.SetRepoID(1)
	_stepOne.SetBuildID(1)
	_stepOne.SetNumber(1)
	_stepOne.SetName("foo")
	_stepOne.SetImage("bar")
	_stepOne.SetCreated(1)
	_stepOne.SetStatus("running")

	_stepTwo := testutils.APIStep()
	_stepTwo.SetID(2)
	_stepTwo.SetRepoID(1)
	_stepTwo.SetBuildID(1)
	_stepTwo.SetNumber(2)
	_stepTwo.SetName("foo")
	_stepTwo.SetImage("bar")
	_stepTwo.SetCreated(1)
	_stepTwo.SetStatus("pending")

	_stepThree := testutils.APIStep()
	_stepThree.SetID(3)
	_stepThree.SetRepoID(1)
	_stepThree.SetBuildID(1)
	_stepThree.SetNumber(3)
	_stepThree.SetName("foo")
	_stepThree.SetImage("bar")
	_stepThree.SetCreated(1)
	_stepThree.SetStatus("success")

	_stepFour := testutils.APIStep()
	_stepFour.SetID(4)
	_stepFour.SetRepoID(1)
	_stepFour.SetBuildID(1)
	_stepFour.SetNumber(4)
	_stepFour.SetName("foo")
	_stepFour.SetImage("bar")
	_stepFour.SetCreated(5)
	_stepFour.SetStatus("pending")

	_postgres, _mock := testPostgres(t)

	ctx := context.TODO()
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the name query
	_mock.ExpectExec(`UPDATE "steps" SET "status"=$1,"error"=$2,"finished"=$3 WHERE created < $4 AND (status = 'running' OR status = 'pending')`).
		WithArgs("error", "msg", NowTimestamp{}, 3).
		WillReturnResult(sqlmock.NewResult(1, 2))

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

	_, err = _sqlite.CreateStep(ctx, _stepThree)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	_, err = _sqlite.CreateStep(ctx, _stepFour)
	if err != nil {
		t.Errorf("unable to create test step for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *Engine
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     2,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     2,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CleanSteps(ctx, "msg", 3)

			if test.failure {
				if err == nil {
					t.Errorf("CleanSteps for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanSteps for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CleanSteps for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
