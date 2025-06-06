// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"testing"

	"github.com/go-vela/server/compiler/types/raw"
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
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Comment: []string{"rerun"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Comment: []string{"rerun"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"deployment"}, Target: []string{"production"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "deployment", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "production"},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"deployment"}, Target: []string{"production"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "deployment", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "stage"},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"schedule"}, Target: []string{"weekly"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "schedule", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "weekly"},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Event: []string{"schedule"}, Target: []string{"weekly"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "", Event: "schedule", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "nightly"},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Status: []string{"success", "failure"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Status: "failure", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Sender: []string{"octocat"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "ok to test", Event: "push", Repo: "octocat/hello-world", Sender: "octocat", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// If with or operator
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{If: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		// Unless with and operator
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Operator: "and"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}, Operator: "and"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// Unless with or operator
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    false,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"}},
			data:    &RuleData{Branch: "dev", Comment: "rerun", Event: "pull_request", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Path: []string{"foo.txt", "/foo/bar.txt"}, Operator: "or"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "pull_request", Path: []string{}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:    true,
		},
		// Advanced Rulesets
		{
			ruleset: &Ruleset{
				If: Rules{
					Event:    []string{"push", "pull_request"},
					Tag:      []string{"release/*"},
					Operator: "or",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Event:    []string{"push", "pull_request"},
					Tag:      []string{"release/*"},
					Operator: "or",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Event:    []string{"push", "pull_request"},
					Tag:      []string{"release/*"},
					Operator: "or",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want: false,
		},
		// Eval
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "VELA_BUILD_AUTHOR == 'Octocat'",
					Operator: "and",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{"VELA_BUILD_AUTHOR": "Octocat"}},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "VELA_BUILD_AUTHOR == 'Octocat'",
					Operator: "and",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{"VELA_BUILD_AUTHOR": "test"}},
			want: false,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "VELA_MISSING_VAR == 'Octocat'",
					Operator: "and",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{}},
			want: false,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "VELA_BUILD_AUTHOR == 'Octocat'",
					Branch:   []string{"main"},
					Event:    []string{"push"},
					Operator: "and",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{"VELA_BUILD_AUTHOR": "Octocat"}},
			want: true,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "VELA_BUILD_AUTHOR == 'Octocat'",
					Branch:   []string{"main"},
					Event:    []string{"pull_request"},
					Operator: "and",
				},
			},
			data: &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{"VELA_BUILD_AUTHOR": "Octocat"}},
			want: false,
		},
		{
			ruleset: &Ruleset{
				If: Rules{
					Eval:     "1 + bad-eval",
					Branch:   []string{"main"},
					Event:    []string{"pull_request"},
					Operator: "and",
				},
			},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: "", Env: map[string]string{"VELA_BUILD_AUTHOR": "Octocat"}},
			want:    false,
			wantErr: true,
		},
		// Bad regexp
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"*-dev"}, Matcher: "regexp"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			wantErr: true,
		},
		{
			ruleset: &Ruleset{Unless: Rules{Branch: []string{"*-dev"}, Event: []string{"push"}, Operator: "or", Matcher: "regexp"}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			wantErr: true,
		},
		// Bad filepath pattern
		{
			ruleset: &Ruleset{If: Rules{Branch: []string{"\\"}}},
			data:    &RuleData{Branch: "main", Comment: "rerun", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			wantErr: true,
		},
	}

	// run test
	for i, test := range tests {
		got, err := test.data.Match(*test.ruleset)
		if err != nil {
			if !test.wantErr {
				t.Errorf("Ruleset Match test #%d for %s operator returned err: %s", i, test.ruleset.Operator, err)
			}
		} else {
			if test.wantErr {
				t.Errorf("Ruleset Match should have returned an error")
			}
		}

		if got != test.want {
			t.Errorf("Ruleset Match test #%d for %s operator is %v, want %v", i, test.ruleset.Operator, got, test.want)
		}
	}
}

