// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/yaml"
)

const (
	// default name for clone stage.
	cloneStageName = "clone"
	// default name for clone step.
	cloneStepName = "clone"
)

// CloneStage injects the clone stage process into a yaml configuration.
func (c *client) CloneStage(p *yaml.Build) (*yaml.Build, error) {
	// check if the compiler is setup for a local pipeline
	if c.local {
		// skip injecting the clone process
		return p, nil
	}

	c.log.AppendData([]byte("injecting clone stage\n"))

	stages := yaml.StageSlice{}

	// create new clone stage
	clone := &yaml.Stage{
		Name: cloneStageName,
		Steps: yaml.StepSlice{
			&yaml.Step{
				Detach:     false,
				Image:      c.CloneImage,
				Name:       cloneStepName,
				Privileged: false,
				Pull:       constants.PullNotPresent,
			},
		},
	}

	// add clone stage as first stage
	stages = append(stages, clone)

	// add existing stages after clone stage
	stages = append(stages, p.Stages...)

	// overwrite existing stages
	p.Stages = stages

	return p, nil
}

// CloneStep injects the clone step process into a yaml configuration.
func (c *client) CloneStep(p *yaml.Build) (*yaml.Build, error) {
	// check if the compiler is setup for a local pipeline
	if c.local {
		// skip injecting the clone process
		return p, nil
	}

	c.log.AppendData([]byte("injecting clone step\n"))

	steps := yaml.StepSlice{}

	// create new clone step
	clone := &yaml.Step{
		Detach:     false,
		Image:      c.CloneImage,
		Name:       cloneStepName,
		Privileged: false,
		Pull:       constants.PullNotPresent,
	}

	// add clone step as first step
	steps = append(steps, clone)

	// add existing steps after clone step
	steps = append(steps, p.Steps...)

	// overwrite existing steps
	p.Steps = steps

	return p, nil
}
