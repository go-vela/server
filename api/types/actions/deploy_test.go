// SPDX-License-Identifier: Apache-2.0

package actions

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestTypes_Deploy_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		actions *Deploy
		want    *Deploy
	}{
		{
			actions: testDeploy(),
			want:    testDeploy(),
		},
		{
			actions: new(Deploy),
			want:    new(Deploy),
		},
	}

	// run tests
	for _, test := range tests {
		if test.actions.GetCreated() != test.want.GetCreated() {
			t.Errorf("GetCreated is %v, want %v", test.actions.GetCreated(), test.want.GetCreated())
		}
	}
}

func TestTypes_Deploy_Setters(t *testing.T) {
	// setup types
	var a *Deploy

	// setup tests
	tests := []struct {
		actions *Deploy
		want    *Deploy
	}{
		{
			actions: testDeploy(),
			want:    testDeploy(),
		},
		{
			actions: a,
			want:    new(Deploy),
		},
	}

	// run tests
	for _, test := range tests {
		test.actions.SetCreated(test.want.GetCreated())

		if test.actions.GetCreated() != test.want.GetCreated() {
			t.Errorf("SetCreated is %v, want %v", test.actions.GetCreated(), test.want.GetCreated())
		}
	}
}

func TestTypes_Deploy_FromMask(t *testing.T) {
	// setup types
	mask := testMask()

	want := testDeploy()

	// run test
	got := new(Deploy).FromMask(mask)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("FromMask is %v, want %v", got, want)
	}
}

func TestTypes_Deploy_ToMask(t *testing.T) {
	// setup types
	actions := testDeploy()

	want := int64(constants.AllowDeployCreate)

	// run test
	got := actions.ToMask()

	if want != got {
		t.Errorf("ToMask is %v, want %v", got, want)
	}
}

func testDeploy() *Deploy {
	deploy := new(Deploy)
	deploy.SetCreated(true)

	return deploy
}
