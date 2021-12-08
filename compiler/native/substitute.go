// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/buildkite/yaml"

	"github.com/drone/envsubst"

	types "github.com/go-vela/types/yaml"
)

// SubstituteStages replaces every declared environment
// variable with it's corresponding value for each step
// in every stage in a yaml configuration.
func (c *client) SubstituteStages(s types.StageSlice) (types.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		// inject the scripts into the steps for the stage
		steps, err := c.SubstituteSteps(stage.Steps)
		if err != nil {
			return nil, err
		}

		stage.Steps = steps
	}

	return s, nil
}

// SubstituteSteps replaces every declared environment
// variable with it's corresponding value for each step
// in a yaml configuration.
func (c *client) SubstituteSteps(s types.StepSlice) (types.StepSlice, error) {
	// iterate through all steps
	for _, step := range s {
		// marshal step configuration
		body, err := yaml.Marshal(step)
		if err != nil {
			return nil, fmt.Errorf("unable to marshal configuration: %v", err)
		}

		// create substitute function
		subFunc := func(name string) string {
			// check for the environment variable
			env, ok := step.Environment[name]
			if !ok {
				// return the original declaration if
				// the environment variable isn't found
				return fmt.Sprintf("${%s}", name)
			}

			// check for a new line
			if strings.Contains(env, "\n") {
				// escape the environment variable
				env = fmt.Sprintf("%q", env)
			}

			return env
		}

		// substitute the environment variables
		subStep, err := envsubst.Eval(string(body), subFunc)
		if err != nil {
			return nil, fmt.Errorf("unable to substitute environment variables: %v", err)
		}

		// unmarshal step configuration
		err = yaml.Unmarshal([]byte(subStep), step)
		if err != nil {
			return nil, fmt.Errorf("unable to unmarshal configuration: %v", err)
		}
	}

	return s, nil
}
