// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types/settings"
)

func TestTypes_Platform_Nullify(t *testing.T) {
	// setup types
	var ps *Platform

	want := &Platform{
		ID: sql.NullInt32{Int32: 0, Valid: false},
	}

	// setup tests
	tests := []struct {
		repo *Platform
		want *Platform
	}{
		{
			repo: testPlatform(),
			want: testPlatform(),
		},
		{
			repo: ps,
			want: nil,
		},
		{
			repo: new(Platform),
			want: want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.repo.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_Platform_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Platform)
	want.SetID(1)
	want.SetRepoAllowlist([]string{"github/octocat"})
	want.SetScheduleAllowlist([]string{"*"})
	want.SetMaxDashboardRepos(10)
	want.SetQueueRestartLimit(30)
	want.SetCreatedAt(0)
	want.SetUpdatedAt(0)
	want.SetUpdatedBy("")

	want.Compiler = new(api.Compiler)
	want.SetCloneImage("target/vela-git-slim:latest")
	want.SetTemplateDepth(10)
	want.SetStarlarkExecLimit(100)

	want.Queue = new(api.Queue)
	want.SetRoutes([]string{"vela"})

	want.SCM = new(api.SCM)
	want.SetRepoRoleMap(map[string]string{
		"admin":  "admin",
		"triage": "read",
	})
	want.SetOrgRoleMap(map[string]string{
		"admin":  "admin",
		"member": "read",
	})
	want.SetTeamRoleMap(map[string]string{
		"admin": "admin",
	})

	// run test
	got := testPlatform().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_Platform_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		settings *Platform
	}{
		{
			failure:  false,
			settings: testPlatform(),
		},
		{ // no CloneImage set for settings
			failure: true,
			settings: &Platform{
				ID:                sql.NullInt32{Int32: 1, Valid: true},
				MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
				Compiler: Compiler{
					TemplateDepth:     sql.NullInt64{Int64: 10, Valid: true},
					StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
				},
			},
		},
		{ // no TemplateDepth set for settings
			failure: true,
			settings: &Platform{
				ID:                sql.NullInt32{Int32: 1, Valid: true},
				MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
				Compiler: Compiler{
					CloneImage:        sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
					StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
				},
			},
		},
		{ // no StarlarkExecLimit set for settings
			failure: true,
			settings: &Platform{
				ID:                sql.NullInt32{Int32: 1, Valid: true},
				MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
				Compiler: Compiler{
					CloneImage:    sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
					TemplateDepth: sql.NullInt64{Int64: 10, Valid: true},
				},
			},
		},
		{ // no MaxDashboardRepos set for settings
			failure: true,
			settings: &Platform{
				ID: sql.NullInt32{Int32: 1, Valid: true},
				Compiler: Compiler{
					CloneImage:        sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
					TemplateDepth:     sql.NullInt64{Int64: 10, Valid: true},
					StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
				},
			},
		},
		{ // negative QueueRestartLimit set for settings
			failure: true,
			settings: &Platform{
				ID:                sql.NullInt32{Int32: 1, Valid: true},
				MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
				QueueRestartLimit: sql.NullInt32{Int32: -1, Valid: true},
				Compiler: Compiler{
					CloneImage:        sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
					TemplateDepth:     sql.NullInt64{Int64: 10, Valid: true},
					StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
				},
			},
		},
		{ // no queue fields set for settings
			failure: false,
			settings: &Platform{
				ID:                sql.NullInt32{Int32: 1, Valid: true},
				MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
				Compiler: Compiler{
					CloneImage:        sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
					TemplateDepth:     sql.NullInt64{Int64: 10, Valid: true},
					StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
				},
				Queue: Queue{},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.settings.Validate()

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

func TestTypes_Platform_PlatformFromAPI(t *testing.T) {
	// setup types
	s := new(api.Platform)
	s.SetID(1)
	s.SetRepoAllowlist([]string{"github/octocat"})
	s.SetScheduleAllowlist([]string{"*"})
	s.SetMaxDashboardRepos(10)
	s.SetQueueRestartLimit(30)
	s.SetCreatedAt(0)
	s.SetUpdatedAt(0)
	s.SetUpdatedBy("")

	s.Compiler = new(api.Compiler)
	s.SetCloneImage("target/vela-git-slim:latest")
	s.SetTemplateDepth(10)
	s.SetStarlarkExecLimit(100)

	s.Queue = new(api.Queue)
	s.SetRoutes([]string{"vela"})

	s.SCM = new(api.SCM)
	s.SetRepoRoleMap(map[string]string{
		"admin":  "admin",
		"triage": "read",
	})
	s.SetOrgRoleMap(map[string]string{
		"admin":  "admin",
		"member": "read",
	})
	s.SetTeamRoleMap(map[string]string{
		"admin": "admin",
	})

	want := testPlatform()

	// run test
	got := SettingsFromAPI(s)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("PlatformFromAPI is %v, want %v", got, want)
	}
}

// testPlatform is a test helper function to create a Platform
// type with all fields set to a fake value.
func testPlatform() *Platform {
	return &Platform{
		ID: sql.NullInt32{Int32: 1, Valid: true},
		Compiler: Compiler{
			CloneImage:        sql.NullString{String: "target/vela-git-slim:latest", Valid: true},
			TemplateDepth:     sql.NullInt64{Int64: 10, Valid: true},
			StarlarkExecLimit: sql.NullInt64{Int64: 100, Valid: true},
		},
		Queue: Queue{
			Routes: []string{"vela"},
		},
		SCM: SCM{
			RepoRoleMap: map[string]string{
				"admin":  "admin",
				"triage": "read",
			},
			OrgRoleMap: map[string]string{
				"admin":  "admin",
				"member": "read",
			},
			TeamRoleMap: map[string]string{
				"admin": "admin",
			},
		},
		RepoAllowlist:     []string{"github/octocat"},
		ScheduleAllowlist: []string{"*"},
		MaxDashboardRepos: sql.NullInt32{Int32: 10, Valid: true},
		QueueRestartLimit: sql.NullInt32{Int32: 30, Valid: true},
		CreatedAt:         sql.NullInt64{Int64: 0, Valid: true},
		UpdatedAt:         sql.NullInt64{Int64: 0, Valid: true},
		UpdatedBy:         sql.NullString{String: "", Valid: true},
	}
}
