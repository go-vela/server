// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_CountActiveSchedules(t *testing.T) {
	_scheduleOne := testAPISchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetRepoID(1)
	_scheduleOne.SetActive(true)
	_scheduleOne.SetName("nightly")
	_scheduleOne.SetEntry("0 0 * * *")
	_scheduleOne.SetCreatedAt(1)
	_scheduleOne.SetCreatedBy("user1")
	_scheduleOne.SetUpdatedAt(1)
	_scheduleOne.SetUpdatedBy("user2")
	_scheduleOne.SetBranch("main")

	_scheduleTwo := testAPISchedule()
	_scheduleTwo.SetID(2)
	_scheduleTwo.SetRepoID(2)
	_scheduleTwo.SetActive(false)
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
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules" WHERE active = $1`).WithArgs(true).WillReturnRows(_rows)

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
		want     int64
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     1,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     1,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.CountActiveSchedules(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("CountActiveSchedules for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CountActiveSchedules for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("CountActiveSchedules for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
