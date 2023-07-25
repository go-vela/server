// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_CountActiveSchedules(t *testing.T) {
	_scheduleOne := testSchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetRepoID(1)
	_scheduleOne.SetActive(true)
	_scheduleOne.SetName("nightly")
	_scheduleOne.SetEntry("0 0 * * *")
	_scheduleOne.SetCreatedAt(1)
	_scheduleOne.SetCreatedBy("user1")
	_scheduleOne.SetUpdatedAt(1)
	_scheduleOne.SetUpdatedBy("user2")

	_scheduleTwo := testSchedule()
	_scheduleTwo.SetID(2)
	_scheduleTwo.SetRepoID(2)
	_scheduleTwo.SetActive(false)
	_scheduleTwo.SetName("hourly")
	_scheduleTwo.SetEntry("0 * * * *")
	_scheduleTwo.SetCreatedAt(1)
	_scheduleTwo.SetCreatedBy("user1")
	_scheduleTwo.SetUpdatedAt(1)
	_scheduleTwo.SetUpdatedBy("user2")

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules" WHERE active = $1`).WithArgs(true).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSchedule(_scheduleOne)
	if err != nil {
		t.Errorf("unable to create test schedule for sqlite: %v", err)
	}

	_, err = _sqlite.CreateSchedule(_scheduleTwo)
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
			got, err := test.database.CountActiveSchedules()

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
