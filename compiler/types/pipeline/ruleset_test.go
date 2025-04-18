// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"testing"

	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/constants"
)

func TestPipeline_Ruleset_Match(t *testing.T) {
	// setup types
	tests := []struct {
		ruleset *Ruleset
		data    *RuleData
		envs    raw.StringSliceMap
		want    bool
		wantErr bool
	}{
		// Empty
		{ruleset: &Ruleset{}, data: &RuleData{Branch: "main"}, want: true},
		// If with and operator
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Comment: []string{"rerun"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Comment: []string{"rerun"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"deployment"}, Target: []string{"production"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "deployment", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "production"},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"deployment"}, Target: []string{"production"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "deployment", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "stage"},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"schedule"}, Target: []string{"weekly"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "schedule", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "weekly"},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"schedule"}, Target: []string{"weekly"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "schedule", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "nightly"},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Status: []string{"success", "failure"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Status: "failure", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Sender: []string{"octocat"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Sender: "octocat", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// If with or operator
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		// Unless with and operator
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}}, Operator: "and"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}}, Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// Unless with or operator
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}}, Operator: "or"},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}}, Operator: "or"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// Advanced Rulesets
		{
			ruleset: &Ruleset{
				If: Rules{
					Event: []string{"push", "pull_request"},
					Tag:   []string{"release/*"},
				},
				Operator: "or",
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Event: []string{"push", "pull_request"},
					Tag:   []string{"release/*"},
				},
				Operator: "or",
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Event: []string{"push", "pull_request"},
					Tag:   []string{"release/*"},
				},
				Operator: "or",
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want: false,
		},
		// Bad regexp
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"*-dev"}}, Matcher: "regexp"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			wantErr: true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"*-dev"}, Event: []string{"push"}}, Operator: "or", Matcher: "regexp"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			wantErr: true,
		},
		// Eval
		{
			ruleset: &Ruleset{Eval: "VELA_BUILD_AUTHOR == 'Octocat'", Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			envs: map[string]string{
				"VELA_BUILD_AUTHOR": "Octocat",
			},
			want:    true,
			wantErr: false,
		},
		{
			ruleset: &Ruleset{Eval: "VELA_BUILD_AUTHOR == 'Octocat'", Operator: "and"},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			envs: map[string]string{
				"VELA_BUILD_AUTHOR": "test",
			},
			want:    false,
			wantErr: false,
		},
	}

	// run test
	for _, test := range tests {
		got, err := test.ruleset.Match(test.data, test.envs)
		if err != nil {
			if !test.wantErr {
				t.Errorf("Ruleset Match for %s operator returned err: %s", test.ruleset.Operator, err)
			}
		} else {
			if test.wantErr {
				t.Errorf("Ruleset Match should have returned an error")
			}
		}

		if got != test.want {
			t.Errorf("Ruleset Match for %s operator is %v, want %v", test.ruleset.Operator, got, test.want)
		}
	}
}

func TestPipeline_Rules_NoStatus(t *testing.T) {
	// setup types
	r := Rules{}

	// run test
	got := r.Empty()

	if !got {
		t.Errorf("Rule NoStatus is %v, want true", got)
	}
}

func TestPipeline_Rules_Empty(t *testing.T) {
	// setup types
	r := Rules{}

	// run test
	got := r.Empty()

	if !got {
		t.Errorf("Rule IsEmpty is %v, want true", got)
	}
}

func TestPipeline_Rules_Empty_Invalid(t *testing.T) {
	// setup types
	r := Rules{Branch: []string{"main"}}

	// run test
	got := r.Empty()

	if got {
		t.Errorf("Rule IsEmpty is %v, want false", got)
	}
}

