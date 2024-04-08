// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
)

func TestActions_Pull_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *Pull
		want    *Pull
	}{
		{
			actions: testPull(),
			want:    testPull(),
		},
		{
			actions: new(Pull),
			want:    new(Pull),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetOpened() != test.want.GetOpened() {
			t.Errorf("GetOpened is %v, want %v", test.actions.GetOpened(), test.want.GetOpened())
		}

		if test.actions.GetSynchronize() != test.want.GetSynchronize() {
			t.Errorf("GetSynchronize is %v, want %v", test.actions.GetSynchronize(), test.want.GetSynchronize())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("GetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}

		if test.actions.GetReopened() != test.want.GetReopened() {
			t.Errorf("GetReopened is %v, want %v", test.actions.GetReopened(), test.want.GetReopened())
		}

		if test.actions.GetLabeled() != test.want.GetLabeled() {
			t.Errorf("GetLabeled is %v, want %v", test.actions.GetLabeled(), test.want.GetLabeled())
		}

		if test.actions.GetUnlabeled() != test.want.GetUnlabeled() {
			t.Errorf("GetUnlabeled is %v, want %v", test.actions.GetUnlabeled(), test.want.GetUnlabeled())
		}
	}
}

func TestActions_Pull_Setters(t *testing.T) {
	// setup types
	var a *Pull

	// setup tests
	tests := []struct {
		actions *Pull
		want    *Pull
	}{
		{
			actions: testPull(),
			want:    testPull(),
		},
		{
			actions: a,
			want:    new(Pull),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetOpened(test.want.GetOpened())
		test.actions.SetSynchronize(test.want.GetSynchronize())
		test.actions.SetEdited(test.want.GetEdited())
		test.actions.SetReopened(test.want.GetReopened())
		test.actions.SetLabeled(test.want.GetLabeled())
		test.actions.SetUnlabeled(test.want.GetUnlabeled())

		if test.actions.GetOpened() != test.want.GetOpened() {
			t.Errorf("SetOpened is %v, want %v", test.actions.GetOpened(), test.want.GetOpened())
		}

		if test.actions.GetSynchronize() != test.want.GetSynchronize() {
			t.Errorf("SetSynchronize is %v, want %v", test.actions.GetSynchronize(), test.want.GetSynchronize())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("SetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}

		if test.actions.GetReopened() != test.want.GetReopened() {
			t.Errorf("SetReopened is %v, want %v", test.actions.GetReopened(), test.want.GetReopened())
		}

		if test.actions.GetLabeled() != test.want.GetLabeled() {
			t.Errorf("SetLabeled is %v, want %v", test.actions.GetLabeled(), test.want.GetLabeled())
		}

		if test.actions.GetUnlabeled() != test.want.GetUnlabeled() {
			t.Errorf("SetUnlabeled is %v, want %v", test.actions.GetUnlabeled(), test.want.GetUnlabeled())
		}
	}
}

func TestActions_Pull_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testPull()

	// run test
	got := new(Pull).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestActions_Pull_ToMask(t *testing.T) {
	// setup types
	actions := testPull()

	want := int64(constants.AllowPullOpen | constants.AllowPullSync | constants.AllowPullReopen | constants.AllowPullUnlabel)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testPull() *Pull {
	pr := new(Pull)
	pr.SetOpened(true)
	pr.SetSynchronize(true)
	pr.SetEdited(false)
	pr.SetReopened(true)
	pr.SetLabeled(false)
	pr.SetUnlabeled(true)

	return pr
}
