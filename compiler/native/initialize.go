// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/yaml"
)

const (
	// default image for init process.
	initImage = "#init"
	// default name for init stage.
	initStageName = "init"
	// default name for init step.
	initStepName = "init"
)

// InitStage injects the init stage process into a yaml configuration.
func (c *client) InitStage(p *yaml.Build) (*yaml.Build, error) {
	stages := yaml.StageSlice{}

	// create new clone stage
	init := &yaml.Stage{
		Name: initStageName,
		Steps: yaml.StepSlice{
			&yaml.Step{
				Detach:     false,
				Image:      initImage,
				Name:       initStepName,
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
func (c *client) InitStep(p *yaml.Build) (*yaml.Build, error) {
	steps := yaml.StepSlice{}

	// create new init step
	init := &yaml.Step{
		Detach:     false,
		Image:      initImage,
		Name:       initStepName,
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
