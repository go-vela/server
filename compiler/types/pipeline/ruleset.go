// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-vela/server/compiler/types/raw"

	"github.com/expr-lang/expr"

	"github.com/go-vela/server/constants"
)

type (
	// Ruleset is the pipeline representation of
	// a ruleset block for a step in a pipeline.
	//
	// swagger:model PipelineRuleset
	Ruleset struct {
		If       Rules  `json:"if,omitempty"       yaml:"if,omitempty"`
		Unless   Rules  `json:"unless,omitempty"   yaml:"unless,omitempty"`
		Matcher  string `json:"matcher,omitempty"  yaml:"matcher,omitempty"`
		Operator string `json:"operator,omitempty" yaml:"operator,omitempty"`
		Continue bool   `json:"continue,omitempty" yaml:"continue,omitempty"`
		Eval     string `json:"eval,omitempty"     yaml:"eval,omitempty"`
	}

	// Rules is the pipeline representation of the ruletypes
	// from a ruleset block for a step in a pipeline.
	//
	// swagger:model PipelineRules
	Rules struct {
		Branch   Ruletype `json:"branch,omitempty"   yaml:"branch,omitempty"`
		Comment  Ruletype `json:"comment,omitempty"  yaml:"comment,omitempty"`
		Event    Ruletype `json:"event,omitempty"    yaml:"event,omitempty"`
		Path     Ruletype `json:"path,omitempty"     yaml:"path,omitempty"`
		Repo     Ruletype `json:"repo,omitempty"     yaml:"repo,omitempty"`
		Sender   Ruletype `json:"sender,omitempty"   yaml:"sender,omitempty"`
		Status   Ruletype `json:"status,omitempty"   yaml:"status,omitempty"`
		Tag      Ruletype `json:"tag,omitempty"      yaml:"tag,omitempty"`
		Target   Ruletype `json:"target,omitempty"   yaml:"target,omitempty"`
		Label    Ruletype `json:"label,omitempty"    yaml:"label,omitempty"`
		Instance Ruletype `json:"instance,omitempty" yaml:"instance,omitempty"`
		Parallel bool     `json:"-"                  yaml:"-"`
	}

	// Ruletype is the pipeline representation of an element
	// for a ruleset block for a step in a pipeline.
	//
	// swagger:model PipelineRuletype
	Ruletype []string

	// RuleData is the data to check our ruleset
	// against for a step in a pipeline.
	RuleData struct {
		Branch   string   `json:"branch,omitempty"   yaml:"branch,omitempty"`
		Comment  string   `json:"comment,omitempty"  yaml:"comment,omitempty"`
		Event    string   `json:"event,omitempty"    yaml:"event,omitempty"`
		Path     []string `json:"path,omitempty"     yaml:"path,omitempty"`
		Repo     string   `json:"repo,omitempty"     yaml:"repo,omitempty"`
		Sender   string   `json:"sender,omitempty"   yaml:"sender,omitempty"`
		Status   string   `json:"status,omitempty"   yaml:"status,omitempty"`
		Tag      string   `json:"tag,omitempty"      yaml:"tag,omitempty"`
		Target   string   `json:"target,omitempty"   yaml:"target,omitempty"`
		Label    []string `json:"label,omitempty"    yaml:"label,omitempty"`
		Instance string   `json:"instance,omitempty" yaml:"instance,omitempty"`
		Parallel bool     `json:"-"                  yaml:"-"`
	}
)

// Match returns true when the provided ruledata matches
// the if rules and does not match any of the unless rules.
// When the provided if rules are empty, the function returns
// true. When both the provided if and unless rules are empty,
// the function also returns true.
func (r *Ruleset) Match(from *RuleData, envs raw.StringSliceMap) (bool, error) {
	// return true when the if and unless rules are empty
	if r.If.Empty() && r.Unless.Empty() && r.Eval == "" {
		return true, nil
	}

	// return false when the unless rules are not empty and match
	if !r.Unless.Empty() {
		match, err := r.Unless.Match(from, r.Matcher, r.Operator)
		if err != nil {
			return false, err
		}

		if match {
			return false, nil
		}
	}

	// return true when the if rules are empty
	if r.If.Empty() && r.Eval == "" {
		return true, nil
	}

	// return true when the if rules match
	match, err := r.If.Match(from, r.Matcher, r.Operator)
	if match && r.Eval != "" {
		eval, err := expr.Compile(r.Eval, expr.Env(envs), expr.AllowUndefinedVariables(), expr.AsBool())
		if err != nil {
			return false, fmt.Errorf("failed to compile eval of %s: %w", r.Eval, err)
		}

		result, err := expr.Run(eval, envs)
		if err != nil {
			return false, fmt.Errorf("failed to run eval of %s: %w", r.Eval, err)
		}

		bResult, ok := result.(bool)
		if !ok {
			return false, fmt.Errorf("failed to parse eval of %s: expected bool but got %v", r.Eval, bResult)
		}

		match = bResult
	}

	return match, err
}

// NoStatus returns true if the status field is empty.
func (r *Rules) NoStatus() bool {
	// return true if every ruletype is empty
	return len(r.Status) == 0
}

