// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/types/pipeline"
)

func TestTypes_Secret_Sanitize(t *testing.T) {
	// setup types
	s := testSecret()

	want := testSecret()
	want.SetValue(constants.SecretMask)

	// run test
	got := s.Sanitize()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sanitize is %v, want %v", got, want)
	}
}

func TestTypes_Secret_Match(t *testing.T) {
	// setup types
	v := "foo"
	fBool := false
	tBool := true

	testEvents := &Events{
		Push: &actions.Push{
			Branch: &tBool,
			Tag:    &tBool,
		},
		PullRequest: &actions.Pull{
			Opened:      &fBool,
			Synchronize: &tBool,
			Edited:      &fBool,
		},
		Deployment: &actions.Deploy{
			Created: &tBool,
		},
		Comment: &actions.Comment{
			Created: &tBool,
			Edited:  &tBool,
		},
		Schedule: &actions.Schedule{
			Run: &tBool,
		},
	}

	// setup tests
	tests := []struct {
		name string
		step *pipeline.Container
		sec  *Secret
		want bool
	}{
		{ // test matching secret events
			name: "push",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "pull request opened fail",
			step: &pipeline.Container{
				Image: "alpine:latest",
				Environment: map[string]string{
					"VELA_BUILD_EVENT":        "pull_request",
					"VELA_BUILD_EVENT_ACTION": "opened",
				},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: false,
		},
		{
			name: "tag",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "tag"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "deployment",
			step: &pipeline.Container{
				Image: "alpine:latest",
				Environment: map[string]string{
					"VELA_BUILD_EVENT":        "deployment",
					"VELA_BUILD_EVENT_ACTION": "created",
				},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "comment created",
			step: &pipeline.Container{
				Image: "alpine:latest",
				Environment: map[string]string{
					"VELA_BUILD_EVENT":        "comment",
					"VELA_BUILD_EVENT_ACTION": "created",
				},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "fake event",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "fake_event"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: false,
		},
		{
			name: "push with empty image constraint",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "schedule",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "schedule"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{},
				AllowEvents: testEvents,
			},
			want: true,
		},

		{ // test matching secret images
			name: "image basic",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine"},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "image and tag",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine:latest"},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "mismatch tag with same image",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine:1"},
				AllowEvents: testEvents,
			},
			want: false,
		},
		{
			name: "multiple allowed images",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine", "centos"},
				AllowEvents: testEvents,
			},
			want: true,
		},

		{ // test matching secret events and images
			name: "push and image pass",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine"},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "push and image tag pass",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine:latest"},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "push and bad image tag",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine:1"},
				AllowEvents: testEvents,
			},
			want: false,
		},
		{
			name: "mismatch event and match image",
			step: &pipeline.Container{
				Image: "alpine:latest",
				Environment: map[string]string{
					"VELA_BUILD_EVENT":        "pull_request",
					"VELA_BUILD_EVENT_ACTION": "edited",
				},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine:latest"},
				AllowEvents: testEvents,
			},
			want: false,
		},
		{
			name: "event and multi image allowed pass",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine", "centos"},
				AllowEvents: testEvents,
			},
			want: true,
		},

		{ // test build events with image ACLs and rulesets
			name: "rulesets and push pass",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Ruleset: pipeline.Ruleset{
					If: pipeline.Rules{
						Event: []string{"push"},
					},
				},
			},
			sec: &Secret{
				Name:        &v,
				Value:       &v,
				Images:      &[]string{"alpine"},
				AllowEvents: testEvents,
			},
			want: true,
		},
		{
			name: "no commands allowed",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Ruleset: pipeline.Ruleset{
					If: pipeline.Rules{
						Event: []string{"push"},
					},
				},
				Commands: []string{"echo hi"},
			},
			sec: &Secret{
				Name:         &v,
				Value:        &v,
				Images:       &[]string{"alpine"},
				AllowEvents:  testEvents,
				AllowCommand: &fBool,
			},
			want: false,
		},
		{
			name: "no commands allowed - entrypoint provided",
			step: &pipeline.Container{
				Image:       "alpine:latest",
				Environment: map[string]string{"VELA_BUILD_EVENT": "push"},
				Ruleset: pipeline.Ruleset{
					If: pipeline.Rules{
						Event: []string{"push"},
					},
				},
				Entrypoint: []string{"sh", "-c", "echo hi"},
			},
			sec: &Secret{
				Name:         &v,
				Value:        &v,
				Images:       &[]string{"alpine"},
				AllowEvents:  testEvents,
				AllowCommand: &fBool,
			},
			want: false,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.sec.Match(test.step)

		if got != test.want {
			t.Errorf("Match for %s is %v, want %v", test.name, got, test.want)
		}
	}
}

