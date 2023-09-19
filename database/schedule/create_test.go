// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSchedule_Engine_CreateSchedule(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "schedules"
("repo_id","active","name","entry","created_at","created_by","updated_at","updated_by","scheduled_at","branch","id")
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING "id"`).
		WithArgs(1, false, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil, "main", 1).
		WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

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
			got, err := test.database.CreateSchedule(context.TODO(), _schedule)

			if test.failure {
				if err == nil {
					t.Errorf("CreateSchedule for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateSchedule for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _schedule) {
				t.Errorf("CreateSchedule for %s returned %s, want %s", test.name, got, _schedule)
			}
		})
	}
}
