// SPDX-License-Identifier: Apache-2.0

package types

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/constants"
)

func TestTypes_Events_Getters(t *testing.T) {
	// setup types
	eventsOne, eventsTwo := testEvents()

	// setup tests
	tests := []struct {
		events *Events
		want   *Events
	}{
		{
			events: eventsOne,
			want:   eventsOne,
		},
		{
			events: eventsTwo,
			want:   eventsTwo,
		},
		{
			events: new(Events),
			want:   new(Events),
		},
	}

	// run tests
	for _, test := range tests {
		if !reflect.DeepEqual(test.events.GetPush(), test.want.GetPush()) {
			t.Errorf("GetPush is %v, want %v", test.events.GetPush(), test.want.GetPush())
		}

		if !reflect.DeepEqual(test.events.GetPullRequest(), test.want.GetPullRequest()) {
			t.Errorf("GetPullRequest is %v, want %v", test.events.GetPush(), test.want.GetPush())
		}

		if !reflect.DeepEqual(test.events.GetDeployment(), test.want.GetDeployment()) {
			t.Errorf("GetDeployment is %v, want %v", test.events.GetPush(), test.want.GetPush())
		}

		if !reflect.DeepEqual(test.events.GetComment(), test.want.GetComment()) {
			t.Errorf("GetComment is %v, want %v", test.events.GetPush(), test.want.GetPush())
		}

		if !reflect.DeepEqual(test.events.GetSchedule(), test.want.GetSchedule()) {
			t.Errorf("GetSchedule is %v, want %v", test.events.GetSchedule(), test.want.GetSchedule())
		}
	}
}

func TestTypes_Events_Setters(t *testing.T) {
	// setup types
	var e *Events

	eventsOne, eventsTwo := testEvents()

	// setup tests
	tests := []struct {
		events *Events
		want   *Events
	}{
		{
			events: eventsOne,
			want:   eventsOne,
		},
		{
			events: eventsTwo,
			want:   eventsTwo,
		},
		{
			events: e,
			want:   new(Events),
		},
	}

	// run tests
	for _, test := range tests {
		test.events.SetPush(test.want.GetPush())
		test.events.SetPullRequest(test.want.GetPullRequest())
		test.events.SetDeployment(test.want.GetDeployment())
		test.events.SetComment(test.want.GetComment())
		test.events.SetSchedule(test.want.GetSchedule())

		if !reflect.DeepEqual(test.events.GetPush(), test.want.GetPush()) {
			t.Errorf("SetPush is %v, want %v", test.events.GetPush(), test.want.GetPush())
		}

		if !reflect.DeepEqual(test.events.GetPullRequest(), test.want.GetPullRequest()) {
			t.Errorf("SetPullRequest is %v, want %v", test.events.GetPullRequest(), test.want.GetPullRequest())
		}

		if !reflect.DeepEqual(test.events.GetDeployment(), test.want.GetDeployment()) {
			t.Errorf("SetDeployment is %v, want %v", test.events.GetDeployment(), test.want.GetDeployment())
		}

		if !reflect.DeepEqual(test.events.GetComment(), test.want.GetComment()) {
			t.Errorf("SetComment is %v, want %v", test.events.GetComment(), test.want.GetComment())
		}

		if !reflect.DeepEqual(test.events.GetSchedule(), test.want.GetSchedule()) {
			t.Errorf("SetSchedule is %v, want %v", test.events.GetSchedule(), test.want.GetSchedule())
		}
	}
}

func TestTypes_Events_List(t *testing.T) {
	// setup types
	eventsOne, eventsTwo := testEvents()

	wantOne := []string{
		"push",
		"pull_request:opened",
		"pull_request:synchronize",
		"pull_request:reopened",
		"pull_request:unlabeled",
		"tag",
		"comment:created",
		"schedule",
		"delete:branch",
	}

	wantTwo := []string{
		"pull_request:edited",
		"pull_request:labeled",
		"deployment",
		"comment:edited",
		"delete:tag",
	}

	// run test
	gotOne := eventsOne.List()

	if diff := cmp.Diff(wantOne, gotOne); diff != "" {
		t.Errorf("(List: -want +got):\n%s", diff)
	}

	gotTwo := eventsTwo.List()

	if diff := cmp.Diff(wantTwo, gotTwo); diff != "" {
		t.Errorf("(List Inverse: -want +got):\n%s", diff)
	}
}

