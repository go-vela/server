// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package template

import "github.com/go-vela/types/yaml"

// Engine represents the interface for Vela integrating
// with the different supported template engines.
type Engine interface {
	// RenderBuild defines a function that combines
	// the template with the build.
	RenderBuild(template string, step *yaml.Step) (yaml.StepSlice, error)
	// RenderStep defines a function that combines
	// the template with the step.
	RenderStep(template string, step *yaml.Step) (yaml.StepSlice, error)
}
