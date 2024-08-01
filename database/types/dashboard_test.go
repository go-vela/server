// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_Dashboard_Nullify(t *testing.T) {
	// setup types
	var h *Dashboard

	want := &Dashboard{
		Name:      sql.NullString{String: "", Valid: false},
		CreatedAt: sql.NullInt64{Int64: 0, Valid: false},
		CreatedBy: sql.NullString{String: "", Valid: false},
		UpdatedAt: sql.NullInt64{Int64: 0, Valid: false},
		UpdatedBy: sql.NullString{String: "", Valid: false},
	}

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
			dashboard: h,
			want:      nil,
		},
		{
			dashboard: new(Dashboard),
			want:      want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.dashboard.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_Dashboard_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Dashboard)
	want.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	want.SetName("vela")
	want.SetCreatedAt(1)
	want.SetCreatedBy("octocat")
	want.SetUpdatedAt(2)
	want.SetUpdatedBy("octokitty")
	want.SetAdmins(testAdminsJSON())
	want.SetRepos(testDashReposJSON())

	uuid, _ := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")
	h := &Dashboard{
		ID:        uuid,
		Name:      sql.NullString{String: "vela", Valid: true},
		CreatedAt: sql.NullInt64{Int64: 1, Valid: true},
		CreatedBy: sql.NullString{String: "octocat", Valid: true},
		UpdatedAt: sql.NullInt64{Int64: 2, Valid: true},
		UpdatedBy: sql.NullString{String: "octokitty", Valid: true},
		Admins:    testAdminsJSON(),
		Repos:     testDashReposJSON(),
	}

	// run test
	got := h.ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_Dashboard_Validate(t *testing.T) {
	uuid, _ := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")

	dashRepo := new(api.DashboardRepo)
	dashRepo.SetName("dashboard-repo")

	dashRepos := []*api.DashboardRepo{}
	for i := 0; i < 11; i++ {
		dashRepos = append(dashRepos, dashRepo)
	}

	exceededReposDashboard := testDashboard()
	exceededReposDashboard.Repos = DashReposJSON(dashRepos)

	// setup tests
	tests := []struct {
		failure   bool
		dashboard *Dashboard
	}{
		{
			failure:   false,
			dashboard: testDashboard(),
		},
		{ // no name set for dashboard
			failure: true,
			dashboard: &Dashboard{
				ID: uuid,
			},
		},
		{ // hit repo limit
			failure:   true,
			dashboard: exceededReposDashboard,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.dashboard.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestTypes_DashboardFromAPI(t *testing.T) {
	uuid, err := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")
	if err != nil {
		t.Errorf("error parsing uuid: %v", err)
	}

	// setup types
	want := &Dashboard{
		ID:        uuid,
		Name:      sql.NullString{String: "vela", Valid: true},
		CreatedAt: sql.NullInt64{Int64: 1, Valid: true},
		CreatedBy: sql.NullString{String: "octocat", Valid: true},
		UpdatedAt: sql.NullInt64{Int64: 2, Valid: true},
		UpdatedBy: sql.NullString{String: "octokitty", Valid: true},
		Admins:    testAdminsJSON(),
		Repos:     testDashReposJSON(),
	}

	d := new(api.Dashboard)
	d.SetID("c8da1302-07d6-11ea-882f-4893bca275b8")
	d.SetName("vela")
	d.SetCreatedAt(1)
	d.SetCreatedBy("octocat")
	d.SetUpdatedAt(2)
	d.SetUpdatedBy("octokitty")
	d.SetAdmins(testAdminsJSON())
	d.SetRepos(testDashReposJSON())

	// run test
	got := DashboardFromAPI(d)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("DashboardFromAPI() mismatch (-want +got):\n%s", diff)
	}

	// test empty uuid results in generated uuid
	d.SetID("")

	//nolint:staticcheck // linter is lying
	got = DashboardFromAPI(d)

	if len(got.ID) != 16 {
		t.Errorf("Length is %d", len(got.ID))
	}

	// test poorly formed uuid results in nil dashboard
	d.SetID("123-abc")

	got = DashboardFromAPI(d)

	if got != nil {
		t.Errorf("DashboardFromAPI should have returned nil")
	}
}

// testDashboard is a test helper function to create a Dashboard
// type with all fields set to a fake value.
func testDashboard() *Dashboard {
	uuid, _ := uuid.Parse("c8da1302-07d6-11ea-882f-4893bca275b8")

	return &Dashboard{
		ID:        uuid,
		Name:      sql.NullString{String: "vela", Valid: true},
		CreatedAt: sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
		CreatedBy: sql.NullString{String: "octocat", Valid: true},
		UpdatedAt: sql.NullInt64{Int64: time.Now().UTC().Unix(), Valid: true},
		UpdatedBy: sql.NullString{String: "octokitty", Valid: true},
		Admins:    testAdminsJSON(),
		Repos:     testDashReposJSON(),
	}
}

// testDashReposJSON is a test helper function to create a DashReposJSON
// type with all fields set to a fake value.
func testDashReposJSON() DashReposJSON {
	d := new(api.DashboardRepo)

	d.SetName("go-vela/server")
	d.SetID(1)
	d.SetBranches([]string{"main"})
	d.SetEvents([]string{"push", "tag"})

	return DashReposJSON{d}
}

// testAdminsJSON is a test helper function to create a DashReposJSON
// type with all fields set to a fake value.
func testAdminsJSON() AdminsJSON {
	u1 := new(api.User)

	u1.SetName("octocat")
	u1.SetID(1)
	u1.SetActive(true)
	u1.SetToken("foo")

	u2 := new(api.User)

	u2.SetName("octokitty")
	u2.SetID(2)
	u2.SetActive(true)
	u2.SetToken("bar")

	return AdminsJSON{u1, u2}
}
