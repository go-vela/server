// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_DeleteSchedule(t *testing.T) {
	_schedule := testSchedule()
	_schedule.SetID(1)
	_schedule.SetRepoID(1)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1)
	_schedule.SetCreatedBy("user1")
	_schedule.SetUpdatedAt(1)
	_schedule.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "schedules" WHERE "schedules"."id" = $1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	ctx := context.TODO()

	err := _sqlite.CreateSchedule(ctx, _schedule)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err = test.database.DeleteSchedule(_schedule)

			if test.failure {
				if err == nil {
					t.Errorf("DeleteSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteSchedule for %s returned err: %v", test.name, err)
			}
		})
	}
}