func TestPipeline_Rules_Match_Regex_Tag(t *testing.T) {
	// setup types
	tests := []struct {
		rules    *Rules
		data     *RuleData
		operator string
		want     bool
	}{
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"refs/tags/20.*"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"[0-9][0-9].[0-9].[0-9][0-9].[0-9][0-9][0-9]"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+$"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.165-prod", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+$"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.165-prod", Target: ""},
			operator: "and",
			want:     false,
		},
	}

	// run test
	for _, test := range tests {
		got, _ := test.rules.Match(test.data, "regexp", test.operator)

		if got != test.want {
			t.Errorf("Rules Match for %s operator is %v, want %v", test.operator, got, test.want)
		}
	}
}

func TestPipeline_Rules_Match(t *testing.T) {
	// setup types
	tests := []struct {
		rules    *Rules
		data     *RuleData
		operator string
		want     bool
	}{
		// Empty
		{
			rules:    &Rules{},
			data:     &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{},
			data:     &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     false,
		},
		// and operator
		{
			rules:    &Rules{Branch: []string{"main"}},
			data:     &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Branch: []string{"main"}},
			data:     &RuleData{Branch: "dev", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Path: []string{"foob.txt"}},
			data:     &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Status: []string{"success", "failure"}},
			data:     &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Tag: "refs/heads/main", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"refs/tags/[0-9].*-prod"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.167-prod", Target: ""},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"tag"}, Tag: []string{"path/to/thing/*/*"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "path/to/thing/stage/1.0.2-rc", Target: ""},
			operator: "and",
			want:     true,
		},
		// or operator
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     true,
		},
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "dev", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     true,
		},
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     true,
		},
		{
			rules:    &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:     &RuleData{Branch: "dev", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     false,
		},
		{
			rules:    &Rules{Path: []string{"foob.txt"}},
			data:     &RuleData{Branch: "dev", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     false,
		},
		// Advanced Rulesets
		{
			rules:    &Rules{Event: []string{"push", "pull_request"}, Tag: []string{"release/*"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			operator: "or",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"push", "pull_request"}, Tag: []string{"release/*"}},
			data:     &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			operator: "or",
			want:     false,
		},
		{
			rules:    &Rules{Event: []string{"pull_request:labeled"}, Label: []string{"enhancement", "documentation"}},
			data:     &RuleData{Branch: "main", Event: "pull_request:labeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"pull_request:labeled"}, Label: []string{"enhancement", "documentation"}},
			data:     &RuleData{Branch: "main", Event: "pull_request:labeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"support"}},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Event: []string{"pull_request:unlabeled"}, Label: []string{"enhancement", "documentation"}},
			data:     &RuleData{Branch: "main", Event: "pull_request:unlabeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"pull_request:unlabeled"}, Label: []string{"enhancement"}},
			data:     &RuleData{Branch: "main", Event: "pull_request:unlabeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Event: []string{"push"}, Label: []string{"enhancement", "documentation"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "and",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"push"}, Label: []string{"enhancement"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Event: []string{"push"}, Label: []string{"enhancement"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			operator: "or",
			want:     true,
		},
		{
			rules:    &Rules{Event: []string{"push"}, Instance: []string{"http://localhost:8080"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Instance: "http://localhost:5432"},
			operator: "and",
			want:     false,
		},
		{
			rules:    &Rules{Event: []string{"push"}, Instance: []string{"http://localhost:8080"}},
			data:     &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Instance: "http://localhost:8080"},
			operator: "and",
			want:     true,
		},
	}

	// run test
	for _, test := range tests {
		got, _ := test.rules.Match(test.data, "filepath", test.operator)

		if got != test.want {
			t.Errorf("Rules Match for %s operator is %v, want %v", test.operator, got, test.want)
		}
	}
}

func TestPipeline_Ruletype_MatchAnd(t *testing.T) {
	// setup types
	tests := []struct {
		matcher string
		rule    Ruletype
		pattern string
		want    bool
	}{
		// Empty with filepath matcher
		{matcher: "filepath", rule: []string{}, pattern: "main", want: true},
		{matcher: "filepath", rule: []string{}, pattern: "push", want: true},
		{matcher: "filepath", rule: []string{}, pattern: "foo/bar", want: true},
		{matcher: "filepath", rule: []string{}, pattern: "success", want: true},
		{matcher: "filepath", rule: []string{}, pattern: "release/*", want: true},
		// Branch with filepath matcher
		{matcher: "filepath", rule: []string{"main"}, pattern: "main", want: true},
		{matcher: "filepath", rule: []string{"main"}, pattern: "dev", want: false},
		// Comment with filepath matcher
		{matcher: "filepath", rule: []string{"ok to test"}, pattern: "ok to test", want: true},
		{matcher: "filepath", rule: []string{"ok to test"}, pattern: "rerun", want: false},
		// Event with filepath matcher
		{matcher: "filepath", rule: []string{"push"}, pattern: "push", want: true},
		{matcher: "filepath", rule: []string{"push"}, pattern: "pull_request", want: false},
		// Repo with filepath matcher
		{matcher: "filepath", rule: []string{"foo/bar"}, pattern: "foo/bar", want: true},
		{matcher: "filepath", rule: []string{"foo/bar"}, pattern: "test/foobar", want: false},
		// Status with filepath matcher
		{matcher: "filepath", rule: []string{"success"}, pattern: "success", want: true},
		{matcher: "filepath", rule: []string{"success"}, pattern: "failure", want: false},
		// Tag with filepath matcher
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/*", want: true},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "stage/*", want: false},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/111.2.3-rc", want: true},
		{matcher: "filepath", rule: []string{"release/**"}, pattern: "release/1.2.3-rc-hold", want: true},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/stage/1.2.3-rc", want: false},
		{matcher: "filepath", rule: []string{"release/*/*"}, pattern: "release/stage/1.2.3-rc", want: true},
		{matcher: "filepath", rule: []string{"release/stage/*"}, pattern: "release/stage/1.2.3-rc", want: true},
		{matcher: "filepath", rule: []string{"release/prod/*"}, pattern: "release/stage/1.2.3-rc", want: false},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/1.2.3-rc", want: true},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/1.2.3", want: true},
		// Target with filepath matcher
		{matcher: "filepath", rule: []string{"production"}, pattern: "production", want: true},
		{matcher: "filepath", rule: []string{"stage"}, pattern: "production", want: false},
		// Label with filepath matcher
		{matcher: "filepath", rule: []string{"enhancement", "documentation"}, pattern: "documentation", want: true},
		{matcher: "filepath", rule: []string{"enhancement", "documentation"}, pattern: "question", want: false},
		// Empty with regex matcher
		{matcher: "regexp", rule: []string{}, pattern: "main", want: true},
		{matcher: "regexp", rule: []string{}, pattern: "push", want: true},
		{matcher: "regexp", rule: []string{}, pattern: "foo/bar", want: true},
		{matcher: "regexp", rule: []string{}, pattern: "success", want: true},
		{matcher: "regexp", rule: []string{}, pattern: "release/*", want: true},
		// Branch with regex matcher
		{matcher: "regexp", rule: []string{"main"}, pattern: "main", want: true},
		{matcher: "regexp", rule: []string{"main"}, pattern: "dev", want: false},
		// Comment with regex matcher
		{matcher: "regexp", rule: []string{"ok to test"}, pattern: "ok to test", want: true},
		{matcher: "regexp", rule: []string{"ok to test"}, pattern: "rerun", want: false},
		// Event with regex matcher
		{matcher: "regexp", rule: []string{"push"}, pattern: "push", want: true},
		{matcher: "regexp", rule: []string{"push"}, pattern: "pull_request", want: false},
		// Repo with regex matcher
		{matcher: "regexp", rule: []string{"foo/bar"}, pattern: "foo/bar", want: true},
		{matcher: "regexp", rule: []string{"foo/bar"}, pattern: "test/foobar", want: false},
		// Status with regex matcher
		{matcher: "regexp", rule: []string{"success"}, pattern: "success", want: true},
		{matcher: "regexp", rule: []string{"success"}, pattern: "failure", want: false},
		// Tag with regex matcher
		{matcher: "regexp", rule: []string{"release/*"}, pattern: "release/*", want: true},
		{matcher: "regexp", rule: []string{"release/*"}, pattern: "stage/*", want: false},
		{matcher: "regex", rule: []string{"release/[0-9]+.*-rc$"}, pattern: "release/111.2.3-rc", want: true},
		{matcher: "regex", rule: []string{"release/[0-9]+.*-rc$"}, pattern: "release/1.2.3-rc-hold", want: false},
		{matcher: "regexp", rule: []string{"release/*"}, pattern: "release/stage/1.2.3-rc", want: true},
		{matcher: "regexp", rule: []string{"release/*/*"}, pattern: "release/stage/1.2.3-rc", want: true},
		{matcher: "regex", rule: []string{"release/stage/*"}, pattern: "release/stage/1.2.3-rc", want: true},
		{matcher: "regex", rule: []string{"release/prod/*"}, pattern: "release/stage/1.2.3-rc", want: false},
		{matcher: "regexp", rule: []string{"release/[0-9]+.[0-9]+.[0-9]+$"}, pattern: "release/1.2.3-rc", want: false},
		{matcher: "regexp", rule: []string{"release/[0-9]+.[0-9]+.[0-9]+$"}, pattern: "release/1.2.3", want: true},
		// Target with regex matcher
		{matcher: "regexp", rule: []string{"production"}, pattern: "production", want: true},
		{matcher: "regexp", rule: []string{"stage"}, pattern: "production", want: false},
		// Label with regexp matcher
		{matcher: "regexp", rule: []string{"enhancement", "documentation"}, pattern: "documentation", want: true},
		{matcher: "regexp", rule: []string{"enhancement", "documentation"}, pattern: "question", want: false},
		// Instance with regexp matcher
		{matcher: "regexp", rule: []string{"http://localhost:8080", "http://localhost:1234"}, pattern: "http://localhost:5432", want: false},
		{matcher: "regexp", rule: []string{"http://localhost:8080", "http://localhost:1234"}, pattern: "http://localhost:8080", want: true},
	}

	// run test
	for _, test := range tests {
		got, _ := test.rule.MatchSingle(test.pattern, test.matcher, constants.OperatorAnd)

		if got != test.want {
			t.Errorf("MatchAnd for %s matcher is %v, want %v", test.matcher, got, test.want)
		}
	}
}

