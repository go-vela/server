// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_OrgBuildLimit_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		orgBuildLimit *OrgBuildLimit
		want          *OrgBuildLimit
	}{
		{
			orgBuildLimit: testOrgBuildLimit(),
			want:          testOrgBuildLimit(),
		},
		{
			orgBuildLimit: new(OrgBuildLimit),
			want:          new(OrgBuildLimit),
		},
	}

	// run tests
	for _, test := range tests {
		if test.orgBuildLimit.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.orgBuildLimit.GetID(), test.want.GetID())
		}

		if test.orgBuildLimit.GetOrg() != test.want.GetOrg() {
			t.Errorf("GetOrg is %v, want %v", test.orgBuildLimit.GetOrg(), test.want.GetOrg())
		}

		if test.orgBuildLimit.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("GetBuildLimit is %v, want %v", test.orgBuildLimit.GetBuildLimit(), test.want.GetBuildLimit())
		}

		if test.orgBuildLimit.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("GetCreatedAt is %v, want %v", test.orgBuildLimit.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.orgBuildLimit.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("GetUpdatedAt is %v, want %v", test.orgBuildLimit.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.orgBuildLimit.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("GetUpdatedBy is %v, want %v", test.orgBuildLimit.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_OrgBuildLimit_Setters(t *testing.T) {
	// setup types
	var o *OrgBuildLimit

	// setup tests
	tests := []struct {
		orgBuildLimit *OrgBuildLimit
		want          *OrgBuildLimit
	}{
		{
			orgBuildLimit: testOrgBuildLimit(),
			want:          testOrgBuildLimit(),
		},
		{
			orgBuildLimit: o,
			want:          new(OrgBuildLimit),
		},
	}

	// run tests
	for _, test := range tests {
		test.orgBuildLimit.SetID(test.want.GetID())
		test.orgBuildLimit.SetOrg(test.want.GetOrg())
		test.orgBuildLimit.SetBuildLimit(test.want.GetBuildLimit())
		test.orgBuildLimit.SetCreatedAt(test.want.GetCreatedAt())
		test.orgBuildLimit.SetUpdatedAt(test.want.GetUpdatedAt())
		test.orgBuildLimit.SetUpdatedBy(test.want.GetUpdatedBy())

		if test.orgBuildLimit.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.orgBuildLimit.GetID(), test.want.GetID())
		}

		if test.orgBuildLimit.GetOrg() != test.want.GetOrg() {
			t.Errorf("SetOrg is %v, want %v", test.orgBuildLimit.GetOrg(), test.want.GetOrg())
		}

		if test.orgBuildLimit.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("SetBuildLimit is %v, want %v", test.orgBuildLimit.GetBuildLimit(), test.want.GetBuildLimit())
		}

		if test.orgBuildLimit.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("SetCreatedAt is %v, want %v", test.orgBuildLimit.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.orgBuildLimit.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("SetUpdatedAt is %v, want %v", test.orgBuildLimit.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.orgBuildLimit.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("SetUpdatedBy is %v, want %v", test.orgBuildLimit.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_OrgBuildLimit_String(t *testing.T) {
	// setup types
	o := testOrgBuildLimit()

	want := fmt.Sprintf(`{
  ID: %d,
  Org: %s,
  BuildLimit: %d,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		o.GetID(),
		o.GetOrg(),
		o.GetBuildLimit(),
		o.GetCreatedAt(),
		o.GetUpdatedAt(),
		o.GetUpdatedBy(),
	)

	// run test
	got := o.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testOrgBuildLimit is a test helper function to create an OrgBuildLimit
// type with all fields set to a fake value.
func testOrgBuildLimit() *OrgBuildLimit {
	o := new(OrgBuildLimit)

	o.SetID(1)
	o.SetOrg("github")
	o.SetBuildLimit(30)
	o.SetCreatedAt(1)
	o.SetUpdatedAt(1)
	o.SetUpdatedBy("octocat")

	return o
}
