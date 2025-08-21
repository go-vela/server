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
)

func TestSchedule_Engine_UpdateSchedule_Config(t *testing.T) {
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

	currTime := time.Now().UTC()
	nextTime, _ := gronx.NextTickAfter("0 0 * * *", currTime, false)

	_schedule := testutils.APISchedule()
	_schedule.SetID(1)
	_schedule.SetRepo(_repo)
	_schedule.SetActive(true)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1713476291)
	_schedule.SetCreatedBy("octocat")
	_schedule.SetUpdatedAt(3013476291)
	_schedule.SetUpdatedBy("octokitty")
	_schedule.SetScheduledAt(2013476291)
	_schedule.SetBranch("main")
	_schedule.SetError("no version: YAML property provided")
	_schedule.SetNextRun(nextTime.Unix())

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "schedules"
SET "repo_id"=$1,"active"=$2,"name"=$3,"entry"=$4,"created_at"=$5,"created_by"=$6,"updated_at"=$7,"updated_by"=$8,"scheduled_at"=$9,"branch"=$10,"error"=$11
WHERE "id" = $12`).
		WithArgs(1, true, "nightly", "0 0 * * *", 1713476291, "octocat", NowTimestamp{}, "octokitty", 2013476291, "main", "no version: YAML property provided", 1).
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
		database *Engine
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

			if diff := cmp.Diff(_schedule, got); diff != "" {
				t.Errorf("UpdateSchedule for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}

func TestSchedule_Engine_UpdateSchedule_NotConfig(t *testing.T) {
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

	currTime := time.Now().UTC()
	nextTime, _ := gronx.NextTickAfter("0 0 * * *", currTime, false)

	_schedule := testutils.APISchedule()
	_schedule.SetID(1)
	_schedule.SetRepo(_repo)
	_schedule.SetActive(true)
	_schedule.SetName("nightly")
	_schedule.SetEntry("0 0 * * *")
	_schedule.SetCreatedAt(1713476291)
	_schedule.SetCreatedBy("octocat")
	_schedule.SetUpdatedAt(3013476291)
	_schedule.SetUpdatedBy("octokitty")
	_schedule.SetScheduledAt(2013476291)
	_schedule.SetBranch("main")
	_schedule.SetError("no version: YAML property provided")
	_schedule.SetNextRun(nextTime.Unix())

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "schedules" SET "scheduled_at"=$1 WHERE "id" = $2`).
		WithArgs(2013476291, 1).
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
		database *Engine
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

			if diff := cmp.Diff(_schedule, got); diff != "" {
				t.Errorf("UpdateSchedule for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
