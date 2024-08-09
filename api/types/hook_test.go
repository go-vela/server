// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTypes_Hook_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		hook *Hook
		want *Hook
	}{
		{
			hook: testHook(),
			want: testHook(),
		},
		{
			hook: new(Hook),
			want: new(Hook),
		},
	}

	// run tests
	for _, test := range tests {
		if test.hook.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.hook.GetID(), test.want.GetID())
		}

		if diff := cmp.Diff(test.hook.GetRepo(), test.want.GetRepo()); diff != "" {
			t.Errorf("GetRepo is a mismatch (-got +want):\n%s", diff)
		}

		if diff := cmp.Diff(test.hook.GetBuild(), test.want.GetBuild()); diff != "" {
			t.Errorf("GetBuild is a mismatch (-got +want):\n%s", diff)
		}

		if test.hook.GetNumber() != test.want.GetNumber() {
			t.Errorf("GetNumber is %v, want %v", test.hook.GetNumber(), test.want.GetNumber())
		}

		if test.hook.GetSourceID() != test.want.GetSourceID() {
			t.Errorf("GetSourceID is %v, want %v", test.hook.GetSourceID(), test.want.GetSourceID())
		}

		if test.hook.GetCreated() != test.want.GetCreated() {
			t.Errorf("GetCreated is %v, want %v", test.hook.GetCreated(), test.want.GetCreated())
		}

		if test.hook.GetHost() != test.want.GetHost() {
			t.Errorf("GetHost is %v, want %v", test.hook.GetHost(), test.want.GetHost())
		}

		if test.hook.GetEvent() != test.want.GetEvent() {
			t.Errorf("GetEvent is %v, want %v", test.hook.GetEvent(), test.want.GetEvent())
		}

		if test.hook.GetEventAction() != test.want.GetEventAction() {
			t.Errorf("GetEventAction is %v, want %v", test.hook.GetEventAction(), test.want.GetEventAction())
		}

		if test.hook.GetBranch() != test.want.GetBranch() {
			t.Errorf("GetBranch is %v, want %v", test.hook.GetBranch(), test.want.GetBranch())
		}

		if test.hook.GetError() != test.want.GetError() {
			t.Errorf("GetError is %v, want %v", test.hook.GetError(), test.want.GetError())
		}

		if test.hook.GetStatus() != test.want.GetStatus() {
			t.Errorf("GetStatus is %v, want %v", test.hook.GetStatus(), test.want.GetStatus())
		}

		if test.hook.GetLink() != test.want.GetLink() {
			t.Errorf("GetLink is %v, want %v", test.hook.GetLink(), test.want.GetLink())
		}

		if test.hook.GetWebhookID() != test.want.GetWebhookID() {
			t.Errorf("GetWebhookID is %v, want %v", test.hook.GetWebhookID(), test.want.GetWebhookID())
		}
	}
}

