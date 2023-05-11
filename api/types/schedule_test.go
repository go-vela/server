// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package types

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-vela/types/library"
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
			if !reflect.DeepEqual(test.schedule.GetRepo(), test.want.GetRepo()) {
				t.Errorf("GetRepo is %v, want %v", test.schedule.GetRepo(), test.want.GetRepo())
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
			if test.schedule.GetID() != test.want.GetID() {
				t.Errorf("SetID is %v, want %v", test.schedule.GetID(), test.want.GetID())
			}
			test.schedule.SetActive(test.want.GetActive())
			if test.schedule.GetActive() != test.want.GetActive() {
				t.Errorf("SetActive is %v, want %v", test.schedule.GetActive(), test.want.GetActive())
			}
			test.schedule.SetName(test.want.GetName())
			if test.schedule.GetName() != test.want.GetName() {
				t.Errorf("SetName is %v, want %v", test.schedule.GetName(), test.want.GetName())
			}
			test.schedule.SetEntry(test.want.GetEntry())
			if test.schedule.GetEntry() != test.want.GetEntry() {
				t.Errorf("SetEntry is %v, want %v", test.schedule.GetEntry(), test.want.GetEntry())
			}
			test.schedule.SetCreatedAt(test.want.GetCreatedAt())
			if test.schedule.GetCreatedAt() != test.want.GetCreatedAt() {
				t.Errorf("SetCreatedAt is %v, want %v", test.schedule.GetCreatedAt(), test.want.GetCreatedAt())
			}
			test.schedule.SetCreatedBy(test.want.GetCreatedBy())
			if test.schedule.GetCreatedBy() != test.want.GetCreatedBy() {
				t.Errorf("SetCreatedBy is %v, want %v", test.schedule.GetCreatedBy(), test.want.GetCreatedBy())
			}
			test.schedule.SetUpdatedAt(test.want.GetUpdatedAt())
			if test.schedule.GetUpdatedAt() != test.want.GetUpdatedAt() {
				t.Errorf("SetUpdatedAt is %v, want %v", test.schedule.GetUpdatedAt(), test.want.GetUpdatedAt())
			}
			test.schedule.SetUpdatedBy(test.want.GetUpdatedBy())
			if test.schedule.GetUpdatedBy() != test.want.GetUpdatedBy() {
				t.Errorf("SetUpdatedBy is %v, want %v", test.schedule.GetUpdatedBy(), test.want.GetUpdatedBy())
			}
			test.schedule.SetScheduledAt(test.want.GetScheduledAt())
			if test.schedule.GetScheduledAt() != test.want.GetScheduledAt() {
				t.Errorf("SetScheduledAt is %v, want %v", test.schedule.GetScheduledAt(), test.want.GetScheduledAt())
			}
			test.schedule.SetRepo(test.want.GetRepo())
			if !reflect.DeepEqual(test.schedule.GetRepo(), test.want.GetRepo()) {
				t.Errorf("SetRepo is %v, want %v", test.schedule.GetRepo(), test.want.GetRepo())
			}
		})
	}
}

func TestLibrary_Schedule_String(t *testing.T) {
	s := testSchedule()

	want := fmt.Sprintf(`{
  Active: %t,
  CreatedAt: %d,
  CreatedBy: %s,
  Entry: %s,
  ID: %d,
  Name: %s,
  ScheduledAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Repo: %v,
}`,
		s.GetActive(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetEntry(),
		s.GetID(),
		s.GetName(),
		s.GetScheduledAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
		s.GetRepo(),
	)

	got := s.String()
	if !strings.EqualFold(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testSchedule is a test helper function to create a Schedule type with all fields set to a fake value.
func testSchedule() *Schedule {
	r := new(library.Repo)
	r.SetID(1)

	s := new(Schedule)
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

	return s
}
