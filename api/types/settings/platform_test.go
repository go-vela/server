// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTypes_Platform_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		platform *Platform
		want     *Platform
	}{
		{
			platform: testPlatformSettings(),
			want:     testPlatformSettings(),
		},
		{
			platform: new(Platform),
			want:     new(Platform),
		},
	}

	// run tests
	for _, test := range tests {
		if !reflect.DeepEqual(test.platform.GetCompiler(), test.want.GetCompiler()) {
			t.Errorf("GetCompiler is %v, want %v", test.platform.GetCompiler(), test.want.GetCompiler())
		}

		if !reflect.DeepEqual(test.platform.GetQueue(), test.want.GetQueue()) {
			t.Errorf("GetQueue is %v, want %v", test.platform.GetQueue(), test.want.GetQueue())
		}

		if !reflect.DeepEqual(test.platform.GetSCM(), test.want.GetSCM()) {
			t.Errorf("GetSCM is %v, want %v", test.platform.GetSCM(), test.want.GetSCM())
		}

		if !reflect.DeepEqual(test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist()) {
			t.Errorf("GetRepoAllowlist is %v, want %v", test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist())
		}

		if !reflect.DeepEqual(test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist()) {
			t.Errorf("GetScheduleAllowlist is %v, want %v", test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist())
		}

		if test.platform.GetMaxDashboardRepos() != test.want.GetMaxDashboardRepos() {
			t.Errorf("GetMaxDashboardRepos is %v, want %v", test.platform.GetMaxDashboardRepos(), test.want.GetMaxDashboardRepos())
		}

		if test.platform.GetQueueRestartLimit() != test.want.GetQueueRestartLimit() {
			t.Errorf("GetQueueRestartLimit is %v, want %v", test.platform.GetQueueRestartLimit(), test.want.GetQueueRestartLimit())
		}

		if test.platform.GetEnableRepoSecrets() != test.want.GetEnableRepoSecrets() {
			t.Errorf("GetEnableRepoSecrets is %v, want %v", test.platform.GetEnableRepoSecrets(), test.want.GetEnableRepoSecrets())
		}

		if test.platform.GetEnableOrgSecrets() != test.want.GetEnableOrgSecrets() {
			t.Errorf("GetEnableOrgSecrets is %v, want %v", test.platform.GetEnableOrgSecrets(), test.want.GetEnableOrgSecrets())
		}

		if test.platform.GetEnableSharedSecrets() != test.want.GetEnableSharedSecrets() {
			t.Errorf("GetEnableSharedSecrets is %v, want %v", test.platform.GetEnableSharedSecrets(), test.want.GetEnableSharedSecrets())
		}
	}
}

func TestTypes_Platform_Setters(t *testing.T) {
	// setup types
	var ps *Platform

	// setup tests
	tests := []struct {
		platform *Platform
		want     *Platform
	}{
		{
			platform: testPlatformSettings(),
			want:     testPlatformSettings(),
		},
		{
			platform: ps,
			want:     new(Platform),
		},
	}

	// run tests
	for _, test := range tests {
		test.platform.SetCompiler(test.want.GetCompiler())

		if !reflect.DeepEqual(test.platform.GetCompiler(), test.want.GetCompiler()) {
			t.Errorf("SetCompiler is %v, want %v", test.platform.GetCompiler(), test.want.GetCompiler())
		}

		test.platform.SetQueue(test.want.GetQueue())

		if !reflect.DeepEqual(test.platform.GetQueue(), test.want.GetQueue()) {
			t.Errorf("SetQueue is %v, want %v", test.platform.GetQueue(), test.want.GetQueue())
		}

		test.platform.SetSCM(test.want.GetSCM())

		if !reflect.DeepEqual(test.platform.GetSCM(), test.want.GetSCM()) {
			t.Errorf("SetSCM is %v, want %v", test.platform.GetSCM(), test.want.GetSCM())
		}

		test.platform.SetRepoAllowlist(test.want.GetRepoAllowlist())

		if !reflect.DeepEqual(test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist()) {
			t.Errorf("SetRepoAllowlist is %v, want %v", test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist())
		}

		test.platform.SetScheduleAllowlist(test.want.GetScheduleAllowlist())

		if !reflect.DeepEqual(test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist()) {
			t.Errorf("SetScheduleAllowlist is %v, want %v", test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist())
		}

		test.platform.SetMaxDashboardRepos(test.want.GetMaxDashboardRepos())

		if test.platform.GetMaxDashboardRepos() != test.want.GetMaxDashboardRepos() {
			t.Errorf("SetMaxDashboardRepos is %v, want %v", test.platform.GetMaxDashboardRepos(), test.want.GetMaxDashboardRepos())
		}

		test.platform.SetQueueRestartLimit(test.want.GetQueueRestartLimit())

		if test.platform.GetQueueRestartLimit() != test.want.GetQueueRestartLimit() {
			t.Errorf("SetQueueRestartLimit is %v, want %v", test.platform.GetQueueRestartLimit(), test.want.GetQueueRestartLimit())
		}

		test.platform.SetEnableRepoSecrets(test.want.GetEnableRepoSecrets())

		if test.platform.GetEnableRepoSecrets() != test.want.GetEnableRepoSecrets() {
			t.Errorf("SetEnableRepoSecrets is %v, want %v", test.platform.GetEnableRepoSecrets(), test.want.GetEnableRepoSecrets())
		}

		test.platform.SetEnableOrgSecrets(test.want.GetEnableOrgSecrets())

		if test.platform.GetEnableOrgSecrets() != test.want.GetEnableOrgSecrets() {
			t.Errorf("SetEnableOrgSecrets is %v, want %v", test.platform.GetEnableOrgSecrets(), test.want.GetEnableOrgSecrets())
		}

		test.platform.SetEnableSharedSecrets(test.want.GetEnableSharedSecrets())

		if test.platform.GetEnableSharedSecrets() != test.want.GetEnableSharedSecrets() {
			t.Errorf("SetEnableSharedSecrets is %v, want %v", test.platform.GetEnableSharedSecrets(), test.want.GetEnableSharedSecrets())
		}
	}
}

