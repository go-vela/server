// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package types

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/go-vela/server/api/types"

	"github.com/go-vela/types/library"
)

func TestTypes_ScheduleFromAPI(t *testing.T) {
	r := new(library.Repo)
	r.SetID(1)

	s := new(types.Schedule)
	s.SetID(1)
	s.SetActive(true)
	s.SetName("nightly")
	s.SetEntry("0 0 * * *")
	s.SetCreatedAt(time.Now().UTC().Unix())
	s.SetCreatedBy("user1")
	s.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	s.SetUpdatedBy("user2")
	s.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	s.SetRepo(r)

	want := testSchedule()

	got := ScheduleFromAPI(s)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ScheduleFromAPI is %v, want %v", got, want)
	}
}

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
	r := new(library.Repo)
	r.SetID(1)

	want := new(types.Schedule)
	want.SetID(1)
	want.SetActive(true)
	want.SetName("nightly")
	want.SetEntry("0 0 * * *")
	want.SetCreatedAt(time.Now().UTC().Unix())
	want.SetCreatedBy("user1")
	want.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	want.SetUpdatedBy("user2")
	want.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	want.SetRepo(r)

	got := testSchedule().ToAPI(r)
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
	}
}
