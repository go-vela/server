// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestLibrary_Schedule_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *Schedule
		want    *Schedule
	}{
		{
			actions: testSchedule(),
			want:    testSchedule(),
		},
		{
			actions: new(Schedule),
			want:    new(Schedule),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetRun() != test.want.GetRun() {
			t.Errorf("GetRun is %v, want %v", test.actions.GetRun(), test.want.GetRun())
		}
	}
}

func TestLibrary_Schedule_Setters(t *testing.T) {
	// setup types
	var a *Schedule

	// setup tests
	tests := []struct {
		actions *Schedule
		want    *Schedule
	}{
		{
			actions: testSchedule(),
			want:    testSchedule(),
		},
		{
			actions: a,
			want:    new(Schedule),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetRun(test.want.GetRun())

		if test.actions.GetRun() != test.want.GetRun() {
			t.Errorf("SetRun is %v, want %v", test.actions.GetRun(), test.want.GetRun())
		}
	}
}

func TestLibrary_Schedule_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testSchedule()

	// run test
	got := new(Schedule).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestLibrary_Schedule_ToMask(t *testing.T) {
	// setup types
	actions := testSchedule()

	want := int64(constants.AllowSchedule)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testSchedule() *Schedule {
	schedule := new(Schedule)
	schedule.SetRun(true)

	return schedule
}