func TestTypes_Platform_Update(t *testing.T) {
	// setup types
	s := testPlatformSettings()

	// update fields
	sUpdate := testPlatformSettings()
	sUpdate.SetCompiler(Compiler{})
	sUpdate.SetQueue(Queue{})
	sUpdate.SetSCM(SCM{})
	sUpdate.SetRepoAllowlist([]string{"foo"})
	sUpdate.SetScheduleAllowlist([]string{"bar"})
	sUpdate.SetMaxDashboardRepos(20)
	sUpdate.SetQueueRestartLimit(60)
	sUpdate.SetEnableRepoSecrets(true)
	sUpdate.SetEnableOrgSecrets(true)
	sUpdate.SetEnableSharedSecrets(true)

	// setup tests
	tests := []struct {
		platform *Platform
		want     *Platform
	}{
		{
			platform: s,
			want:     testPlatformSettings(),
		},
		{
			platform: s,
			want:     sUpdate,
		},
	}

	// run tests
	for _, test := range tests {
		test.platform.FromSettings(test.want)

		if diff := cmp.Diff(test.want, test.platform); diff != "" {
			t.Errorf("(Update: -want +got):\n%s", diff)
		}
	}
}

func TestTypes_Platform_String(t *testing.T) {
	// setup types
	s := testPlatformSettings()
	cs := s.GetCompiler()
	qs := s.GetQueue()
	scms := s.GetSCM()

	want := fmt.Sprintf(`{
  ID: %d,
  Compiler: %v,
  Queue: %v,
  SCM: %v,
  RepoAllowlist: %v,
  ScheduleAllowlist: %v,
  MaxDashboardRepos: %d,
  QueueRestartLimit: %d,
  EnableRepoSecrets: %t,
  EnableOrgSecrets: %t,
  EnableSharedSecrets: %t,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		s.GetID(),
		cs.String(),
		qs.String(),
		scms.String(),
		s.GetRepoAllowlist(),
		s.GetScheduleAllowlist(),
		s.GetMaxDashboardRepos(),
		s.GetQueueRestartLimit(),
		s.GetEnableRepoSecrets(),
		s.GetEnableOrgSecrets(),
		s.GetEnableSharedSecrets(),
		s.GetCreatedAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
	)

	// run test
	got := s.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testPlatformSettings is a test helper function to create a Platform
// type with all fields set to a fake value.
func testPlatformSettings() *Platform {
	// setup platform
	s := new(Platform)
	s.SetID(1)
	s.SetCreatedAt(1)
	s.SetUpdatedAt(1)
	s.SetUpdatedBy("vela-server")
	s.SetRepoAllowlist([]string{"foo", "bar"})
	s.SetScheduleAllowlist([]string{"*"})
	s.SetMaxDashboardRepos(10)
	s.SetQueueRestartLimit(30)
	s.SetEnableRepoSecrets(false)
	s.SetEnableOrgSecrets(false)
	s.SetEnableSharedSecrets(false)

	// setup types
	// setup compiler
	cs := new(Compiler)

	cs.SetCloneImage("target/vela-git-slim:latest")
	cs.SetTemplateDepth(1)
	cs.SetStarlarkExecLimit(100)

	// setup queue
	qs := new(Queue)

	qs.SetRoutes([]string{"vela"})

	// setup scm
	scms := new(SCM)

	scms.SetRepoRoleMap(map[string]string{"foo": "bar"})
	scms.SetOrgRoleMap(map[string]string{"foo": "bar"})
	scms.SetTeamRoleMap(map[string]string{"foo": "bar"})

	s.SetCompiler(*cs)
	s.SetQueue(*qs)
	s.SetSCM(*scms)

	return s
}
