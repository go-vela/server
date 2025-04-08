// SPDX-License-Identifier: Apache-2.0

package buildkite

import (
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
)

type (
	// Ruleset is the yaml representation of a
	// ruleset block for a step in a pipeline.
	Ruleset struct {
		If       Rules  `yaml:"if,omitempty"       json:"if,omitempty"       jsonschema:"description=Limit execution to when all rules match.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Unless   Rules  `yaml:"unless,omitempty"   json:"unless,omitempty"   jsonschema:"description=Limit execution to when all rules do not match.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Matcher  string `yaml:"matcher,omitempty"  json:"matcher,omitempty"  jsonschema:"enum=filepath,enum=regexp,default=filepath,description=Use the defined matching method.\nReference: coming soon"`
		Operator string `yaml:"operator,omitempty" json:"operator,omitempty" jsonschema:"enum=or,enum=and,default=and,description=Whether all rule conditions must be met or just any one of them.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Continue bool   `yaml:"continue,omitempty" json:"continue,omitempty" jsonschema:"default=false,description=Limits the execution of a step to continuing on any failure.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
	}

	// Rules is the yaml representation of the ruletypes
	// from a ruleset block for a step in a pipeline.
	Rules struct {
		Branch   []string `yaml:"branch,omitempty,flow"   json:"branch,omitempty"   jsonschema:"description=Limits the execution of a step to matching build branches.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Comment  []string `yaml:"comment,omitempty,flow"  json:"comment,omitempty"  jsonschema:"description=Limits the execution of a step to matching a pull request comment.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Event    []string `yaml:"event,omitempty,flow"    json:"event,omitempty"    jsonschema:"description=Limits the execution of a step to matching build events.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Path     []string `yaml:"path,omitempty,flow"     json:"path,omitempty"     jsonschema:"description=Limits the execution of a step to matching files changed in a repository.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Repo     []string `yaml:"repo,omitempty,flow"     json:"repo,omitempty"     jsonschema:"description=Limits the execution of a step to matching repos.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Sender   []string `yaml:"sender,omitempty,flow"   json:"sender,omitempty"   jsonschema:"description=Limits the execution of a step to matching build senders.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Status   []string `yaml:"status,omitempty,flow"   json:"status,omitempty"   jsonschema:"enum=[failure],enum=[success],description=Limits the execution of a step to matching build statuses.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Tag      []string `yaml:"tag,omitempty,flow"      json:"tag,omitempty"      jsonschema:"description=Limits the execution of a step to matching build tag references.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Target   []string `yaml:"target,omitempty,flow"   json:"target,omitempty"   jsonschema:"description=Limits the execution of a step to matching build deployment targets.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Label    []string `yaml:"label,omitempty,flow"    json:"label,omitempty"    jsonschema:"description=Limits step execution to match on pull requests labels.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Instance []string `yaml:"instance,omitempty,flow" json:"instance,omitempty" jsonschema:"description=Limits step execution to match on certain instances.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Operator string   `yaml:"operator,omitempty"      json:"operator,omitempty" jsonschema:"description=Whether all rule conditions must be met or just any one of them.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
		Matcher  string   `yaml:"matcher,omitempty"       json:"matcher,omitempty"  jsonschema:"description=Use the defined matching method.\nReference: https://go-vela.github.io/docs/reference/yaml/steps/#the-ruleset-key"`
	}
)