func TestTypes_Secret_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		secret *Secret
		want   *Secret
	}{
		{
			secret: testSecret(),
			want:   testSecret(),
		},
		{
			secret: new(Secret),
			want:   new(Secret),
		},
	}

	// run tests
	for _, test := range tests {
		if test.secret.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.secret.GetID(), test.want.GetID())
		}

		if test.secret.GetOrg() != test.want.GetOrg() {
			t.Errorf("GetOrg is %v, want %v", test.secret.GetOrg(), test.want.GetOrg())
		}

		if test.secret.GetRepo() != test.want.GetRepo() {
			t.Errorf("GetRepo is %v, want %v", test.secret.GetRepo(), test.want.GetRepo())
		}

		if test.secret.GetTeam() != test.want.GetTeam() {
			t.Errorf("GetTeam is %v, want %v", test.secret.GetTeam(), test.want.GetTeam())
		}

		if test.secret.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.secret.GetName(), test.want.GetName())
		}

		if test.secret.GetValue() != test.want.GetValue() {
			t.Errorf("GetValue is %v, want %v", test.secret.GetValue(), test.want.GetValue())
		}

		if test.secret.GetType() != test.want.GetType() {
			t.Errorf("GetType is %v, want %v", test.secret.GetType(), test.want.GetType())
		}

		if !reflect.DeepEqual(test.secret.GetImages(), test.want.GetImages()) {
			t.Errorf("GetImages is %v, want %v", test.secret.GetImages(), test.want.GetImages())
		}

		if !reflect.DeepEqual(test.secret.GetAllowEvents(), test.want.GetAllowEvents()) {
			t.Errorf("GetAllowEvents is %v, want %v", test.secret.GetAllowEvents(), test.want.GetAllowEvents())
		}

		if test.secret.GetAllowCommand() != test.want.GetAllowCommand() {
			t.Errorf("GetAllowCommand is %v, want %v", test.secret.GetAllowCommand(), test.want.GetAllowCommand())
		}

		if test.secret.GetAllowSubstitution() != test.want.GetAllowSubstitution() {
			t.Errorf("GetAllowSubstitution is %v, want %v", test.secret.GetAllowSubstitution(), test.want.GetAllowSubstitution())
		}

		if test.secret.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("GetCreatedAt is %v, want %v", test.secret.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.secret.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("GetCreatedBy is %v, want %v", test.secret.GetCreatedBy(), test.want.GetCreatedBy())
		}

		if test.secret.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("GetUpdatedAt is %v, want %v", test.secret.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.secret.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("GetUpdatedBy is %v, want %v", test.secret.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_Secret_Setters(t *testing.T) {
	// setup types
	var s *Secret

	// setup tests
	tests := []struct {
		secret *Secret
		want   *Secret
	}{
		{
			secret: testSecret(),
			want:   testSecret(),
		},
		{
			secret: s,
			want:   new(Secret),
		},
	}

	// run tests
	for _, test := range tests {
		test.secret.SetID(test.want.GetID())
		test.secret.SetOrg(test.want.GetOrg())
		test.secret.SetRepo(test.want.GetRepo())
		test.secret.SetTeam(test.want.GetTeam())
		test.secret.SetName(test.want.GetName())
		test.secret.SetValue(test.want.GetValue())
		test.secret.SetType(test.want.GetType())
		test.secret.SetImages(test.want.GetImages())
		test.secret.SetAllowEvents(test.want.GetAllowEvents())
		test.secret.SetAllowCommand(test.want.GetAllowCommand())
		test.secret.SetAllowSubstitution(test.want.GetAllowSubstitution())
		test.secret.SetCreatedAt(test.want.GetCreatedAt())
		test.secret.SetCreatedBy(test.want.GetCreatedBy())
		test.secret.SetUpdatedAt(test.want.GetUpdatedAt())
		test.secret.SetUpdatedBy(test.want.GetUpdatedBy())

		if test.secret.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.secret.GetID(), test.want.GetID())
		}

		if test.secret.GetOrg() != test.want.GetOrg() {
			t.Errorf("SetOrg is %v, want %v", test.secret.GetOrg(), test.want.GetOrg())
		}

		if test.secret.GetRepo() != test.want.GetRepo() {
			t.Errorf("SetRepo is %v, want %v", test.secret.GetRepo(), test.want.GetRepo())
		}

		if test.secret.GetTeam() != test.want.GetTeam() {
			t.Errorf("SetTeam is %v, want %v", test.secret.GetTeam(), test.want.GetTeam())
		}

		if test.secret.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.secret.GetName(), test.want.GetName())
		}

		if test.secret.GetValue() != test.want.GetValue() {
			t.Errorf("SetValue is %v, want %v", test.secret.GetValue(), test.want.GetValue())
		}

		if test.secret.GetType() != test.want.GetType() {
			t.Errorf("SetType is %v, want %v", test.secret.GetType(), test.want.GetType())
		}

		if !reflect.DeepEqual(test.secret.GetImages(), test.want.GetImages()) {
			t.Errorf("SetImages is %v, want %v", test.secret.GetImages(), test.want.GetImages())
		}

		if !reflect.DeepEqual(test.secret.GetAllowEvents(), test.want.GetAllowEvents()) {
			t.Errorf("SetAllowEvents is %v, want %v", test.secret.GetAllowEvents(), test.want.GetAllowEvents())
		}

		if test.secret.GetAllowCommand() != test.want.GetAllowCommand() {
			t.Errorf("SetAllowCommand is %v, want %v", test.secret.GetAllowCommand(), test.want.GetAllowCommand())
		}

		if test.secret.GetAllowSubstitution() != test.want.GetAllowSubstitution() {
			t.Errorf("SetAllowSubstitution is %v, want %v", test.secret.GetAllowSubstitution(), test.want.GetAllowSubstitution())
		}

		if test.secret.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("SetCreatedAt is %v, want %v", test.secret.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.secret.GetCreatedBy() != test.want.GetCreatedBy() {
			t.Errorf("SetCreatedBy is %v, want %v", test.secret.GetCreatedBy(), test.want.GetCreatedBy())
		}

		if test.secret.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("SetUpdatedAt is %v, want %v", test.secret.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.secret.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("SetUpdatedBy is %v, want %v", test.secret.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_Secret_String(t *testing.T) {
	// setup types
	s := testSecret()

	want := fmt.Sprintf(`{
	AllowCommand: %t,
	AllowEvents: %v,
	AllowSubstitution: %t,
	ID: %d,
	Images: %s,
	Name: %s,
	Org: %s,
	Repo: %s,
	Team: %s,
	Type: %s,
	Value: %s,
	CreatedAt: %d,
	CreatedBy: %s,
	UpdatedAt: %d,
	UpdatedBy: %s,
}`,
		s.GetAllowCommand(),
		s.GetAllowEvents().List(),
		s.GetAllowSubstitution(),
		s.GetID(),
		s.GetImages(),
		s.GetName(),
		s.GetOrg(),
		s.GetRepo(),
		s.GetTeam(),
		s.GetType(),
		s.GetValue(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
	)

	// run test
	got := s.String()

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("String Mismatch: -want +got):\n%s", diff)
	}
}

// testSecret is a test helper function to create a Secret
// type with all fields set to a fake value.
func testSecret() *Secret {
	currentTime := time.Now()
	tsCreate := currentTime.UTC().Unix()
	tsUpdate := currentTime.Add(time.Hour * 1).UTC().Unix()
	s := new(Secret)

	s.SetID(1)
	s.SetOrg("github")
	s.SetRepo("octocat")
	s.SetTeam("octokitties")
	s.SetName("foo")
	s.SetValue("bar")
	s.SetType("repo")
	s.SetImages([]string{"alpine"})
	s.SetAllowEvents(NewEventsFromMask(1))
	s.SetAllowCommand(true)
	s.SetAllowSubstitution(true)
	s.SetCreatedAt(tsCreate)
	s.SetCreatedBy("octocat")
	s.SetUpdatedAt(tsUpdate)
	s.SetUpdatedBy("octocat2")

	return s
}