// Empty returns true if the provided ruletypes are empty.
func (r *Rules) Empty() bool {
	// return true if every ruletype is empty
	if len(r.Branch) == 0 &&
		len(r.Comment) == 0 &&
		len(r.Event) == 0 &&
		len(r.Path) == 0 &&
		len(r.Repo) == 0 &&
		len(r.Sender) == 0 &&
		len(r.Status) == 0 &&
		len(r.Tag) == 0 &&
		len(r.Target) == 0 &&
		len(r.Label) == 0 &&
		len(r.Instance) == 0 {
		return true
	}

	// return false if any of the ruletype is provided
	return false
}

// Match returns true for the `or` operator when one of the
// ruletypes from the rules match the provided ruledata.
// Match returns true for the `and` operator when all of the
// ruletypes from the rules match the provided ruledata. For
// both operators, when none of the ruletypes from the rules
// match the provided ruledata, the function returns false.
func (r *Rules) Match(from *RuleData, matcher, op string) (bool, error) {
	status, err := r.matchStatus(from, matcher, op)
	if err != nil {
		return false, err
	}

	matchBranch, err := r.Branch.MatchSingle(from.Branch, matcher, op)
	if err != nil {
		return false, err
	}

	matchComment, err := r.Comment.MatchSingle(from.Comment, matcher, op)
	if err != nil {
		return false, err
	}

	matchEvent, err := r.Event.MatchSingle(from.Event, matcher, op)
	if err != nil {
		return false, err
	}

	matchPath, err := r.Path.MatchMultiple(from.Path, matcher, op)
	if err != nil {
		return false, err
	}

	matchRepo, err := r.Repo.MatchSingle(from.Repo, matcher, op)
	if err != nil {
		return false, err
	}

	matchSender, err := r.Sender.MatchSingle(from.Sender, matcher, op)
	if err != nil {
		return false, err
	}

	matchTag, err := r.Tag.MatchSingle(from.Tag, matcher, op)
	if err != nil {
		return false, err
	}

	matchTarget, err := r.Target.MatchSingle(from.Target, matcher, op)
	if err != nil {
		return false, err
	}

	matchLabel, err := r.Label.MatchMultiple(from.Label, matcher, op)
	if err != nil {
		return false, err
	}

	matchInstance, err := r.Instance.MatchSingle(from.Instance, matcher, op)
	if err != nil {
		return false, err
	}

	return r.evaluateMatches(op, status, matchBranch, matchComment, matchEvent, matchPath, matchRepo, matchSender, matchTag, matchTarget, matchLabel, matchInstance), nil
}

func (r *Rules) matchStatus(from *RuleData, matcher, op string) (bool, error) {
	if len(from.Status) == 0 {
		return true, nil
	}

	return r.Status.MatchSingle(from.Status, matcher, op)
}

func (r *Rules) evaluateMatches(op string, matches ...bool) bool {
	switch op {
	case constants.OperatorOr:
		for _, match := range matches {
			if match {
				return true
			}
		}

		return false
	default:
		for _, match := range matches {
			if !match {
				return false
			}
		}

		return true
	}
}

// MatchSingle returns true when the provided ruletype
// matches the provided ruledata. When the provided
// ruletype is empty, the function returns true for
// the `and` operator and false for the `or` operator.
func (r *Ruletype) MatchSingle(data, matcher, logic string) (bool, error) {
	// return true for `and`, false for `or` if an empty ruletype is provided
	if len(*r) == 0 {
		return strings.EqualFold(logic, constants.OperatorAnd), nil
	}

	// iterate through each pattern in the ruletype
	for _, pattern := range *r {
		match, err := match(data, matcher, pattern)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	// return false if no match is found
	return false, nil
}

// MatchMultiple returns true when the provided ruletype
// matches the provided ruledata. When the provided
// ruletype is empty, the function returns true for
// the `and` operator and false for the `or` operator.
func (r *Ruletype) MatchMultiple(data []string, matcher, logic string) (bool, error) {
	// return true for `and`, false for `or` if an empty ruletype is provided
	if len(*r) == 0 {
		return strings.EqualFold(logic, constants.OperatorAnd), nil
	}

	// iterate through each pattern in the ruletype
	for _, pattern := range *r {
		for _, value := range data {
			match, err := match(value, matcher, pattern)
			if err != nil {
				return false, err
			}

			if match {
				return true, nil
			}
		}
	}

	// return false if no match is found
	return false, nil
}

// match is a helper function that compares data against a pattern
// and returns true if the data matches the pattern, depending on
// matcher specified.
func match(data, matcher, pattern string) (bool, error) {
	// handle the pattern based off the matcher provided
	switch matcher {
	case constants.MatcherRegex, "regex":
		regExpPattern, err := regexp.Compile(pattern)
		if err != nil {
			return false, fmt.Errorf("error in regex pattern %s: %w", pattern, err)
		}

		// return true if the regexp pattern matches the ruledata
		if regExpPattern.MatchString(data) {
			return true, nil
		}
	case constants.MatcherFilepath:
		fallthrough
	default:
		// return true if the pattern matches the ruledata
		ok, _ := filepath.Match(pattern, data)
		if ok {
			return true, nil
		}
	}

	// return false if no match is found
	return false, nil
}