func TestPipeline_Ruletype_MatchOr(t *testing.T) {
	// setup types
	tests := []struct {
		matcher string
		rule    Ruletype
		pattern string
		want    bool
	}{
		// Empty with filepath matcher
		{matcher: "filepath", rule: []string{}, pattern: "main", want: false},
		{matcher: "filepath", rule: []string{}, pattern: "push", want: false},
		{matcher: "filepath", rule: []string{}, pattern: "foo/bar", want: false},
		{matcher: "filepath", rule: []string{}, pattern: "success", want: false},
		{matcher: "filepath", rule: []string{}, pattern: "release/*", want: false},
		// Branch with filepath matcher
		{matcher: "filepath", rule: []string{"main"}, pattern: "main", want: true},
		{matcher: "filepath", rule: []string{"main"}, pattern: "dev", want: false},
		// Comment with filepath matcher
		{matcher: "filepath", rule: []string{"ok to test"}, pattern: "ok to test", want: true},
		{matcher: "filepath", rule: []string{"ok to test"}, pattern: "rerun", want: false},
		// Event with filepath matcher
		{matcher: "filepath", rule: []string{"push"}, pattern: "push", want: true},
		{matcher: "filepath", rule: []string{"push"}, pattern: "pull_request", want: false},
		// Repo with filepath matcher
		{matcher: "filepath", rule: []string{"foo/bar"}, pattern: "foo/bar", want: true},
		{matcher: "filepath", rule: []string{"foo/bar"}, pattern: "test/foobar", want: false},
		// Status with filepath matcher
		{matcher: "filepath", rule: []string{"success"}, pattern: "success", want: true},
		{matcher: "filepath", rule: []string{"success"}, pattern: "failure", want: false},
		// Tag with filepath matcher
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "release/*", want: true},
		{matcher: "filepath", rule: []string{"release/*"}, pattern: "stage/*", want: false},
		// Target with filepath matcher
		{matcher: "filepath", rule: []string{"production"}, pattern: "production", want: true},
		{matcher: "filepath", rule: []string{"stage"}, pattern: "production", want: false},
		// Label with filepath matcher
		{matcher: "filepath", rule: []string{"enhancement", "documentation"}, pattern: "documentation", want: true},
		{matcher: "filepath", rule: []string{"enhancement", "documentation"}, pattern: "question", want: false},
		// Empty with regexp matcher
		{matcher: "regexp", rule: []string{}, pattern: "main", want: false},
		{matcher: "regexp", rule: []string{}, pattern: "push", want: false},
		{matcher: "regexp", rule: []string{}, pattern: "foo/bar", want: false},
		{matcher: "regexp", rule: []string{}, pattern: "success", want: false},
		{matcher: "regexp", rule: []string{}, pattern: "release/*", want: false},
		// Branch with regexp matcher
		{matcher: "regexp", rule: []string{"main"}, pattern: "main", want: true},
		{matcher: "regexp", rule: []string{"main"}, pattern: "dev", want: false},
		// Comment with regexp matcher
		{matcher: "regexp", rule: []string{"ok to test"}, pattern: "ok to test", want: true},
		{matcher: "regexp", rule: []string{"ok to test"}, pattern: "rerun", want: false},
		// Event with regexp matcher
		{matcher: "regexp", rule: []string{"push"}, pattern: "push", want: true},
		{matcher: "regexp", rule: []string{"push"}, pattern: "pull_request", want: false},
		// Repo with regexp matcher
		{matcher: "regexp", rule: []string{"foo/bar"}, pattern: "foo/bar", want: true},
		{matcher: "regexp", rule: []string{"foo/bar"}, pattern: "test/foobar", want: false},
		// Status with regexp matcher
		{matcher: "regexp", rule: []string{"success"}, pattern: "success", want: true},
		{matcher: "regexp", rule: []string{"success"}, pattern: "failure", want: false},
		// Tag with regexp matcher
		{matcher: "regexp", rule: []string{"release/*"}, pattern: "release/*", want: true},
		{matcher: "regexp", rule: []string{"release/*"}, pattern: "stage/*", want: false},
		// Target with regexp matcher
		{matcher: "regexp", rule: []string{"production"}, pattern: "production", want: true},
		{matcher: "regexp", rule: []string{"stage"}, pattern: "production", want: false},
		// Label with regexp matcher
		{matcher: "regexp", rule: []string{"enhancement", "documentation"}, pattern: "documentation", want: true},
		{matcher: "regexp", rule: []string{"enhancement", "documentation"}, pattern: "question", want: false},
		// Instance with regexp matcher
		{matcher: "regexp", rule: []string{"http://localhost:8080", "http://localhost:1234"}, pattern: "http://localhost:5432", want: false},
		{matcher: "regexp", rule: []string{"http://localhost:8080", "http://localhost:1234"}, pattern: "http://localhost:8080", want: true},
	}

	// run test
	for _, test := range tests {
		got, _ := test.rule.MatchSingle(test.pattern, test.matcher, constants.OperatorOr)

		if got != test.want {
			t.Errorf("MatchOr for %s matcher is %v, want %v", test.matcher, got, test.want)
		}
	}
}
