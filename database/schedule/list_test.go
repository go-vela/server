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

func TestSchedule_Engine_ListSchedules(t *testing.T) {
	_repoOne := testRepo()
	_repoOne.SetID(1)
	_repoOne.SetOrg("foo")
	_repoOne.SetName("bar")
	_repoOne.SetFullName("foo/bar")

	_repoTwo := testRepo()
	_repoTwo.SetID(2)
	_repoTwo.SetOrg("bar")
	_repoTwo.SetName("foo")
	_repoTwo.SetFullName("bar/foo")

	_scheduleOne := testSchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetName("nightly")
	_scheduleOne.SetEntry("0 0 * * *")
	_scheduleOne.SetCreatedAt(1)
	_scheduleOne.SetCreatedBy("user1")
	_scheduleOne.SetUpdatedAt(1)
	_scheduleOne.SetUpdatedBy("user2")
	_scheduleOne.SetRepo(_repoOne)

	_scheduleTwo := testSchedule()
	_scheduleTwo.SetID(2)
	_scheduleTwo.SetName("hourly")
	_scheduleTwo.SetEntry("0 * * * *")
	_scheduleTwo.SetCreatedAt(1)
	_scheduleTwo.SetCreatedBy("user1")
	_scheduleTwo.SetUpdatedAt(1)
	_scheduleTwo.SetUpdatedBy("user2")
	_scheduleTwo.SetRepo(_repoTwo)

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at"}).
		AddRow(1, 1, false, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil).
		AddRow(2, 1, false, "hourly", "0 * * * *", 1, "user1", 1, "user2", nil)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateSchedule(_scheduleOne)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	err = _sqlite.CreateSchedule(_scheduleTwo)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	_scheduleOne.SetRepo(nil)
	_scheduleTwo.SetRepo(nil)

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*types.Schedule
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*types.Schedule{_scheduleOne, _scheduleTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*types.Schedule{_scheduleOne, _scheduleTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSchedules()

			if test.failure {
				if err == nil {
					t.Errorf("ListSchedules for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSchedules for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSchedules for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
