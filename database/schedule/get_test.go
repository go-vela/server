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

func TestSchedule_Engine_GetSchedule(t *testing.T) {
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

	_schedule := testAPISchedule()
	_schedule.SetID(1)
	_schedule.SetRepo(_repo)
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
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at", "branch"}).
		AddRow(1, 1, false, "nightly", "0 0 * * *", 1, "user1", 1, "user2", nil, "main")

	_repoRows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "baz", "foo", "bar", "foo/bar", "", "", "", "{}", 0, 0, 0, "public", false, false, false, 1, "yaml", "", "")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules" WHERE id = $1 LIMIT $2`).WithArgs(1, 1).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateSchedule(context.TODO(), _schedule)
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
		want     *api.Schedule
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
			got, err := test.database.GetSchedule(context.TODO(), 1)

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
