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

		if !reflect.DeepEqual(test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist()) {
			t.Errorf("GetRepoAllowlist is %v, want %v", test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist())
		}

		if !reflect.DeepEqual(test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist()) {
			t.Errorf("GetScheduleAllowlist is %v, want %v", test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist())
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

		test.platform.SetRepoAllowlist(test.want.GetRepoAllowlist())

		if !reflect.DeepEqual(test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist()) {
			t.Errorf("SetRepoAllowlist is %v, want %v", test.platform.GetRepoAllowlist(), test.want.GetRepoAllowlist())
		}

		test.platform.SetScheduleAllowlist(test.want.GetScheduleAllowlist())

		if !reflect.DeepEqual(test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist()) {
			t.Errorf("SetScheduleAllowlist is %v, want %v", test.platform.GetScheduleAllowlist(), test.want.GetScheduleAllowlist())
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
	sUpdate.SetRepoAllowlist([]string{"foo"})
	sUpdate.SetScheduleAllowlist([]string{"bar"})

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

	want := fmt.Sprintf(`{
  ID: %d,
  Compiler: %v,
  Queue: %v,
  RepoAllowlist: %v,
  ScheduleAllowlist: %v,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		s.GetID(),
		cs.String(),
		qs.String(),
		s.GetRepoAllowlist(),
		s.GetScheduleAllowlist(),
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

	// setup types
	// setup compiler
	cs := new(Compiler)

	cs.SetCloneImage("target/vela-git:latest")
	cs.SetTemplateDepth(1)
	cs.SetStarlarkExecLimit(100)

	// setup queue
	qs := new(Queue)

	qs.SetRoutes([]string{"vela"})

	s.SetCompiler(*cs)
	s.SetQueue(*qs)

	return s
}