func TestTypes_Events_NewEventsFromMask_ToDatabase(t *testing.T) {
	// setup mask
	maskOne := int64(
		constants.AllowPushBranch |
			constants.AllowPushTag |
			constants.AllowPushDeleteBranch |
			constants.AllowPullOpen |
			constants.AllowPullSync |
			constants.AllowPullReopen |
			constants.AllowPullUnlabel |
			constants.AllowCommentCreate |
			constants.AllowSchedule,
	)

	maskTwo := int64(
		constants.AllowPushDeleteTag |
			constants.AllowPullEdit |
			constants.AllowCommentEdit |
			constants.AllowPullLabel |
			constants.AllowDeployCreate,
	)

	wantOne, wantTwo := testEvents()

	// run test
	gotOne := NewEventsFromMask(maskOne)

	if diff := cmp.Diff(wantOne, gotOne); diff != "" {
		t.Errorf("(NewEventsFromMask: -want +got):\n%s", diff)
	}

	gotTwo := NewEventsFromMask(maskTwo)

	if diff := cmp.Diff(wantTwo, gotTwo); diff != "" {
		t.Errorf("(NewEventsFromMask Inverse: -want +got):\n%s", diff)
	}

	// ensure ToDatabase maps back to masks
	if gotOne.ToDatabase() != maskOne {
		t.Errorf("ToDatabase returned %d, want %d", gotOne.ToDatabase(), maskOne)
	}

	if gotTwo.ToDatabase() != maskTwo {
		t.Errorf("ToDatabase returned %d, want %d", gotTwo.ToDatabase(), maskTwo)
	}
}

