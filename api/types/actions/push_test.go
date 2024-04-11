// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestTypes_Push_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *Push
		want    *Push
	}{
		{
			actions: testPush(),
			want:    testPush(),
		},
		{
			actions: new(Push),
			want:    new(Push),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetBranch() != test.want.GetBranch() {
			t.Errorf("GetBranch is %v, want %v", test.actions.GetBranch(), test.want.GetBranch())
		}

		if test.actions.GetTag() != test.want.GetTag() {
			t.Errorf("GetTag is %v, want %v", test.actions.GetTag(), test.want.GetTag())
		}
	}
}

func TestTypes_Push_Setters(t *testing.T) {
	// setup types
	var a *Push

	// setup tests
	tests := []struct {
		actions *Push
		want    *Push
	}{
		{
			actions: testPush(),
			want:    testPush(),
		},
		{
			actions: a,
			want:    new(Push),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetBranch(test.want.GetBranch())
		test.actions.SetTag(test.want.GetTag())
		test.actions.SetDeleteBranch(test.want.GetDeleteBranch())
		test.actions.SetDeleteTag(test.want.GetDeleteTag())

		if test.actions.GetBranch() != test.want.GetBranch() {
			t.Errorf("SetBranch is %v, want %v", test.actions.GetBranch(), test.want.GetBranch())
		}

		if test.actions.GetTag() != test.want.GetTag() {
			t.Errorf("SetTag is %v, want %v", test.actions.GetTag(), test.want.GetTag())
		}

		if test.actions.GetDeleteBranch() != test.want.GetDeleteBranch() {
			t.Errorf("SetDeleteBranch is %v, want %v", test.actions.GetDeleteBranch(), test.want.GetDeleteBranch())
		}

		if test.actions.GetDeleteTag() != test.want.GetDeleteTag() {
			t.Errorf("SetDeleteTag is %v, want %v", test.actions.GetDeleteTag(), test.want.GetDeleteTag())
		}
	}
}

func TestTypes_Push_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testPush()

	// run test
	got := new(Push).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestTypes_Push_ToMask(t *testing.T) {
	// setup types
	actions := testPush()

	want := int64(constants.AllowPushBranch | constants.AllowPushTag | constants.AllowPushDeleteBranch | constants.AllowPushDeleteTag)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testPush() *Push {
	push := new(Push)
	push.SetBranch(true)
	push.SetTag(true)
	push.SetDeleteBranch(true)
	push.SetDeleteTag(true)

	return push
}

func testMask() int64 {
	return int64(
		constants.AllowPushBranch |
			constants.AllowPushTag |
			constants.AllowPushDeleteBranch |
			constants.AllowPushDeleteTag |
			constants.AllowPullOpen |
			constants.AllowPullSync |
			constants.AllowPullReopen |
			constants.AllowPullUnlabel |
			constants.AllowDeployCreate |
			constants.AllowCommentCreate |
			constants.AllowSchedule,
	)
}
