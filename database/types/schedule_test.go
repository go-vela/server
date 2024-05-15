// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
	"github.com/go-vela/types/constants"
)

func TestTypes_Schedule_Nullify(t *testing.T) {
	tests := []struct {
		name     string
		schedule *Schedule
		want     *Schedule
	}{
		{
			name:     "schedule with fields",
			schedule: testSchedule(),
			want:     testSchedule(),
		},
		{
			name:     "schedule with empty fields",
			schedule: new(Schedule),
			want: &Schedule{
				ID:          sql.NullInt64{Int64: 0, Valid: false},
				RepoID:      sql.NullInt64{Int64: 0, Valid: false},
				Active:      sql.NullBool{Bool: false, Valid: false},
				Name:        sql.NullString{String: "", Valid: false},
				Entry:       sql.NullString{String: "", Valid: false},
				CreatedAt:   sql.NullInt64{Int64: 0, Valid: false},
				CreatedBy:   sql.NullString{String: "", Valid: false},
				UpdatedAt:   sql.NullInt64{Int64: 0, Valid: false},
				UpdatedBy:   sql.NullString{String: "", Valid: false},
				ScheduledAt: sql.NullInt64{Int64: 0, Valid: false},
				Branch:      sql.NullString{String: "", Valid: false},
				Error:       sql.NullString{String: "", Valid: false},
			},
		},
		{
			name:     "empty schedule",
			schedule: nil,
			want:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.schedule.Nullify()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Nullify is %v, want %v", got, test.want)
			}
		})
	}
}

func TestTypes_Schedule_ToAPI(t *testing.T) {
	// setup types
	e := api.NewEventsFromMask(1)

	owner := testutils.APIUser().Crop()
	owner.SetID(1)
	owner.SetName("octocat")
	owner.SetActive(true)
	owner.SetToken("superSecretToken")
	owner.SetRefreshToken("superSecretRefreshToken")

	repo := testutils.APIRepo()
	repo.SetID(1)
	repo.SetOwner(owner)
	repo.SetHash("superSecretHash")
	repo.SetOrg("github")
	repo.SetName("octocat")
	repo.SetFullName("github/octocat")
	repo.SetLink("https://github.com/github/octocat")
	repo.SetClone("https://github.com/github/octocat.git")
	repo.SetBranch("main")
	repo.SetTopics([]string{"cloud", "security"})
	repo.SetBuildLimit(10)
	repo.SetTimeout(30)
	repo.SetCounter(0)
	repo.SetVisibility("public")
	repo.SetPrivate(false)
	repo.SetTrusted(false)
	repo.SetActive(true)
	repo.SetAllowEvents(e)
	repo.SetPipelineType("yaml")
	repo.SetPreviousName("oldName")
	repo.SetApproveBuild(constants.ApproveNever)

	want := testutils.APISchedule()
	want.SetID(1)
	want.SetActive(true)
	want.SetRepo(repo)
	want.SetName("nightly")
	want.SetEntry("0 0 * * *")
	want.SetCreatedAt(time.Now().UTC().Unix())
	want.SetCreatedBy("user1")
	want.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	want.SetUpdatedBy("user2")
	want.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	want.SetBranch("main")
	want.SetError("unable to trigger build for schedule nightly: unknown character")
	want.SetNextRun(0)

	// run test
	got := testSchedule().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_Schedule_Validate(t *testing.T) {
	tests := []struct {
		name     string
		failure  bool
		schedule *Schedule
	}{
		{
			name:     "schedule with valid fields",
			failure:  false,
			schedule: testSchedule(),
		},
		{
			name:    "schedule with invalid entry",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Name:   sql.NullString{String: "invalid", Valid: false},
				Entry:  sql.NullString{String: "!@#$%^&*()", Valid: false},
			},
		},
		{
			name:    "schedule with missing entry",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Name:   sql.NullString{String: "nightly", Valid: false},
			},
		},
		{
			name:    "schedule with missing name",
			failure: true,
			schedule: &Schedule{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Entry:  sql.NullString{String: "0 0 * * *", Valid: false},
			},
		},
		{
			name:    "schedule with missing repo_id",
			failure: true,
			schedule: &Schedule{
				ID:    sql.NullInt64{Int64: 1, Valid: true},
				Name:  sql.NullString{String: "nightly", Valid: false},
				Entry: sql.NullString{String: "0 0 * * *", Valid: false},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.schedule.Validate()
			if test.failure {
				if err == nil {
					t.Errorf("Validate should have returned err")
				}

				return
			}

			if err != nil {
				t.Errorf("Validate returned err: %v", err)
			}
		})
	}
}

func TestTypes_ScheduleFromAPI(t *testing.T) {
	// setup types
	s := new(api.Schedule)
	repo := testutils.APIRepo()
	repo.SetID(1)

	s.SetID(1)
	s.SetActive(true)
	s.SetRepo(repo)
	s.SetName("nightly")
	s.SetEntry("0 0 * * *")
	s.SetCreatedAt(time.Now().UTC().Unix())
	s.SetCreatedBy("user1")
	s.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	s.SetUpdatedBy("user2")
	s.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	s.SetBranch("main")
	s.SetError("unable to trigger build for schedule nightly: unknown character")

	want := testSchedule()
	want.Repo = Repo{}

	// run test
	got := ScheduleFromAPI(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ScheduleFromAPI is %v, want %v", got, want)
	}
}

// testSchedule is a test helper function to create a Schedule type with all fields set to a fake value.
func testSchedule() *Schedule {
	return &Schedule{
		ID:          sql.NullInt64{Int64: 1, Valid: true},
		RepoID:      sql.NullInt64{Int64: 1, Valid: true},
		Active:      sql.NullBool{Bool: true, Valid: true},
		Name:        sql.NullString{String: "nightly", Valid: true},
		Entry:       sql.NullString{String: "0 0 * * *", Valid: true},
		CreatedAt:   sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
		CreatedBy:   sql.NullString{String: "user1", Valid: true},
		UpdatedAt:   sql.NullInt64{Int64: time.Now().Add(time.Hour * 1).UTC().Unix(), Valid: true},
		UpdatedBy:   sql.NullString{String: "user2", Valid: true},
		ScheduledAt: sql.NullInt64{Int64: time.Now().Add(time.Hour * 2).UTC().Unix(), Valid: true},
		Branch:      sql.NullString{String: "main", Valid: true},
		Error:       sql.NullString{String: "unable to trigger build for schedule nightly: unknown character", Valid: true},

		Repo: *testRepo(),
	}
}
