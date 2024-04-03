// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestSchedule_Engine_ListSchedulesForRepo(t *testing.T) {
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")

	_scheduleOne := testSchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetRepoID(1)
	_scheduleOne.SetName("nightly")
	_scheduleOne.SetEntry("0 0 * * *")
	_scheduleOne.SetCreatedAt(1)
	_scheduleOne.SetCreatedBy("user1")
	_scheduleOne.SetUpdatedAt(1)
	_scheduleOne.SetUpdatedBy("user2")
	_scheduleOne.SetBranch("main")

	_scheduleTwo := testSchedule()
	_scheduleTwo.SetID(2)
	_scheduleTwo.SetRepoID(2)
	_scheduleTwo.SetName("hourly")
	_scheduleTwo.SetEntry("0 * * * *")
	_scheduleTwo.SetCreatedAt(1)
	_scheduleTwo.SetCreatedBy("user1")
	_scheduleTwo.SetUpdatedAt(1)
	_scheduleTwo.SetUpdatedBy("user2")
	_scheduleTwo.SetBranch("main")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules" WHERE repo_id = $1`).WithArgs(1).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at", "branch"}).
		AddRow(1, 1, false, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil, "main")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules" WHERE repo_id = $1 ORDER BY id DESC LIMIT $2`).WithArgs(1, 10).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSchedule(context.TODO(), _scheduleOne)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSchedule(context.TODO(), _scheduleTwo)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*library.Schedule
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*library.Schedule{_scheduleOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*library.Schedule{_scheduleOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _, err := test.database.ListSchedulesForRepo(context.TODO(), _repo, 1, 10)

			if test.failure {
				if err == nil {
					t.Errorf("ListSchedulesForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSchedulesForRepo for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListSchedulesForRepo for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
