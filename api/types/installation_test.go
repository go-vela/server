// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"
)

func TestTypes_Installation_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		installation *Installation
		want         *Installation
	}{
		{
			installation: testInstallation(),
			want:         testInstallation(),
		},
		{
			installation: new(Installation),
			want:         new(Installation),
		},
	}

	// run tests
	for _, test := range tests {
		if test.installation.GetInstallID() != test.want.GetInstallID() {
			t.Errorf("GetInstallID is %v, want %v", test.installation.GetInstallID(), test.want.GetInstallID())
		}

		if test.installation.GetTarget() != test.want.GetTarget() {
			t.Errorf("GetTarget is %v, want %v", test.installation.GetTarget(), test.want.GetTarget())
		}
	}
}

func TestTypes_Installation_Setters(t *testing.T) {
	// setup types
	var i *Installation

	// setup tests
	tests := []struct {
		installation *Installation
		want         *Installation
	}{
		{
			installation: testInstallation(),
			want:         testInstallation(),
		},
		{
			installation: i,
			want:         new(Installation),
		},
	}

	// run tests
	for _, test := range tests {
		test.installation.SetInstallID(test.want.GetInstallID())
		test.installation.SetTarget(test.want.GetTarget())

		if test.installation.GetInstallID() != test.want.GetInstallID() {
			t.Errorf("SetInstallID is %v, want %v", test.installation.GetInstallID(), test.want.GetInstallID())
		}

		if test.installation.GetTarget() != test.want.GetTarget() {
			t.Errorf("SetTarget is %v, want %v", test.installation.GetTarget(), test.want.GetTarget())
		}
	}
}

// testInstallation is a test helper function to create an Installation
// type with all fields set to a fake value.
func testInstallation() *Installation {
	i := new(Installation)

	i.SetInstallID(1)
	i.SetTarget("octocat")

	return i
}
