// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"strings"

	"github.com/buildkite/yaml"

	"github.com/go-vela/types/raw"
)

// convertPlatformVars takes the platform injected variables
// within the step environment block and modifies them to be returned
// within the template.
func convertPlatformVars(slice raw.StringSliceMap, name string) raw.StringSliceMap {
	envs := make(map[string]string)

	// iterate through the list of key/value pairs provided
	for key, value := range slice {
		// lowercase the key
		key = strings.ToLower(key)

		// iterate through the list of possible prefixes to look for
		for _, prefix := range []string{"deployment_parameter_", "vela_"} {
			// check if the key has the prefix
			if strings.HasPrefix(key, prefix) {
				// add the non-prefixed key/value pair
				envs[strings.TrimPrefix(key, prefix)] = value
			}
		}
	}

	envs["template_name"] = name

	return envs
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This code is under copyright (full attribution in NOTICE) and is from:

// https://github.com/helm/helm/blob/a499b4b179307c267bdf3ec49b880e3dbd2a5591/pkg/engine/funcs.go#L83
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}

	return strings.TrimSuffix(string(data), "\n")
}

type funcHandler struct {
	envs raw.StringSliceMap
}

// returnPlatformVar returns the value of the platform
// variable if it exists within the environment map.
func (h funcHandler) returnPlatformVar(input string) string {
	// lowercase the input keyx
	input = strings.ToLower(input)

	// iterate through the list of possible prefixes to look for
	for _, prefix := range []string{"deployment_parameter_", "vela_"} {
		// trim the prefix from the input key
		trimmed := strings.TrimPrefix(input, prefix)
		// check if the key exists within map
		if _, ok := h.envs[trimmed]; ok {
			// return the non-prefixed value if exists
			return h.envs[trimmed]
		}
	}

	// return empty string if not exists
	return ""
}