func TestTypes_Hook_Setters(t *testing.T) {
	// setup types
	var h *Hook

	// setup tests
	tests := []struct {
		hook *Hook
		want *Hook
	}{
		{
			hook: testHook(),
			want: testHook(),
		},
		{
			hook: h,
			want: new(Hook),
		},
	}

	// run tests
	for _, test := range tests {
		test.hook.SetID(test.want.GetID())
		test.hook.SetRepo(test.want.GetRepo())
		test.hook.SetBuild(test.want.GetBuild())
		test.hook.SetNumber(test.want.GetNumber())
		test.hook.SetSourceID(test.want.GetSourceID())
		test.hook.SetCreated(test.want.GetCreated())
		test.hook.SetHost(test.want.GetHost())
		test.hook.SetEvent(test.want.GetEvent())
		test.hook.SetEventAction(test.want.GetEventAction())
		test.hook.SetBranch(test.want.GetBranch())
		test.hook.SetError(test.want.GetError())
		test.hook.SetStatus(test.want.GetStatus())
		test.hook.SetLink(test.want.GetLink())
		test.hook.SetWebhookID(test.want.GetWebhookID())

		if test.hook.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.hook.GetID(), test.want.GetID())
		}

		if diff := cmp.Diff(test.hook.GetRepo(), test.want.GetRepo()); diff != "" {
			t.Errorf("SetRepo is a mismatch (-got +want):\n%s", diff)
		}

		if diff := cmp.Diff(test.hook.GetBuild(), test.want.GetBuild()); diff != "" {
			t.Errorf("SetBuild is a mismatch (-got +want):\n%s", diff)
		}

		if test.hook.GetNumber() != test.want.GetNumber() {
			t.Errorf("SetNumber is %v, want %v", test.hook.GetNumber(), test.want.GetNumber())
		}

		if test.hook.GetSourceID() != test.want.GetSourceID() {
			t.Errorf("SetSourceID is %v, want %v", test.hook.GetSourceID(), test.want.GetSourceID())
		}

		if test.hook.GetCreated() != test.want.GetCreated() {
			t.Errorf("SetCreated is %v, want %v", test.hook.GetCreated(), test.want.GetCreated())
		}

		if test.hook.GetHost() != test.want.GetHost() {
			t.Errorf("SetHost is %v, want %v", test.hook.GetHost(), test.want.GetHost())
		}

		if test.hook.GetEvent() != test.want.GetEvent() {
			t.Errorf("SetEvent is %v, want %v", test.hook.GetEvent(), test.want.GetEvent())
		}

		if test.hook.GetEventAction() != test.want.GetEventAction() {
			t.Errorf("SetEventAction is %v, want %v", test.hook.GetEventAction(), test.want.GetEventAction())
		}

		if test.hook.GetBranch() != test.want.GetBranch() {
			t.Errorf("SetBranch is %v, want %v", test.hook.GetBranch(), test.want.GetBranch())
		}

		if test.hook.GetError() != test.want.GetError() {
			t.Errorf("SetError is %v, want %v", test.hook.GetError(), test.want.GetError())
		}

		if test.hook.GetStatus() != test.want.GetStatus() {
			t.Errorf("SetStatus is %v, want %v", test.hook.GetStatus(), test.want.GetStatus())
		}

		if test.hook.GetLink() != test.want.GetLink() {
			t.Errorf("SetLink is %v, want %v", test.hook.GetLink(), test.want.GetLink())
		}

		if test.hook.GetWebhookID() != test.want.GetWebhookID() {
			t.Errorf("SetWebhookID is %v, want %v", test.hook.GetWebhookID(), test.want.GetWebhookID())
		}
	}
}

func TestTypes_Hook_String(t *testing.T) {
	// setup types
	h := testHook()

	want := fmt.Sprintf(`{
  Branch: %s,
  Build: %v,
  Created: %d,
  Error: %s,
  Event: %s,
  EventAction: %s,
  Host: %s,
  ID: %d,
  Link: %s,
  Number: %d,
  Repo: %v,
  SourceID: %s,
  Status: %s,
  WebhookID: %d,
}`,
		h.GetBranch(),
		h.GetBuild(),
		h.GetCreated(),
		h.GetError(),
		h.GetEvent(),
		h.GetEventAction(),
		h.GetHost(),
		h.GetID(),
		h.GetLink(),
		h.GetNumber(),
		h.GetRepo(),
		h.GetSourceID(),
		h.GetStatus(),
		h.GetWebhookID(),
	)

	// run test
	got := h.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testHook is a test helper function to create a Hook
// type with all fields set to a fake value.
func testHook() *Hook {
	h := new(Hook)

	h.SetID(1)
	h.SetRepo(testRepo())
	h.SetBuild(testBuild())
	h.SetNumber(1)
	h.SetSourceID("c8da1302-07d6-11ea-882f-4893bca275b8")
	h.SetCreated(time.Now().UTC().Unix())
	h.SetHost("github.com")
	h.SetEvent("push")
	h.SetEventAction("")
	h.SetBranch("main")
	h.SetError("")
	h.SetStatus("success")
	h.SetLink("https://github.com/github/octocat/settings/hooks/1")
	h.SetWebhookID(123456)

	return h
}
