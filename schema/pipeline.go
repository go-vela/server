// SPDX-License-Identifier: Apache-2.0

// This program utilizes the json/jsonschema tags on structs in compiler/types/yaml
// to generate the majority of the final jsonschema for a Vela pipeline.
//
// Some manual intervention is needed for custom types and/or custom unmarshaling
// that is in place. For reference, we use the mechanisms provided by the schema lib:
// https://github.com/invopop/jsonschema?tab=readme-ov-file#custom-type-definitions
// for hooking into the schema generation process. Some types will have a JSONSchema
// or JSONSchemaExtend method attached to handle these overrides.

package schema

import (
	"fmt"

	"github.com/invopop/jsonschema"

	types "github.com/go-vela/server/compiler/types/yaml/yaml"
)

// NewPipelineSchema generates the JSON schema object for a Vela pipeline configuration.
//
// The returned value can be marshaled into actual JSON.
func NewPipelineSchema() (*jsonschema.Schema, error) {
	ref := jsonschema.Reflector{
		ExpandedStruct: true,
	}
	s := ref.Reflect(types.Build{})

	// very unlikely scenario
	if s == nil {
		return nil, fmt.Errorf("schema generation failed")
	}

	s.Title = "Vela Pipeline Configuration"

	// allows folks to have other top level arbitrary
	// keys without validation errors
	s.AdditionalProperties = nil

	// apply Ruleset modification
	//
	// note: we have to do the modification here,
	// because the custom type hooks can't provide
	// access to the top level definitions, even if
	// they were already processed, so we have to
	// do it at this top level.
	modRulesetSchema(s)

	return s, nil
}

// modRulesetSchema applies modifications to the Ruleset definition.
//
// rules can currently live at ruleset level or nested within
// 'if' (default) or 'unless'. without changes the struct would
// only allow the nested version.
func modRulesetSchema(schema *jsonschema.Schema) {
	if schema.Definitions == nil {
		return
	}

	rules, hasRules := schema.Definitions["Rules"]
	ruleSet, hasRuleset := schema.Definitions["Ruleset"]

	// exit early if we don't have what we need
	if !hasRules || !hasRuleset {
		return
	}

	// create copies
	_rulesWithRuleset := *rules
	_ruleSet := *ruleSet

	// copy every property from Ruleset, other than `if` and `unless`
	for item := _ruleSet.Properties.Newest(); item != nil; item = item.Prev() {
		if item.Key != "if" && item.Key != "unless" {
			_rulesWithRuleset.Properties.Set(item.Key, item.Value)
		}
	}

	// create a new definition for Ruleset
	schema.Definitions["Ruleset"].AnyOf = []*jsonschema.Schema{
		&_ruleSet,
		&_rulesWithRuleset,
	}
	schema.Definitions["Ruleset"].Properties = nil
	schema.Definitions["Ruleset"].Type = ""
	schema.Definitions["Ruleset"].AdditionalProperties = nil
}
