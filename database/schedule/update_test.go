// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_UpdateSchedule(t *testing.T) {
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_schedule := testSchedule()
	_schedule.SetID(1)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1)
	_schedule.SetCreatedBy("user1")
	_schedule.SetUpdatedAt(1)
	_schedule.SetUpdatedBy("user2")
	_schedule.SetRepo(_repo)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "schedules"
SET "repo_id"=$1,"active"=$2,"name"=$3,"entry"=$4,"created_at"=$5,"created_by"=$6,"updated_at"=$7,"updated_by"=$8,"scheduled_at"=$9
WHERE "id" = $10`).
		WithArgs(1, false, "nightly", "0 0 * * *", 1, "user1", time.Now().UTC().Unix(), "user2", nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateSchedule(_schedule)
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
			err = test.database.UpdateSchedule(_schedule)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSchedule for %s returned err: %v", test.name, err)
			}
		})
	}
}
