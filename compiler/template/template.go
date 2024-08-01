// SPDX-License-Identifier: Apache-2.0

package template

import "github.com/go-vela/types/yaml"

// Engine represents the interface for Vela integrating
// with the different supported template engines.
type Engine interface {
	// RenderBuild defines a function that combines
	// the template with the build.
	RenderBuild(template string, step *yaml.Step) (yaml.StepSlice, error)
	// Render defines a function that combines
	// the template with the step.
	Render(template string, step *yaml.Step) (yaml.StepSlice, error)
}
