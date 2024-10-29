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

func TestSchedule_Engine_GetScheduleForRepo(t *testing.T) {
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

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "repo_id", "active", "name", "entry", "created_at", "created_by", "updated_at", "updated_by", "scheduled_at", "branch", "error"}).
		AddRow(1, 1, true, "nightly", "0 0 * * *", 1713476291, "octocat", 3013476291, "octokitty", 2013476291, "main", "no version: YAML property provided")

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "schedules" WHERE repo_id = $1 AND name = $2 LIMIT $3`).WithArgs(1, "nightly", 1).WillReturnRows(_rows)

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
			got, err := test.database.GetScheduleForRepo(context.TODO(), _repo, "nightly")

			if test.failure {
				if err == nil {
					t.Errorf("GetScheduleForRepo for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetScheduleForRepo for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetScheduleForRepo for %s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}
