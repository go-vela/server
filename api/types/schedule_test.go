// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTypes_Schedule_Getters(t *testing.T) {
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
			want:     new(Schedule),
		},
		{
			name:     "empty schedule",
			schedule: nil,
			want:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.schedule.GetID() != test.want.GetID() {
				t.Errorf("GetID is %v, want %v", test.schedule.GetID(), test.want.GetID())
			}

			if test.schedule.GetRepo().GetID() != test.want.GetRepo().GetID() {
				t.Errorf("GetRepoID is %v, want %v", test.schedule.GetRepo().GetID(), test.want.GetRepo().GetID())
			}

			if test.schedule.GetActive() != test.want.GetActive() {
				t.Errorf("GetActive is %v, want %v", test.schedule.GetActive(), test.want.GetActive())
			}

			if test.schedule.GetName() != test.want.GetName() {
				t.Errorf("GetName is %v, want %v", test.schedule.GetName(), test.want.GetName())
			}

			if test.schedule.GetEntry() != test.want.GetEntry() {
				t.Errorf("GetEntry is %v, want %v", test.schedule.GetEntry(), test.want.GetEntry())
			}

			if test.schedule.GetCreatedAt() != test.want.GetCreatedAt() {
				t.Errorf("GetCreatedAt is %v, want %v", test.schedule.GetCreatedAt(), test.want.GetCreatedAt())
			}

			if test.schedule.GetCreatedBy() != test.want.GetCreatedBy() {
				t.Errorf("GetCreatedBy is %v, want %v", test.schedule.GetCreatedBy(), test.want.GetCreatedBy())
			}

			if test.schedule.GetUpdatedAt() != test.want.GetUpdatedAt() {
				t.Errorf("GetUpdatedAt is %v, want %v", test.schedule.GetUpdatedAt(), test.want.GetUpdatedAt())
			}

			if test.schedule.GetUpdatedBy() != test.want.GetUpdatedBy() {
				t.Errorf("GetUpdatedBy is %v, want %v", test.schedule.GetUpdatedBy(), test.want.GetUpdatedBy())
			}

			if test.schedule.GetScheduledAt() != test.want.GetScheduledAt() {
				t.Errorf("GetScheduledAt is %v, want %v", test.schedule.GetScheduledAt(), test.want.GetScheduledAt())
			}

			if test.schedule.GetBranch() != test.want.GetBranch() {
				t.Errorf("GetBranch is %v, want %v", test.schedule.GetBranch(), test.want.GetBranch())
			}

			if test.schedule.GetError() != test.want.GetError() {
				t.Errorf("GetError is %v, want %v", test.schedule.GetError(), test.want.GetError())
			}
		})
	}
}

func TestTypes_Schedule_Setters(t *testing.T) {
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
			want:     new(Schedule),
		},
		{
			name:     "empty schedule",
			schedule: nil,
			want:     nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.schedule.SetID(test.want.GetID())
			test.schedule.SetRepo(test.want.GetRepo())
			test.schedule.SetActive(test.want.GetActive())
			test.schedule.SetName(test.want.GetName())
			test.schedule.SetEntry(test.want.GetEntry())
			test.schedule.SetCreatedAt(test.want.GetCreatedAt())
			test.schedule.SetCreatedBy(test.want.GetCreatedBy())
			test.schedule.SetUpdatedAt(test.want.GetUpdatedAt())
			test.schedule.SetUpdatedBy(test.want.GetUpdatedBy())
			test.schedule.SetScheduledAt(test.want.GetScheduledAt())
			test.schedule.SetBranch(test.want.GetBranch())
			test.schedule.SetError(test.want.GetError())

			if test.schedule.GetID() != test.want.GetID() {
				t.Errorf("SetID is %v, want %v", test.schedule.GetID(), test.want.GetID())
			}

			if test.schedule.GetRepo().GetID() != test.want.GetRepo().GetID() {
				t.Errorf("SetRepoID is %v, want %v", test.schedule.GetRepo().GetID(), test.want.GetRepo().GetID())
			}

			if test.schedule.GetActive() != test.want.GetActive() {
				t.Errorf("SetActive is %v, want %v", test.schedule.GetActive(), test.want.GetActive())
			}

			if test.schedule.GetName() != test.want.GetName() {
				t.Errorf("SetName is %v, want %v", test.schedule.GetName(), test.want.GetName())
			}

			if test.schedule.GetEntry() != test.want.GetEntry() {
				t.Errorf("SetEntry is %v, want %v", test.schedule.GetEntry(), test.want.GetEntry())
			}

			if test.schedule.GetCreatedAt() != test.want.GetCreatedAt() {
				t.Errorf("SetCreatedAt is %v, want %v", test.schedule.GetCreatedAt(), test.want.GetCreatedAt())
			}

			if test.schedule.GetCreatedBy() != test.want.GetCreatedBy() {
				t.Errorf("SetCreatedBy is %v, want %v", test.schedule.GetCreatedBy(), test.want.GetCreatedBy())
			}

			if test.schedule.GetUpdatedAt() != test.want.GetUpdatedAt() {
				t.Errorf("SetUpdatedAt is %v, want %v", test.schedule.GetUpdatedAt(), test.want.GetUpdatedAt())
			}

			if test.schedule.GetUpdatedBy() != test.want.GetUpdatedBy() {
				t.Errorf("SetUpdatedBy is %v, want %v", test.schedule.GetUpdatedBy(), test.want.GetUpdatedBy())
			}

			if test.schedule.GetScheduledAt() != test.want.GetScheduledAt() {
				t.Errorf("SetScheduledAt is %v, want %v", test.schedule.GetScheduledAt(), test.want.GetScheduledAt())
			}

			if test.schedule.GetBranch() != test.want.GetBranch() {
				t.Errorf("SetBranch is %v, want %v", test.schedule.GetBranch(), test.want.GetBranch())
			}

			if test.schedule.GetError() != test.want.GetError() {
				t.Errorf("SetError is %v, want %v", test.schedule.GetError(), test.want.GetError())
			}
		})
	}
}

func TestTypes_Schedule_String(t *testing.T) {
	s := testSchedule()

	want := fmt.Sprintf(`{
  Active: %t,
  CreatedAt: %d,
  CreatedBy: %s,
  Entry: %s,
  ID: %d,
  Name: %s,
  Repo: %v,
  ScheduledAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Branch: %s,
  Error: %s,
}`,
		s.GetActive(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetEntry(),
		s.GetID(),
		s.GetName(),
		s.GetRepo(),
		s.GetScheduledAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
		s.GetBranch(),
		s.GetError(),
	)

	got := s.String()
	if !strings.EqualFold(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testSchedule is a test helper function to create a Schedule type with all fields set to a fake value.
func testSchedule() *Schedule {
	s := new(Schedule)
	s.SetID(1)
	s.SetRepo(testRepo())
	s.SetActive(true)
	s.SetName("nightly")
	s.SetEntry("0 0 * * *")
	s.SetCreatedAt(time.Now().UTC().Unix())
	s.SetCreatedBy("user1")
	s.SetUpdatedAt(time.Now().Add(time.Hour * 1).UTC().Unix())
	s.SetUpdatedBy("user2")
	s.SetScheduledAt(time.Now().Add(time.Hour * 2).UTC().Unix())
	s.SetBranch("main")
	s.SetError("unable to trigger build for schedule nightly: unknown character")

	return s
}
