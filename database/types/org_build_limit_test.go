// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_OrgBuildLimit_Nullify(t *testing.T) {
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
			want:          nil,
		},
		{
			orgBuildLimit: new(OrgBuildLimit),
			want:          new(OrgBuildLimit),
		},
	}

	// run tests
	for _, test := range tests {
		got := test.orgBuildLimit.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_OrgBuildLimit_ToAPI(t *testing.T) {
	// setup types
	want := new(api.OrgBuildLimit)
	want.SetID(1)
	want.SetOrg("github")
	want.SetBuildLimit(30)
	want.SetCreatedAt(1)
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("octocat")

	// run test
	got := testOrgBuildLimit().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_OrgBuildLimit_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure       bool
		orgBuildLimit *OrgBuildLimit
	}{
		{
			failure:       false,
			orgBuildLimit: testOrgBuildLimit(),
		},
		{ // no Org set
			failure: true,
			orgBuildLimit: &OrgBuildLimit{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				BuildLimit: sql.NullInt32{Int32: 30, Valid: true},
			},
		},
		{ // no BuildLimit set
			failure: true,
			orgBuildLimit: &OrgBuildLimit{
				ID:  sql.NullInt64{Int64: 1, Valid: true},
				Org: sql.NullString{String: "github", Valid: true},
			},
		},
		{ // negative BuildLimit set
			failure: true,
			orgBuildLimit: &OrgBuildLimit{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				Org:        sql.NullString{String: "github", Valid: true},
				BuildLimit: sql.NullInt32{Int32: -1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.orgBuildLimit.Validate()

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

func TestTypes_OrgBuildLimit_OrgBuildLimitFromAPI(t *testing.T) {
	// setup types
	o := new(api.OrgBuildLimit)
	o.SetID(1)
	o.SetOrg("github")
	o.SetBuildLimit(30)
	o.SetCreatedAt(1)
	o.SetUpdatedAt(1)
	o.SetUpdatedBy("octocat")

	want := testOrgBuildLimit()

	// run test
	got := OrgBuildLimitFromAPI(o)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrgBuildLimitFromAPI is %v, want %v", got, want)
	}
}

// testOrgBuildLimit is a test helper function to create an OrgBuildLimit
// type with all fields set to a fake value.
func testOrgBuildLimit() *OrgBuildLimit {
	return &OrgBuildLimit{
		ID:         sql.NullInt64{Int64: 1, Valid: true},
		Org:        sql.NullString{String: "github", Valid: true},
		BuildLimit: sql.NullInt32{Int32: 30, Valid: true},
		CreatedAt:  sql.NullInt64{Int64: 1, Valid: true},
		UpdatedAt:  sql.NullInt64{Int64: 1, Valid: true},
		UpdatedBy:  sql.NullString{String: "octocat", Valid: true},
	}
}
