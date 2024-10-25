// SPDX-License-Identifier: Apache-2.0

//go:build ignore

// This program will utilize the json/jsonschema tags
// on structs in compiler/types/yaml to generate the
// majority of the final jsonschema for a
// Vela pipeline.
//
// Some manual intervention is needed for custom types
// and/or custom marshalling that is in place. For reference
// we use the provided mechanisms, see:
// https://github.com/invopop/jsonschema?tab=readme-ov-file#custom-type-definitions
// for hooking into the schema generation process. Some
// types will have a JSONSchema or JSONSchemaExtend method
// attached to handle the overrides.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/compiler/types/yaml"
)

func main() {
	ref := jsonschema.Reflector{
		ExpandedStruct: true,
	}
	d := ref.Reflect(yaml.Build{})

	if d == nil {
		logrus.Fatal("reflection failed")
	}

	d.Title = "Vela Pipeline Configuration"

	// allows folks to have other top level arbitrary
	// keys without validation errors
	d.AdditionalProperties = nil

	// rules can currently live at ruleset level or
	// nested within 'if' (default) or 'unless'.
	// without changes the struct would only allow
	// the nested version.
	//
	// note: we have to do the modification here,
	// because the custom type hooks can't provide
	// access to the top level definitions, even if
	// they were already processed, so we have to
	// do it at this top level.
	ruleSetWithRules := *d.Definitions["Rules"]
	ruleSet := *d.Definitions["Ruleset"]

	// copy every property from Ruleset, other than `if` and `unless`
	for item := ruleSet.Properties.Newest(); item != nil; item = item.Prev() {
		if item.Key != "if" && item.Key != "unless" {
			ruleSetWithRules.Properties.Set(item.Key, item.Value)
		}
	}

	// create a new definition for Ruleset
	d.Definitions["Ruleset"].AnyOf = []*jsonschema.Schema{
		&ruleSet,
		&ruleSetWithRules,
	}
	d.Definitions["Ruleset"].Properties = nil
	d.Definitions["Ruleset"].Type = ""
	d.Definitions["Ruleset"].AdditionalProperties = nil

	// output json
	j, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%s\n", j)
}
