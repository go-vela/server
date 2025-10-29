// SPDX-License-Identifier: Apache-2.0

package native

import (
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
)

const (
	// default image for init process.
	initImage = "#init"
)

// InitStage injects the init stage process into a yaml configuration.
func (c *Client) InitStage(p *yaml.Build) (*yaml.Build, error) {
	stages := yaml.StageSlice{}

	// create new clone stage
	init := &yaml.Stage{
		Name: constants.InitName,
		Steps: yaml.StepSlice{
			&yaml.Step{
				Detach:     false,
				Image:      initImage,
				Name:       constants.InitName,
				Privileged: false,
				Pull:       constants.PullNotPresent,
			},
		},
	}

	// add init stage as first stage
	stages = append(stages, init)

	// add existing stages after init stage
	stages = append(stages, p.Stages...)

	// overwrite existing stages
	p.Stages = stages

	return p, nil
}

// InitStep injects the init step process into a yaml configuration.
func (c *Client) InitStep(p *yaml.Build) (*yaml.Build, error) {
	steps := yaml.StepSlice{}

	// create new init step
	init := &yaml.Step{
		Detach:     false,
		Image:      initImage,
		Name:       constants.InitName,
		Privileged: false,
		Pull:       constants.PullNotPresent,
	}

	// add init step as first step
	steps = append(steps, init)

	// add existing steps after init step
	steps = append(steps, p.Steps...)

	// overwrite existing steps
	p.Steps = steps

	return p, nil
}
