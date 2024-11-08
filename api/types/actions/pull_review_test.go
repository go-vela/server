// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestActions_PullReview_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *PullReview
		want    *PullReview
	}{
		{
			actions: testPullReview(),
			want:    testPullReview(),
		},
		{
			actions: new(PullReview),
			want:    new(PullReview),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetSubmitted() != test.want.GetSubmitted() {
			t.Errorf("GetSubmitted is %v, want %v", test.actions.GetSubmitted(), test.want.GetSubmitted())
		}

		if test.actions.GetDismissed() != test.want.GetDismissed() {
			t.Errorf("GetDismissed is %v, want %v", test.actions.GetDismissed(), test.want.GetDismissed())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("GetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}
	}
}

func TestActions_PullReview_Setters(t *testing.T) {
	// setup types
	var a *PullReview

	// setup tests
	tests := []struct {
		actions *PullReview
		want    *PullReview
	}{
		{
			actions: testPullReview(),
			want:    testPullReview(),
		},
		{
			actions: a,
			want:    new(PullReview),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetSubmitted(test.want.GetSubmitted())
		test.actions.SetDismissed(test.want.GetDismissed())
		test.actions.SetEdited(test.want.GetEdited())

		if test.actions.GetSubmitted() != test.want.GetSubmitted() {
			t.Errorf("SetSubmitted is %v, want %v", test.actions.GetSubmitted(), test.want.GetSubmitted())
		}

		if test.actions.GetDismissed() != test.want.GetDismissed() {
			t.Errorf("SetDismissed is %v, want %v", test.actions.GetDismissed(), test.want.GetDismissed())
		}

		if test.actions.GetEdited() != test.want.GetEdited() {
			t.Errorf("SetEdited is %v, want %v", test.actions.GetEdited(), test.want.GetEdited())
		}
	}
}

func TestActions_PullReview_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testPullReview()

	// run test
	got := new(PullReview).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestActions_PullReview_ToMask(t *testing.T) {
	// setup types
	actions := testPullReview()

	want := int64(constants.AllowPullReviewSubmit | constants.AllowPullReviewDismiss)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testPullReview() *PullReview {
	pr := new(PullReview)

	pr.SetSubmitted(true)
	pr.SetEdited(false)
	pr.SetDismissed(true)

	return pr
}
