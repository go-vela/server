// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adhocore/gronx"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/server/database/types"
)

func TestSchedule_Engine_ListSchedules(t *testing.T) {
	// setup types
	_owner := testutils.APIUser()
	_owner.SetID(1)
	_owner.SetName("octocat")
	_owner.SetToken("superSecretToken")
	_owner.SetRefreshToken("superSecretRefreshToken")
	_owner.SetFavorites([]string{"github/octocat"})
	_owner.SetActive(true)
	_owner.SetAdmin(false)
	_owner.SetDashboards([]string{"45bcf19b-c151-4e2d-b8c6-80a62ba2eae7"})

	_repo := testutils.APIRepo()
	_repo.SetID(1)
	_repo.SetOwner(_owner.Crop())
	_repo.SetHash("MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy")
	_repo.SetOrg("github")
	_repo.SetName("octocat")
	_repo.SetFullName("github/octocat")
	_repo.SetLink("https://github.com/github/octocat")
	_repo.SetClone("https://github.com/github/octocat.git")
	_repo.SetBranch("main")
	_repo.SetTopics([]string{"cloud", "security"})
	_repo.SetBuildLimit(10)
	_repo.SetTimeout(30)
	_repo.SetCounter(0)
	_repo.SetVisibility("public")
	_repo.SetPrivate(false)
	_repo.SetTrusted(false)
	_repo.SetActive(true)
	_repo.SetAllowEvents(api.NewEventsFromMask(1))
	_repo.SetPipelineType("")
	_repo.SetPreviousName("")
	_repo.SetApproveBuild(constants.ApproveNever)
	_repo.SetInstallID(0)

	currTime := time.Now().UTC()
	nextTime, _ := gronx.NextTickAfter("0 0 * * *", currTime, false)

	_scheduleOne := testutils.APISchedule()
	_scheduleOne.SetID(1)
	_scheduleOne.SetRepo(_repo)
	_scheduleOne.SetActive(true)
	_scheduleOne.SetName("nightly")
	_scheduleOne.SetEntry("0 0 * * *")
	_scheduleOne.SetCreatedAt(1713476291)
	_scheduleOne.SetCreatedBy("octocat")
	_scheduleOne.SetUpdatedAt(3013476291)
	_scheduleOne.SetUpdatedBy("octokitty")
	_scheduleOne.SetScheduledAt(2013476291)
	_scheduleOne.SetBranch("main")
	_scheduleOne.SetError("no version: YAML property provided")
	_scheduleOne.SetNextRun(nextTime.Unix())

	currTime = time.Now().UTC()
	nextTime, _ = gronx.NextTickAfter("0 * * * *", currTime, false)

	_scheduleTwo := testutils.APISchedule()
	_scheduleTwo.SetID(2)
	_scheduleTwo.SetRepo(_repo)
	_scheduleTwo.SetActive(false)
	_scheduleTwo.SetName("hourly")
	_scheduleTwo.SetEntry("0 * * * *")
	_scheduleTwo.SetCreatedAt(1713476291)
	_scheduleTwo.SetCreatedBy("octocat")
	_scheduleTwo.SetUpdatedAt(3013476291)
	_scheduleTwo.SetUpdatedBy("octokitty")
	_scheduleTwo.SetScheduledAt(2013476291)
	_scheduleTwo.SetBranch("main")
	_scheduleTwo.SetError("no version: YAML property provided")
	_scheduleTwo.SetNextRun(nextTime.Unix())

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows([]string{"count"}).AddRow(2)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT count(*) FROM "schedules"`).WillReturnRows(_rows)

	// create expected result in mock
	_rows = sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at", "branch", "error"}).
		AddRow(1, 1, true, "nightly", "0 0 * * *", 1713476291, "octocat", 3013476291, "octokitty", 2013476291, "main", "no version: YAML property provided").
		AddRow(2, 1, false, "hourly", "0 * * * *", 1713476291, "octocat", 3013476291, "octokitty", 2013476291, "main", "no version: YAML property provided")

	_repoRows := sqlmock.NewRows(
		[]string{"id", "user_id", "hash", "org", "name", "full_name", "link", "clone", "branch", "topics", "build_limit", "timeout", "counter", "visibility", "private", "trusted", "active", "allow_events", "pipeline_type", "previous_name", "approve_build"}).
		AddRow(1, 1, "MzM4N2MzMDAtNmY4Mi00OTA5LWFhZDAtNWIzMTlkNTJkODMy", "github", "octocat", "github/octocat", "https://github.com/github/octocat", "https://github.com/github/octocat.git", "main", "{cloud,security}", 10, 30, 0, "public", false, false, true, 1, "", "", constants.ApproveNever)

	_userRows := sqlmock.NewRows(
		[]string{"id", "name", "token", "refresh_token", "favorites", "active", "admin", "dashboards"}).
		AddRow(1, "octocat", "superSecretToken", "superSecretRefreshToken", "{}", true, false, "{}")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules"`).WillReturnRows(_rows)
	_mock.ExpectQuery(`SELECT * FROM "repos" WHERE "repos"."id" = $1`).WithArgs(1).WillReturnRows(_repoRows)
	_mock.ExpectQuery(`SELECT * FROM "users" WHERE "users"."id" = $1`).WithArgs(1).WillReturnRows(_userRows)

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

	err = _sqlite.client.AutoMigrate(&types.Repo{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableRepo).Create(types.RepoFromAPI(_repo)).Error
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	err = _sqlite.client.AutoMigrate(&types.User{})
	if err != nil {
		t.Errorf("unable to create build table for sqlite: %v", err)
	}

	err = _sqlite.client.Table(constants.TableUser).Create(types.UserFromAPI(_owner)).Error
	if err != nil {
		t.Errorf("unable to create test user for sqlite: %v", err)
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
			want:     []*api.Schedule{_scheduleOne, _scheduleTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []*api.Schedule{_scheduleOne, _scheduleTwo},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListSchedules(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListSchedules for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListSchedules for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ListSchedules for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
