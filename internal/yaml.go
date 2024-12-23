// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"

	bkYaml "github.com/buildkite/yaml"
	yaml "gopkg.in/yaml.v3"

	legacyTypes "github.com/go-vela/server/compiler/types/yaml/buildkite"
	types "github.com/go-vela/server/compiler/types/yaml/yaml"
)

// ParseYAML is a helper function for transitioning teams away from legacy buildkite YAML parser.
func ParseYAML(data []byte) (*types.Build, error) {
	var (
		rootNode yaml.Node
		version  string
	)

	err := yaml.Unmarshal(data, &rootNode)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal pipeline version yaml: %w", err)
	}

	if len(rootNode.Content) == 0 || rootNode.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("unable to find pipeline version in yaml")
	}

	for i, subNode := range rootNode.Content[0].Content {
		if subNode.Kind == yaml.ScalarNode && subNode.Value == "version" {
			if len(rootNode.Content[0].Content) > i {
				version = rootNode.Content[0].Content[i+1].Value

				break
			}
		}
	}

	config := new(types.Build)

	switch version {
	case "legacy":
		legacyConfig := new(legacyTypes.Build)

		err := bkYaml.Unmarshal(data, legacyConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal legacy yaml: %w", err)
		}

		config = legacyConfig.ToYAML()

	default:
		// unmarshal the bytes into the yaml configuration
		err := yaml.Unmarshal(data, config)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
		}
	}

	return config, nil
}
