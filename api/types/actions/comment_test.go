// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestTypes_Comment_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *Comment
		want    *Comment
	}{
		{
			actions: testComment(),
			want:    testComment(),
		},
		{
			actions: new(Comment),
			want:    new(Comment),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetCreated() != test.want.GetCreated() {
			t.Errorf("GetCreated is %v, want %v", test.actions.GetCreated(), test.want.GetCreated())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("GetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}
	}
}

func TestTypes_Comment_Setters(t *testing.T) {
	// setup types
	var a *Comment

	// setup tests
	tests := []struct {
		actions *Comment
		want    *Comment
	}{
		{
			actions: testComment(),
			want:    testComment(),
		},
		{
			actions: a,
			want:    new(Comment),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetCreated(test.want.GetCreated())
		test.actions.SetEdited(test.want.GetEdited())

		if test.actions.GetCreated() != test.want.GetCreated() {
			t.Errorf("SetCreated is %v, want %v", test.actions.GetCreated(), test.want.GetCreated())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("SetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}
	}
}

func TestTypes_Comment_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testComment()

	// run test
	got := new(Comment).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestTypes_Comment_ToMask(t *testing.T) {
	// setup types
	actions := testComment()

	want := int64(constants.AllowCommentCreate)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testComment() *Comment {
	comment := new(Comment)
	comment.SetCreated(true)
	comment.SetEdited(false)

	return comment
}
