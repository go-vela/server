// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"strings"

	bkYaml "github.com/buildkite/yaml"
	yaml "gopkg.in/yaml.v3"

	legacyTypes "github.com/go-vela/server/compiler/types/yaml/buildkite"
	types "github.com/go-vela/server/compiler/types/yaml/yaml"
)

// ParseYAML is a helper function for transitioning teams away from legacy buildkite YAML parser.
func ParseYAML(data []byte) (*types.Build, []string, error) {
	var (
		rootNode yaml.Node
		warnings []string
		version  string
	)

	err := yaml.Unmarshal(data, &rootNode)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal pipeline version yaml: %w", err)
	}

	if len(rootNode.Content) == 0 || rootNode.Content[0].Kind != yaml.MappingNode {
		return nil, nil, fmt.Errorf("unable to find pipeline version in yaml")
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
			return nil, nil, fmt.Errorf("unable to unmarshal legacy yaml: %w", err)
		}

		config = legacyConfig.ToYAML()

		warnings = append(warnings, `using legacy version - address any incompatibilities and use "1" instead`)

	default:
		// unmarshal the bytes into the yaml configuration
		err := yaml.Unmarshal(data, config)
		if err != nil {
			// if error is related to duplicate `<<` keys, attempt to fix
			if strings.Contains(err.Error(), "mapping key \"<<\" already defined") {
				root := new(yaml.Node)

				if err := yaml.Unmarshal(data, root); err != nil {
					fmt.Println("error unmarshalling YAML:", err)

					return nil, nil, err
				}

				warnings = collapseMergeAnchors(root.Content[0], warnings)

				newData, err := yaml.Marshal(root)
				if err != nil {
					return nil, nil, err
				}

				err = yaml.Unmarshal(newData, config)
				if err != nil {
					return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
				}
			} else {
				return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
			}
		}
	}

	return config, warnings, nil
}

// collapseMergeAnchors traverses the entire pipeline and replaces duplicate `<<` keys with a single key->sequence.
func collapseMergeAnchors(node *yaml.Node, warnings []string) []string {
	// only replace on maps
	if node.Kind == yaml.MappingNode {
		var (
			anchors      []*yaml.Node
			keysToRemove []int
			firstIndex   int
			firstFound   bool
		)

		// traverse mapping node content
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]

			// anchor found
			if keyNode.Value == "<<" {
				if (i+1) < len(node.Content) && node.Content[i+1].Kind == yaml.AliasNode {
					anchors = append(anchors, node.Content[i+1])
				}

				if !firstFound {
					firstIndex = i
					firstFound = true
				} else {
					keysToRemove = append(keysToRemove, i)
				}
			}
		}

		// only replace if there were duplicates
		if len(anchors) > 1 && firstFound {
			seqNode := &yaml.Node{
				Kind:    yaml.SequenceNode,
				Content: anchors,
			}

			node.Content[firstIndex] = &yaml.Node{Kind: yaml.ScalarNode, Value: "<<"}
			node.Content[firstIndex+1] = seqNode

			for i := len(keysToRemove) - 1; i >= 0; i-- {
				index := keysToRemove[i]

				warnings = append(warnings, fmt.Sprintf("%d:duplicate << keys in single YAML map", node.Content[index].Line))
				node.Content = append(node.Content[:index], node.Content[index+2:]...)
			}
		}

		// go to next level
		for _, content := range node.Content {
			warnings = collapseMergeAnchors(content, warnings)
		}
	} else if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			warnings = collapseMergeAnchors(item, warnings)
		}
	}

	return warnings
}