func Test_NewEventsFromSlice(t *testing.T) {
	// setup types
	tBool := true
	fBool := false

	e1, e2 := testEvents()

	// setup tests
	tests := []struct {
		name    string
		events  []string
		want    *Events
		failure bool
	}{
		{
			name:    "action specific events to e1",
			events:  []string{"push:branch", "push:tag", "delete:branch", "pull_request:opened", "pull_request:synchronize", "pull_request:reopened", "comment:created", "schedule:run", "pull_request:unlabeled"},
			want:    e1,
			failure: false,
		},
		{
			name:    "action specific events to e2",
			events:  []string{"delete:tag", "pull_request:edited", "deployment:created", "comment:edited", "pull_request:labeled"},
			want:    e2,
			failure: false,
		},
		{
			name:   "general events",
			events: []string{"push", "pull", "deploy", "comment", "schedule", "tag", "delete"},
			want: &Events{
				Push: &actions.Push{
					Branch:       &tBool,
					Tag:          &tBool,
					DeleteBranch: &tBool,
					DeleteTag:    &tBool,
				},
				PullRequest: &actions.Pull{
					Opened:      &tBool,
					Reopened:    &tBool,
					Edited:      &fBool,
					Synchronize: &tBool,
					Labeled:     &fBool,
					Unlabeled:   &fBool,
				},
				Deployment: &actions.Deploy{
					Created: &tBool,
				},
				Comment: &actions.Comment{
					Created: &tBool,
					Edited:  &tBool,
				},
				Schedule: &actions.Schedule{
					Run: &tBool,
				},
			},
			failure: false,
		},
		{
			name:   "double events",
			events: []string{"push", "push:branch", "pull_request", "pull_request:opened"},
			want: &Events{
				Push: &actions.Push{
					Branch:       &tBool,
					Tag:          &fBool,
					DeleteBranch: &fBool,
					DeleteTag:    &fBool,
				},
				PullRequest: &actions.Pull{
					Opened:      &tBool,
					Reopened:    &tBool,
					Edited:      &fBool,
					Synchronize: &tBool,
					Labeled:     &fBool,
					Unlabeled:   &fBool,
				},
				Deployment: &actions.Deploy{
					Created: &fBool,
				},
				Comment: &actions.Comment{
					Created: &fBool,
					Edited:  &fBool,
				},
				Schedule: &actions.Schedule{
					Run: &fBool,
				},
			},
			failure: false,
		},
		{
			name:   "empty events",
			events: []string{},
			want:   NewEventsFromMask(0),
		},
		{
			name:    "invalid events",
			events:  []string{"foo:bar"},
			want:    nil,
			failure: true,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := NewEventsFromSlice(test.events)

		if test.failure {
			if err == nil {
				t.Errorf("NewEventsFromSlice should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("NewEventsFromSlice returned err: %v", err)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("PopulateEvents failed for %s mismatch (-want +got):\n%s", test.name, diff)
		}
	}
}

func TestTypes_Events_Allowed(t *testing.T) {
	// setup types
	eventsOne, eventsTwo := testEvents()

	// setup tests
	tests := []struct {
		event  string
		action string
		want   bool
	}{
		{event: "push", want: true},
		{event: "tag", want: true},
		{event: "pull_request", action: "opened", want: true},
		{event: "pull_request", action: "synchronize", want: true},
		{event: "pull_request", action: "edited", want: false},
		{event: "pull_request", action: "reopened", want: true},
		{event: "pull_request", action: "labeled", want: false},
		{event: "pull_request", action: "unlabeled", want: true},
		{event: "deployment", action: "created", want: false},
		{event: "comment", action: "created", want: true},
		{event: "comment", action: "edited", want: false},
		{event: "schedule", want: true},
		{event: "delete", action: "branch", want: true},
		{event: "delete", action: "tag", want: false},
	}

	for _, test := range tests {
		gotOne := eventsOne.Allowed(test.event, test.action)
		gotTwo := eventsTwo.Allowed(test.event, test.action)

		if gotOne != test.want {
			t.Errorf("Allowed for %s/%s is %v, want %v", test.event, test.action, gotOne, test.want)
		}

		if gotTwo == test.want {
			t.Errorf("Allowed Inverse for %s/%s is %v, want %v", test.event, test.action, gotTwo, !test.want)
		}
	}
}

// testEvents is a helper test function that returns an Events struct and its inverse for unit test coverage.
func testEvents() (*Events, *Events) {
	tBool := true
	fBool := false

	e1 := &Events{
		Push: &actions.Push{
			Branch:       &tBool,
			Tag:          &tBool,
			DeleteBranch: &tBool,
			DeleteTag:    &fBool,
		},
		PullRequest: &actions.Pull{
			Opened:      &tBool,
			Synchronize: &tBool,
			Edited:      &fBool,
			Reopened:    &tBool,
			Labeled:     &fBool,
			Unlabeled:   &tBool,
		},
		Deployment: &actions.Deploy{
			Created: &fBool,
		},
		Comment: &actions.Comment{
			Created: &tBool,
			Edited:  &fBool,
		},
		Schedule: &actions.Schedule{
			Run: &tBool,
		},
	}

	e2 := &Events{
		Push: &actions.Push{
			Branch:       &fBool,
			Tag:          &fBool,
			DeleteBranch: &fBool,
			DeleteTag:    &tBool,
		},
		PullRequest: &actions.Pull{
			Opened:      &fBool,
			Synchronize: &fBool,
			Edited:      &tBool,
			Reopened:    &fBool,
			Labeled:     &tBool,
			Unlabeled:   &fBool,
		},
		Deployment: &actions.Deploy{
			Created: &tBool,
		},
		Comment: &actions.Comment{
			Created: &fBool,
			Edited:  &tBool,
		},
		Schedule: &actions.Schedule{
			Run: &fBool,
		},
	}

	return e1, e2
}
