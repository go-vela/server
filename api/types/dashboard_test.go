// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Dashboard_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		dashboard *Dashboard
		want      *Dashboard
	}{
		{
			dashboard: testDashboard(),
			want:      testDashboard(),
		},
		{
			dashboard: new(Dashboard),
			want:      new(Dashboard),
		},
	}

	// run tests
	for _, test := range tests {
		if test.dashboard.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.dashboard.GetID(), test.want.GetID())
		}

		if test.dashboard.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.dashboard.GetName(), test.want.GetName())
		}

		if !reflect.DeepEqual(test.dashboard.GetAdmins(), test.want.GetAdmins()) {
			t.Errorf("GetAdmins is %v, want %v", test.dashboard.GetAdmins(), test.want.GetAdmins())
		}

		if test.dashboard.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("GetCreatedAt is %v, want %v", test.dashboard.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.dashboard.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("GetCreatedBy is %v, want %v", test.dashboard.GetCreatedBy(), test.want.GetCreatedBy())
		}

		if test.dashboard.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("GetUpdatedAt is %v, want %v", test.dashboard.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.dashboard.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("GetUpdatedBy is %v, want %v", test.dashboard.GetUpdatedBy(), test.want.GetUpdatedBy())
		}

		if !reflect.DeepEqual(test.dashboard.GetRepos(), test.want.GetRepos()) {
			t.Errorf("GetRepos is %v, want %v", test.dashboard.GetRepos(), test.want.GetRepos())
		}
	}
}

func TestTypes_Dashboard_Setters(t *testing.T) {
	// setup types
	var d *Dashboard

	// setup tests
	tests := []struct {
		dashboard *Dashboard
		want      *Dashboard
	}{
		{
			dashboard: testDashboard(),
			want:      testDashboard(),
		},
		{
			dashboard: d,
			want:      new(Dashboard),
		},
	}

	// run tests
	for _, test := range tests {
		test.dashboard.SetID(test.want.GetID())
		test.dashboard.SetName(test.want.GetName())
		test.dashboard.SetAdmins(test.want.GetAdmins())
		test.dashboard.SetCreatedAt(test.want.GetCreatedAt())
		test.dashboard.SetCreatedBy(test.want.GetCreatedBy())
		test.dashboard.SetUpdatedAt(test.want.GetUpdatedAt())
		test.dashboard.SetUpdatedBy(test.want.GetUpdatedBy())
		test.dashboard.SetRepos(test.want.GetRepos())

		if test.dashboard.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.dashboard.GetID(), test.want.GetID())
		}

		if test.dashboard.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.dashboard.GetName(), test.want.GetName())
		}

		if !reflect.DeepEqual(test.dashboard.GetAdmins(), test.want.GetAdmins()) {
			t.Errorf("SetAdmins is %v, want %v", test.dashboard.GetAdmins(), test.want.GetAdmins())
		}

		if test.dashboard.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("SetCreatedAt is %v, want %v", test.dashboard.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.dashboard.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("SetCreatedBy is %v, want %v", test.dashboard.GetCreatedBy(), test.want.GetCreatedBy())
		}

		if test.dashboard.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("SetUpdatedAt is %v, want %v", test.dashboard.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.dashboard.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("SetUpdatedBy is %v, want %v", test.dashboard.GetUpdatedBy(), test.want.GetUpdatedBy())
		}

		if !reflect.DeepEqual(test.dashboard.GetRepos(), test.want.GetRepos()) {
			t.Errorf("SetRepos is %v, want %v", test.dashboard.GetRepos(), test.want.GetRepos())
		}
	}
}

func TestTypes_Dashboard_String(t *testing.T) {
	// setup types
	d := testDashboard()

	want := fmt.Sprintf(`{
  Name: %s,
  ID: %s,
  Admins: %v,
  CreatedAt: %d,
  CreatedBy: %s,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Repos: %v,
}`,
		d.GetName(),
		d.GetID(),
		d.GetAdmins(),
		d.GetCreatedAt(),
		d.GetCreatedBy(),
		d.GetUpdatedAt(),
		d.GetUpdatedBy(),
		d.GetRepos(),
	)

	// run test
	got := d.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testDashboard is a test helper function to create a Dashboard
// type with all fields set to a fake value.
func testDashboard() *Dashboard {
	d := new(Dashboard)

	d.SetID("123-abc")
	d.SetName("vela")
	d.SetAdmins([]string{"1", "42"})
	d.SetCreatedAt(1)
	d.SetCreatedBy("octocat")
	d.SetUpdatedAt(2)
	d.SetUpdatedBy("octokitty")
	d.SetRepos([]*DashboardRepo{testDashboardRepo()})

	return d
}
