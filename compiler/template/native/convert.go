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
	for key, value := range slice {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "vela_") {
			envs[strings.TrimPrefix(key, "vela_")] = value
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
	input = strings.ToLower(input)
	input = strings.TrimPrefix(input, "vela_")
	// check if key exists within map
	if _, ok := h.envs[input]; ok {
		// return value if exists
		return h.envs[input]
	}
	// return empty string if not exists
	return ""
}
