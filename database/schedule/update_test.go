// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_UpdateSchedule_Config(t *testing.T) {
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_schedule := testSchedule()
	_schedule.SetID(1)
	_schedule.SetRepoID(1)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1)
	_schedule.SetCreatedBy("user1")
	_schedule.SetUpdatedAt(1)
	_schedule.SetUpdatedBy("user2")
	_schedule.SetBranch("main")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "schedules"
SET "repo_id"=$1,"active"=$2,"name"=$3,"entry"=$4,"created_at"=$5,"created_by"=$6,"updated_at"=$7,"updated_by"=$8,"scheduled_at"=$9,"branch"=$10
WHERE "id" = $11`).
		WithArgs(1, false, "nightly", "0 0 * * *", 1, "user1", NowTimestamp{}, "user2", nil, "main", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSchedule(context.TODO(), _schedule)
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
			got, err := test.database.UpdateSchedule(context.TODO(), _schedule, true)
			_schedule.SetUpdatedAt(got.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSchedule for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _schedule) {
				t.Errorf("UpdateSchedule for %s returned %s, want %s", test.name, got, _schedule)
			}
		})
	}
}

func TestSchedule_Engine_UpdateSchedule_NotConfig(t *testing.T) {
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_schedule := testSchedule()
	_schedule.SetID(1)
	_schedule.SetRepoID(1)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1)
	_schedule.SetCreatedBy("user1")
	_schedule.SetUpdatedAt(1)
	_schedule.SetUpdatedBy("user2")
	_schedule.SetScheduledAt(1)
	_schedule.SetBranch("main")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "schedules" SET "scheduled_at"=$1 WHERE "id" = $2`).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSchedule(context.TODO(), _schedule)
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
			got, err := test.database.UpdateSchedule(context.TODO(), _schedule, false)

			if test.failure {
				if err == nil {
					t.Errorf("UpdateSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateSchedule for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _schedule) {
				t.Errorf("CreateSchedule for %s returned %s, want %s", test.name, got, _schedule)
			}
		})
	}
}