// ToPipeline converts the Ruleset type
// to a pipeline Ruleset type.
func (r *Ruleset) ToPipeline() *pipeline.Ruleset {
	return &pipeline.Ruleset{
		If:      *r.If.ToPipeline(),
		Unless:  *r.Unless.ToPipeline(),
		Matcher: r.Matcher,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the Ruleset type.
func (r *Ruleset) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// simple struct we try unmarshalling to
	simple := new(Rules)

	// advanced struct we try unmarshalling to
	advanced := new(struct {
		If       Rules
		Unless   *Rules
		Matcher  string
		Operator string
		Continue bool
	})

	// attempt to unmarshal simple ruleset
	//nolint:errcheck // intentionally not handling error
	unmarshal(simple)
	// attempt to unmarshal advanced ruleset
	//nolint:errcheck // intentionally not handling error
	unmarshal(advanced)

	// set ruleset `unless` to advanced `unless` rules if they were parsed
	if advanced.Unless != nil {
		r.Unless = *advanced.Unless
	}
	// parse ruleset matcher and set to default if empty
	r.Matcher = advanced.Matcher
	if r.Matcher == "" {
		r.Matcher = constants.MatcherFilepath
	}
	// parse ruleset operator and set to default if empty
	r.Operator = advanced.Operator
	if r.Operator == "" {
		r.Operator = constants.OperatorAnd
	}
	// set ruleset `continue` to advanced `continue`
	r.Continue = advanced.Continue

	// implicitly add simple ruleset to the advanced ruleset for each rule type
	advanced.If.Branch = append(advanced.If.Branch, simple.Branch...)
	advanced.If.Comment = append(advanced.If.Comment, simple.Comment...)
	advanced.If.Event = append(advanced.If.Event, simple.Event...)
	advanced.If.Path = append(advanced.If.Path, simple.Path...)
	advanced.If.Repo = append(advanced.If.Repo, simple.Repo...)
	advanced.If.Sender = append(advanced.If.Sender, simple.Sender...)
	advanced.If.Status = append(advanced.If.Status, simple.Status...)
	advanced.If.Tag = append(advanced.If.Tag, simple.Tag...)
	advanced.If.Target = append(advanced.If.Target, simple.Target...)
	advanced.If.Label = append(advanced.If.Label, simple.Label...)
	advanced.If.Instance = append(advanced.If.Instance, simple.Instance...)

	// set ruleset `if` to advanced `if` rules
	r.If = advanced.If

	// inherit Ruleset operator/matcher if none specified
	if r.If.Operator == "" {
		r.If.Operator = r.Operator
	}

	if r.If.Matcher == "" {
		r.If.Matcher = r.Matcher
	}

	if advanced.Unless != nil {
		if r.Unless.Operator == "" {
			r.Unless.Operator = r.Operator
		}

		if r.Unless.Matcher == "" {
			r.Unless.Matcher = r.Matcher
		}
	}

	// zero out Ruleset level operator/matcher
	r.Operator = ""
	r.Matcher = ""

	return nil
}

// ToPipeline converts the Rules
// type to a pipeline Rules type.
func (r *Rules) ToPipeline() *pipeline.Rules {
	return &pipeline.Rules{
		Branch:   r.Branch,
		Comment:  r.Comment,
		Event:    r.Event,
		Path:     r.Path,
		Repo:     r.Repo,
		Sender:   r.Sender,
		Status:   r.Status,
		Tag:      r.Tag,
		Target:   r.Target,
		Label:    r.Label,
		Instance: r.Instance,
		Operator: r.Operator,
		Matcher:  r.Matcher,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the Rules type.
func (r *Rules) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// rules struct we try unmarshalling to
	rules := new(struct {
		Branch   raw.StringSlice
		Comment  raw.StringSlice
		Event    raw.StringSlice
		Path     raw.StringSlice
		Repo     raw.StringSlice
		Sender   raw.StringSlice
		Status   raw.StringSlice
		Tag      raw.StringSlice
		Target   raw.StringSlice
		Label    raw.StringSlice
		Instance raw.StringSlice
		Operator string
		Matcher  string
	})

	// attempt to unmarshal rules
	err := unmarshal(rules)
	if err == nil {
		r.Branch = rules.Branch
		r.Comment = rules.Comment
		r.Path = rules.Path
		r.Repo = rules.Repo
		r.Sender = rules.Sender
		r.Status = rules.Status
		r.Tag = rules.Tag
		r.Target = rules.Target
		r.Label = rules.Label
		r.Instance = rules.Instance
		r.Operator = rules.Operator
		r.Matcher = rules.Matcher

		// account for users who use non-scoped pull_request event
		events := []string{}

		for _, e := range rules.Event {
			switch e {
			// backwards compatibility
			// pull_request = pull_request:opened + pull_request:synchronize + pull_request:reopened
			// comment = comment:created + comment:edited
			case constants.EventPull:
				events = append(events,
					constants.EventPull+":"+constants.ActionOpened,
					constants.EventPull+":"+constants.ActionSynchronize,
					constants.EventPull+":"+constants.ActionReopened)
			case constants.EventDeploy:
				events = append(events,
					constants.EventDeploy+":"+constants.ActionCreated)
			case constants.EventComment:
				events = append(events,
					constants.EventComment+":"+constants.ActionCreated,
					constants.EventComment+":"+constants.ActionEdited)
			default:
				events = append(events, e)
			}
		}

		r.Event = events
	}

	return err
}

func (r *Ruleset) ToYAML() *yaml.Ruleset {
	if r == nil {
		return nil
	}

	return &yaml.Ruleset{
		If:       *r.If.ToYAML(),
		Unless:   *r.Unless.ToYAML(),
		Matcher:  r.Matcher,
		Operator: r.Operator,
		Continue: r.Continue,
	}
}

func (r *Rules) ToYAML() *yaml.Rules {
	if r == nil {
		return nil
	}

	return &yaml.Rules{
		Branch:   r.Branch,
		Comment:  r.Comment,
		Event:    r.Event,
		Path:     r.Path,
		Repo:     r.Repo,
		Sender:   r.Sender,
		Status:   r.Status,
		Tag:      r.Tag,
		Target:   r.Target,
		Label:    r.Label,
		Instance: r.Instance,
	}
}
