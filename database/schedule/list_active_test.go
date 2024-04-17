// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/repo"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
)

func TestSchedule_Engine_ListActiveSchedules(t *testing.T) {
	_repo := testRepo()
	_repo.SetID(1)
	_repo.SetHash("baz")
	_repo.SetOrg("foo")
	_repo.SetName("bar")
	_repo.SetFullName("foo/bar")
	_repo.SetVisibility("public")
	_repo.SetPipelineType("yaml")
	_repo.SetTopics([]string{})
	_repo.SetAllowEvents(api.NewEventsFromMask(1))

	_scheduleOne := testAPISchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetRepo(_repo)
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
	_scheduleTwo.SetRepo(_repo)
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
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules" WHERE active = $1`).WithArgs(true).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at", "branch"}).
		AddRow(1, 1, true, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil, "main")

	_repoRows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules" WHERE active = $1`).WithArgs(true).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)

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

	err = _sqlite.client.AutoMigrate(&database.Repo{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(repo.FromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []*api.Schedule
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []*api.Schedule{_scheduleOne},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Schedule{_scheduleOne},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListActiveSchedules(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListActiveSchedules for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListActiveSchedules for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListActiveSchedules for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
