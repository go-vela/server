// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_SCM_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		scm  *SCM
		want *SCM
	}{
		{
			scm:  testSCMSettings(),
			want: testSCMSettings(),
		},
		{
			scm:  new(SCM),
			want: new(SCM),
		},
	}

	// run tests
	for _, test := range tests {
		if !reflect.DeepEqual(test.scm.GetRepoRoleMap(), test.want.GetRepoRoleMap()) {
			t.Errorf("GetRepoRoleMap is %v, want %v", test.scm.GetRepoRoleMap(), test.want.GetRepoRoleMap())
		}

		if !reflect.DeepEqual(test.scm.GetOrgRoleMap(), test.want.GetOrgRoleMap()) {
			t.Errorf("GetOrgRoleMap is %v, want %v", test.scm.GetOrgRoleMap(), test.want.GetOrgRoleMap())
		}

		if !reflect.DeepEqual(test.scm.GetTeamRoleMap(), test.want.GetTeamRoleMap()) {
			t.Errorf("GetTeamRoleMap is %v, want %v", test.scm.GetTeamRoleMap(), test.want.GetTeamRoleMap())
		}
	}
}

func TestTypes_SCM_Setters(t *testing.T) {
	// setup types
	var qs *SCM

	// setup tests
	tests := []struct {
		scm  *SCM
		want *SCM
	}{
		{
			scm:  testSCMSettings(),
			want: testSCMSettings(),
		},
		{
			scm:  qs,
			want: new(SCM),
		},
	}

	// run tests
	for _, test := range tests {
		test.scm.SetRepoRoleMap(test.want.GetRepoRoleMap())
		test.scm.SetOrgRoleMap(test.want.GetOrgRoleMap())
		test.scm.SetTeamRoleMap(test.want.GetTeamRoleMap())

		if !reflect.DeepEqual(test.scm.GetRepoRoleMap(), test.want.GetRepoRoleMap()) {
			t.Errorf("SetRepoRoleMap is %v, want %v", test.scm.GetRepoRoleMap(), test.want.GetRepoRoleMap())
		}

		if !reflect.DeepEqual(test.scm.GetOrgRoleMap(), test.want.GetOrgRoleMap()) {
			t.Errorf("SetOrgRoleMap is %v, want %v", test.scm.GetOrgRoleMap(), test.want.GetOrgRoleMap())
		}

		if !reflect.DeepEqual(test.scm.GetTeamRoleMap(), test.want.GetTeamRoleMap()) {
			t.Errorf("SetTeamRoleMap is %v, want %v", test.scm.GetTeamRoleMap(), test.want.GetTeamRoleMap())
		}
	}
}

func TestTypes_SCM_String(t *testing.T) {
	// setup types
	scms := testSCMSettings()

	// setup tests
	tests := []struct {
		scm  *SCM
		want string
	}{
		{
			scm: scms,
			want: fmt.Sprintf(`{
  RepoRoleMap: %v,
  OrgRoleMap: %v,
  TeamRoleMap: %v,
}`,
				scms.GetRepoRoleMap(),
				scms.GetOrgRoleMap(),
				scms.GetTeamRoleMap(),
			),
		},
	}

	// run tests
	for _, test := range tests {
		if test.scm.String() != test.want {
			t.Errorf("String is %s, want %s", test.scm.String(), test.want)
		}
	}
}

// testSCMSettings is a test helper function to create a SCM
// type with all fields set to a fake value.
func testSCMSettings() *SCM {
	scms := new(SCM)

	// set fake values
	scms.SetRepoRoleMap(map[string]string{"admin": "admin", "triage": "read"})
	scms.SetOrgRoleMap(map[string]string{"admin": "admin", "member": "read"})
	scms.SetTeamRoleMap(map[string]string{"admin": "admin"})

	return scms
}