func BenchmarkMatch_FullRuleset(b *testing.B) {
	// create a sample RuleData
	data := &RuleData{
		Branch:   "main",
		Comment:  "test comment",
		Event:    "push",
		Path:     []string{"path/to/file"},
		Repo:     "github.com/go-vela/server",
		Sender:   "user",
		Status:   "success",
		Tag:      "v1.0.0",
		Target:   "production",
		Label:    []string{"label1", "label2"},
		Instance: "instance1",
	}

	// create a sample Ruleset
	ruleset := Ruleset{
		If: Rules{
			Branch:   []string{"main"},
			Comment:  []string{"test comment"},
			Event:    []string{"push"},
			Path:     []string{"path/to/file"},
			Repo:     []string{"github.com/go-vela/server"},
			Sender:   []string{"user"},
			Status:   []string{"success"},
			Tag:      []string{"v1.0.0"},
			Target:   []string{"production"},
			Label:    []string{"label1", "label2"},
			Instance: []string{"instance1"},
			Operator: "and",
			Matcher:  "filepath",
		},
	}

	// run the benchmark
	for b.Loop() {
		_, err := data.Match(ruleset)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkMatch_FullRulesetEarlyAndExit(b *testing.B) {
	// create a sample RuleData
	data := &RuleData{
		Branch:   "main",
		Comment:  "test comment",
		Event:    "push",
		Path:     []string{"path/to/file"},
		Repo:     "github.com/go-vela/server",
		Sender:   "user",
		Status:   "success",
		Tag:      "v1.0.0",
		Target:   "production",
		Label:    []string{"label1", "label2"},
		Instance: "instance1",
	}

	// create a sample Ruleset
	ruleset := Ruleset{
		If: Rules{
			Branch:   []string{"dev"},
			Comment:  []string{"test comment"},
			Event:    []string{"push"},
			Path:     []string{"path/to/file"},
			Repo:     []string{"github.com/go-vela/server"},
			Sender:   []string{"user"},
			Status:   []string{"success"},
			Tag:      []string{"v1.0.0"},
			Target:   []string{"production"},
			Label:    []string{"label1", "label2"},
			Instance: []string{"instance1"},
			Operator: "and",
			Matcher:  "filepath",
		},
	}

	// run the benchmark
	for b.Loop() {
		_, err := data.Match(ruleset)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkMatch_FullRulesetEarlyOrExit(b *testing.B) {
	// create a sample RuleData
	data := &RuleData{
		Branch:   "main",
		Comment:  "test comment",
		Event:    "push",
		Path:     []string{"path/to/file"},
		Repo:     "github.com/go-vela/server",
		Sender:   "user",
		Status:   "success",
		Tag:      "v1.0.0",
		Target:   "production",
		Label:    []string{"label1", "label2"},
		Instance: "instance1",
	}

	// create a sample Ruleset
	ruleset := Ruleset{
		If: Rules{
			Branch:   []string{"main"},
			Comment:  []string{"test comment"},
			Event:    []string{"push"},
			Path:     []string{"path/to/file"},
			Repo:     []string{"github.com/go-vela/server"},
			Sender:   []string{"user"},
			Status:   []string{"success"},
			Tag:      []string{"v1.0.0"},
			Target:   []string{"production"},
			Label:    []string{"label1", "label2"},
			Instance: []string{"instance1"},
			Operator: "or",
			Matcher:  "filepath",
		},
	}

	// run the benchmark
	for b.Loop() {
		_, err := data.Match(ruleset)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
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
		rules *Rules
		data  *RuleData
		want  bool
	}{
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"refs/tags/20.*"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"[0-9][0-9].[0-9].[0-9][0-9].[0-9][0-9][0-9]"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+$"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/20.4.42.167", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.165-prod", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"^refs/tags/(\\d+\\.)+\\d+$"}, Matcher: "regex", Operator: "and"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.165-prod", Target: ""},
			want:  false,
		},
	}

	// run test
	for _, test := range tests {
		got, _ := test.data.MatchRules(*test.rules)

		if got != test.want {
			t.Errorf("Rules Match is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_Rules_Match(t *testing.T) {
	// setup types
	tests := []struct {
		rules *Rules
		data  *RuleData
		want  bool
	}{
		// and operator
		{
			rules: &Rules{Branch: []string{"main"}},
			data:  &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Branch: []string{"main"}},
			data:  &RuleData{Branch: "dev", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:  &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}},
			data:  &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		{
			rules: &Rules{Path: []string{"foob.txt"}},
			data:  &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		{
			rules: &Rules{Status: []string{"success", "failure"}},
			data:  &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"refs/tags/[0-9].*-prod"}},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/tags/2.4.42.167-prod", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"tag"}, Tag: []string{"path/to/thing/*/*"}},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "path/to/thing/stage/1.0.2-rc", Target: ""},
			want:  true,
		},
		// or operator
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"},
			data:  &RuleData{Branch: "main", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"},
			data:  &RuleData{Branch: "dev", Event: "push", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"},
			data:  &RuleData{Branch: "main", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Branch: []string{"main"}, Event: []string{"push"}, Operator: "or"},
			data:  &RuleData{Branch: "dev", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		{
			rules: &Rules{Path: []string{"foob.txt"}, Operator: "or"},
			data:  &RuleData{Branch: "dev", Event: "pull_request", Path: []string{"foo.txt", "/foo/bar.txt"}, Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		// Advanced Rulesets
		{
			rules: &Rules{Event: []string{"push", "pull_request"}, Tag: []string{"release/*"}, Operator: "or"},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Tag: "release/*", Target: ""},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"push", "pull_request"}, Tag: []string{"release/*"}, Operator: "or"},
			data:  &RuleData{Branch: "main", Event: "tag", Repo: "octocat/hello-world", Status: "pending", Tag: "refs/heads/main", Target: ""},
			want:  false,
		},
		{
			rules: &Rules{Event: []string{"pull_request:labeled"}, Label: []string{"enhancement", "documentation"}},
			data:  &RuleData{Branch: "main", Event: "pull_request:labeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"pull_request:labeled"}, Label: []string{"enhancement", "documentation"}},
			data:  &RuleData{Branch: "main", Event: "pull_request:labeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"support"}},
			want:  false,
		},
		{
			rules: &Rules{Event: []string{"pull_request:unlabeled"}, Label: []string{"enhancement", "documentation"}},
			data:  &RuleData{Branch: "main", Event: "pull_request:unlabeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"pull_request:unlabeled"}, Label: []string{"enhancement"}},
			data:  &RuleData{Branch: "main", Event: "pull_request:unlabeled", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  false,
		},
		{
			rules: &Rules{Event: []string{"push"}, Label: []string{"enhancement", "documentation"}},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"push"}, Label: []string{"enhancement"}},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  false,
		},
		{
			rules: &Rules{Event: []string{"push"}, Label: []string{"enhancement"}, Operator: "or"},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Label: []string{"documentation"}},
			want:  true,
		},
		{
			rules: &Rules{Event: []string{"push"}, Instance: []string{"http://localhost:8080"}},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Instance: "http://localhost:5432"},
			want:  false,
		},
		{
			rules: &Rules{Event: []string{"push"}, Instance: []string{"http://localhost:8080"}},
			data:  &RuleData{Branch: "main", Event: "push", Repo: "octocat/hello-world", Status: "pending", Instance: "http://localhost:8080"},
			want:  true,
		},
	}

	// run test
	for i, test := range tests {
		got, _ := test.data.MatchRules(*test.rules)

		if got != test.want {
			t.Errorf("Rules Match for test #%d is %v, want %v", i, got, test.want)
		}
	}
}
