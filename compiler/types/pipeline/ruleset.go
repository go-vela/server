// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/expr-lang/expr"

	"github.com/go-vela/server/compiler/types/raw"
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
		Eval     string   `json:"eval,omitempty"     yaml:"eval,omitempty"`
		Operator string   `json:"operator,omitempty" yaml:"operator,omitempty"`
		Matcher  string   `json:"matcher,omitempty"  yaml:"matcher,omitempty"`
	}

	// Ruletype is the pipeline representation of an element
	// for a ruleset block for a step in a pipeline.
	//
	// swagger:model PipelineRuletype
	Ruletype []string

	// RuleData is the data to check our ruleset
	// against for a step in a pipeline.
	RuleData struct {
		Branch   string             `json:"branch,omitempty"   yaml:"branch,omitempty"`
		Comment  string             `json:"comment,omitempty"  yaml:"comment,omitempty"`
		Event    string             `json:"event,omitempty"    yaml:"event,omitempty"`
		Path     []string           `json:"path,omitempty"     yaml:"path,omitempty"`
		Repo     string             `json:"repo,omitempty"     yaml:"repo,omitempty"`
		Sender   string             `json:"sender,omitempty"   yaml:"sender,omitempty"`
		Status   string             `json:"status,omitempty"   yaml:"status,omitempty"`
		Tag      string             `json:"tag,omitempty"      yaml:"tag,omitempty"`
		Target   string             `json:"target,omitempty"   yaml:"target,omitempty"`
		Label    []string           `json:"label,omitempty"    yaml:"label,omitempty"`
		Instance string             `json:"instance,omitempty" yaml:"instance,omitempty"`
		Env      raw.StringSliceMap `json:"env,omitempty"      yaml:"env,omitempty"`
	}
)

// Match checks if the context data of a given pipeline compilation matches the ruleset defined.
func (data *RuleData) Match(set Ruleset) (bool, error) {
	if set.If.Empty() && set.Unless.Empty() {
		return true, nil
	}

	if !set.Unless.Empty() {
		match, err := data.MatchRules(set.Unless)
		if err != nil {
			return false, err
		}

		if match {
			return false, nil
		}
	}

	if set.If.Empty() {
		return true, nil
	}

	match, err := data.MatchRules(set.If)
	if err != nil {
		return false, err
	}

	return match, nil
}

// MatchRules iterates through the defined rules in a ruleset and determines if the data matches.
func (data *RuleData) MatchRules(rules Rules) (bool, error) {
	isOr := strings.EqualFold(rules.Operator, constants.OperatorOr)

	var (
		match bool
		err   error
	)

	// build is being compiled - status does not matter
	if len(rules.Status) > 0 && data.Status == constants.StatusPending {
		return true, nil
	}

	// set running to success on ruleset evaluations (c.Execute)
	if len(rules.Status) > 0 && data.Status == constants.StatusRunning {
		data.Status = constants.StatusSuccess
	}

	// run eval first
	if len(rules.Eval) > 0 && data.Env != nil {
		eval, err := expr.Compile(rules.Eval, expr.Env(data.Env), expr.AllowUndefinedVariables(), expr.AsBool())
		if err != nil {
			return false, fmt.Errorf("failed to compile eval of %s: %w", rules.Eval, err)
		}

		result, err := expr.Run(eval, data.Env)
		if err != nil {
			return false, fmt.Errorf("failed to run eval of %s: %w", rules.Eval, err)
		}

		bResult, ok := result.(bool)
		if !ok {
			return false, fmt.Errorf("failed to parse eval of %s: expected bool but got %v", rules.Eval, bResult)
		}

		match = bResult

		// early exit if OR + truthy
		if isOr && match {
			return true, nil
		}

		// early exit if AND + falsy
		if !isOr && !match {
			return false, nil
		}
	}

	// define rule data type
	ruleSets := []struct {
		data []string
		rule Ruletype
	}{
		{[]string{data.Branch}, rules.Branch},
		{[]string{data.Comment}, rules.Comment},
		{[]string{data.Event}, rules.Event},
		{data.Path, rules.Path},
		{[]string{data.Repo}, rules.Repo},
		{[]string{data.Sender}, rules.Sender},
		{[]string{data.Status}, rules.Status},
		{[]string{data.Tag}, rules.Tag},
		{[]string{data.Target}, rules.Target},
		{data.Label, rules.Label},
		{[]string{data.Instance}, rules.Instance},
	}

	for _, rs := range ruleSets {
		// if user has specified no rule, continue
		if len(rs.rule) == 0 {
			continue
		}

		match, err = MatchRule(rs.data, rs.rule, rules.Matcher)
		if err != nil {
			return false, err
		}

		// early exit if OR + truthy
		if isOr && match {
			return true, nil
		}

		// early exit if AND + falsy
		if !isOr && !match {
			return false, nil
		}
	}

	return match, nil
}

// MatchRule determines the truthy value of a rule compared to incoming data given the matcher type.
func MatchRule(data []string, comparator []string, matcher string) (bool, error) {
	var err error

	// for each defined pattern by user
	for _, c := range comparator {
		// if matcher is regex, compile the pattern
		var pattern *regexp.Regexp
		if matcher == constants.MatcherRegex || matcher == "regex" {
			pattern, err = regexp.Compile(c)
			if err != nil {
				return false, fmt.Errorf("error in regex pattern %s: %w", c, err)
			}
		}

		// for each rule data value
		for _, d := range data {
			switch matcher {
			case constants.MatcherRegex, "regex":
				if pattern.MatchString(d) {
					return true, nil
				}
			default:
				match, err := filepath.Match(c, d)
				if err != nil {
					return false, err
				}

				if match {
					return true, nil
				}
			}
		}
	}

	return false, nil
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
		len(r.Instance) == 0 &&
		len(r.Eval) == 0 {
		return true
	}

	// return false if any of the ruletype is provided
	return false
}
