// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/server/api/types"
)

func TestSchedule_Engine_GetSchedule(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at"},
	).AddRow(1, 1, false, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules" WHERE id = $1 LIMIT 1`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateSchedule(_schedule)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	_schedule.SetRepo(nil)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *types.Schedule
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _schedule,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _schedule,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetSchedule(1)

			if test.failure {
				if err == nil {
					t.Errorf("GetSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetSchedule for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("GetSchedule for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
