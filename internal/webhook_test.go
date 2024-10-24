// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestWebhook_ShouldSkip(t *testing.T) {
	// set up tests
	tests := []struct {
		hook       *Webhook
		wantBool   bool
		wantString string
	}{
		{
			&Webhook{Build: testPushBuild("testing [SKIP CI]", "", constants.EventPush)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing", "wip [ci skip]", constants.EventPush)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing [skip VELA]", "", constants.EventPush)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing", "wip [vela skip]", constants.EventPush)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing ***NO_CI*** ok", "nothing", constants.EventPush)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing ok", "nothing", constants.EventPush)},
			false,
			"",
		},
		{
			&Webhook{Build: testPushBuild("testing [SKIP CI]", "", constants.EventTag)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing", "wip [ci skip]", constants.EventTag)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing [skip VELA]", "", constants.EventTag)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing", "wip [vela skip]", constants.EventTag)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing ***NO_CI*** ok", "nothing", constants.EventTag)},
			true,
			skipDirectiveMsg,
		},
		{
			&Webhook{Build: testPushBuild("testing ok", "nothing", constants.EventTag)},
			false,
			"",
		},
	}

	// run tests
	for _, test := range tests {
		gotBool, gotString := test.hook.ShouldSkip()

		if gotString != test.wantString {
			t.Errorf("returned an error, wanted %s, but got %s", test.wantString, gotString)
		}

		if gotBool != test.wantBool {
			t.Errorf("returned an error, wanted %v, but got %v", test.wantBool, gotBool)
		}
	}
}

func testPushBuild(message, title, event string) *api.Build {
	b := new(api.Build)

	b.SetEvent(event)

	if len(message) > 0 {
		b.SetMessage(message)
	}

	if len(title) > 0 {
		b.SetTitle(title)
	}

	b.SetCommit("deadbeef")

	return b
}
