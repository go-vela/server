// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLibrary_DashboardRepo_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		dashboardRepo *DashboardRepo
		want          *DashboardRepo
	}{
		{
			dashboardRepo: testDashboardRepo(),
			want:          testDashboardRepo(),
		},
		{
			dashboardRepo: new(DashboardRepo),
			want:          new(DashboardRepo),
		},
	}

	// run tests
	for _, test := range tests {
		if test.dashboardRepo.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.dashboardRepo.GetID(), test.want.GetID())
		}

		if test.dashboardRepo.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.dashboardRepo.GetName(), test.want.GetName())
		}

		if !reflect.DeepEqual(test.dashboardRepo.GetBranches(), test.want.GetBranches()) {
			t.Errorf("GetBranches is %v, want %v", test.dashboardRepo.GetBranches(), test.want.GetBranches())
		}

		if !reflect.DeepEqual(test.dashboardRepo.GetEvents(), test.want.GetEvents()) {
			t.Errorf("GetEvents is %v, want %v", test.dashboardRepo.GetEvents(), test.want.GetEvents())
		}
	}
}

func TestLibrary_DashboardRepo_Setters(t *testing.T) {
	// setup types
	var d *DashboardRepo

	// setup tests
	tests := []struct {
		dashboardRepo *DashboardRepo
		want          *DashboardRepo
	}{
		{
			dashboardRepo: testDashboardRepo(),
			want:          testDashboardRepo(),
		},
		{
			dashboardRepo: d,
			want:          new(DashboardRepo),
		},
	}

	// run tests
	for _, test := range tests {
		test.dashboardRepo.SetID(test.want.GetID())
		test.dashboardRepo.SetName(test.want.GetName())
		test.dashboardRepo.SetBranches(test.want.GetBranches())
		test.dashboardRepo.SetEvents(test.want.GetEvents())

		if test.dashboardRepo.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.dashboardRepo.GetID(), test.want.GetID())
		}

		if test.dashboardRepo.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.dashboardRepo.GetName(), test.want.GetName())
		}

		if !reflect.DeepEqual(test.dashboardRepo.GetBranches(), test.want.GetBranches()) {
			t.Errorf("SetBranches is %v, want %v", test.dashboardRepo.GetBranches(), test.want.GetBranches())
		}

		if !reflect.DeepEqual(test.dashboardRepo.GetEvents(), test.want.GetEvents()) {
			t.Errorf("SetEvents is %v, want %v", test.dashboardRepo.GetEvents(), test.want.GetEvents())
		}
	}
}

func TestLibrary_DashboardRepo_String(t *testing.T) {
	// setup types
	d := testDashboardRepo()

	want := fmt.Sprintf(`{
  Name: %s,
  ID: %d,
  Branches: %v,
  Events: %v,
}`,
		d.GetName(),
		d.GetID(),
		d.GetBranches(),
		d.GetEvents(),
	)

	// run test
	got := d.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testDashboardRepo is a test helper function to create a DashboardRepo
// type with all fields set to a fake value.
func testDashboardRepo() *DashboardRepo {
	d := new(DashboardRepo)

	d.SetName("go-vela/server")
	d.SetID(1)
	d.SetBranches([]string{"main"})
	d.SetEvents([]string{"push", "tag"})

	return d
}
