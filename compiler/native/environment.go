// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"
	"os"
	"strings"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal"
)

// EnvironmentStages injects environment variables
// for each stage in a yaml configuration.
func (c *Client) EnvironmentStages(s yaml.StageSlice, globalEnv raw.StringSliceMap) (yaml.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		_, err := c.EnvironmentStage(stage, globalEnv)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// EnvironmentStage injects environment variables
// for each stage in a yaml configuration.
func (c *Client) EnvironmentStage(s *yaml.Stage, globalEnv raw.StringSliceMap) (*yaml.Stage, error) {
	// make empty map of environment variables
	env := make(map[string]string)

	// gather set of default environment variables
	defaultEnv := environment(c.build, c.metadata, c.repo, c.user, c.netrc)

	// inject the declared global environment
	// WARNING: local env can override global
	env = appendMap(env, globalEnv)

	// inject the declared environment
	// variables to the build stage
	for k, v := range s.Environment {
		env[k] = v
	}

	// inject the default environment
	// variables to the build stage
	// we do this after injecting the declared environment
	// to ensure the default env overrides any conflicts
	for k, v := range defaultEnv {
		env[k] = v
	}

	// overwrite existing build stage environment
	s.Environment = env

	// inject the environment variables into the steps for the stage
	steps, err := c.EnvironmentSteps(s.Steps, env)
	if err != nil {
		return nil, err
	}

	s.Steps = steps

	return s, nil
}

// EnvironmentSteps injects environment variables
// for each step in a stage for the yaml configuration.
func (c *Client) EnvironmentSteps(s yaml.StepSlice, stageEnv raw.StringSliceMap) (yaml.StepSlice, error) {
	// iterate through all steps
	for _, step := range s {
		_, err := c.EnvironmentStep(step, stageEnv)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// EnvironmentStep injects environment variables
// a single step in a yaml configuration.
func (c *Client) EnvironmentStep(s *yaml.Step, stageEnv raw.StringSliceMap) (*yaml.Step, error) {
	// make empty map of environment variables
	env := make(map[string]string)

	// gather set of default environment variables
	defaultEnv := environment(c.build, c.metadata, c.repo, c.user, c.netrc)

	// inject the declared stage environment
	// WARNING: local env can override global + stage
	env = appendMap(env, stageEnv)

	// inject the declared environment
	// variables to the build step
	for k, v := range s.Environment {
		env[k] = v
	}

	// inject the default environment
	// variables to the build step
	// we do this after injecting the declared environment
	// to ensure the default env overrides any conflicts
	for k, v := range defaultEnv {
		env[k] = v
	}

	// check if the compiler is setup for a local pipeline
	if c.local && !s.Detach {
		// capture all environment variables from the local environment
		for _, e := range os.Environ() {
			// split the environment variable on = into a key value pair
			parts := strings.SplitN(e, "=", 2)

			env[parts[0]] = parts[1]
		}
	}

	// inject the declared parameter
	// variables to the build step
	for k, v := range s.Parameters {
		if v == nil {
			continue
		}

		// parameter keys are passed to the image
		// as PARAMETER_ environment variables
		k = "PARAMETER_" + strings.ToUpper(k)

		// parameter values are passed to the image
		// as string environment variables
		env[k] = api.ToString(v)
	}

	// overwrite existing build step environment
	s.Environment = env

	return s, nil
}

// EnvironmentServices injects environment variables
// for each service in a yaml configuration.
func (c *Client) EnvironmentServices(s yaml.ServiceSlice, globalEnv raw.StringSliceMap) (yaml.ServiceSlice, error) {
	// iterate through all services
	for _, service := range s {
		// make empty map of environment variables
		env := make(map[string]string)

		// gather set of default environment variables
		defaultEnv := environment(c.build, c.metadata, c.repo, c.user, c.netrc)

		// inject the declared global environment
		// WARNING: local env can override global
		env = appendMap(env, globalEnv)

		// inject the declared environment
		// variables to the build service
		for k, v := range service.Environment {
			env[k] = v
		}

		// inject the default environment
		// variables to the build service
		// we do this after injecting the declared environment
		// to ensure the default env overrides any conflicts
		for k, v := range defaultEnv {
			env[k] = v
		}

		// overwrite existing build service environment
		service.Environment = env
	}

	return s, nil
}

// EnvironmentSecrets injects environment variables
// for each secret plugin in a yaml configuration.
func (c *Client) EnvironmentSecrets(s yaml.SecretSlice, globalEnv raw.StringSliceMap) (yaml.SecretSlice, error) {
	// iterate through all secrets
	for _, secret := range s {
		// skip non plugin secrets
		if secret.Origin.Empty() {
			continue
		}

		// make empty map of environment variables
		env := make(map[string]string)

		// gather set of default environment variables
		defaultEnv := environment(c.build, c.metadata, c.repo, c.user, c.netrc)

		// inject the declared global environment
		// WARNING: local env can override global
		env = appendMap(env, globalEnv)

		// inject the declared environment
		// variables to the build secret
		for k, v := range secret.Origin.Environment {
			env[k] = v
		}

		// inject the default environment
		// variables to the build secret
		// we do this after injecting the declared environment
		// to ensure the default env overrides any conflicts
		for k, v := range defaultEnv {
			env[k] = v
		}

		// check if the compiler is setup for a local pipeline
		if c.local {
			// capture all environment variables from the local environment
			for _, e := range os.Environ() {
				// split the environment variable on = into a key value pair
				parts := strings.SplitN(e, "=", 2)

				env[parts[0]] = parts[1]
			}
		}

		// inject the declared parameter
		// variables to the build secret
		for k, v := range secret.Origin.Parameters {
			if v == nil {
				continue
			}

			// parameter keys are passed to the image
			// as PARAMETER_ environment variables
			k = "PARAMETER_" + strings.ToUpper(k)

			// parameter values are passed to the image
			// as string environment variables
			env[k] = api.ToString(v)
		}

		// overwrite existing build secret environment
		secret.Origin.Environment = env
	}

	return s, nil
}

// EnvironmentBuild injects environment variables
// for the build in a yaml configuration.
func (c *Client) EnvironmentBuild() map[string]string {
	// make empty map of environment variables
	env := make(map[string]string)

	// gather set of default environment variables
	defaultEnv := environment(c.build, c.metadata, c.repo, c.user, c.netrc)

	// inject the default environment
	// variables to the build
	// we do this after injecting the declared environment
	// to ensure the default env overrides any conflicts
	for k, v := range defaultEnv {
		env[k] = v
	}

	// check if the compiler is setup for a local pipeline
	if c.local {
		// capture all environment variables from the local environment
		for _, e := range os.Environ() {
			// split the environment variable on = into a key value pair
			parts := strings.SplitN(e, "=", 2)

			env[parts[0]] = parts[1]
		}
	}

	return env
}

// helper function to merge two maps together.
func appendMap(originalMap, otherMap map[string]string) map[string]string {
	for key, value := range otherMap {
		originalMap[key] = value
	}

	return originalMap
}

// helper function that creates the standard set of environment variables for a pipeline.
func environment(b *api.Build, m *internal.Metadata, r *api.Repo, u *api.User, netrc *string) map[string]string {
	// set default workspace
	workspace := constants.WorkspaceDefault
	notImplemented := "TODO"

	env := make(map[string]string)

	// vela specific environment variables
	env["VELA"] = api.ToString(true)
	env["VELA_ADDR"] = notImplemented
	env["VELA_DATABASE"] = notImplemented
	env["VELA_DISTRIBUTION"] = notImplemented
	env["VELA_HOST"] = notImplemented
	env["VELA_NETRC_MACHINE"] = notImplemented
	env["VELA_NETRC_PASSWORD"] = notImplemented
	env["VELA_NETRC_USERNAME"] = "x-oauth-basic"
	env["VELA_QUEUE"] = notImplemented
	env["VELA_RUNTIME"] = notImplemented
	env["VELA_SOURCE"] = notImplemented
	env["VELA_VERSION"] = notImplemented
	env["CI"] = "true"

	// populate environment variables from metadata
	if m != nil {
		env["VELA_ADDR"] = m.Vela.WebAddress
		env["VELA_SERVER_ADDR"] = m.Vela.Address
		env["VELA_DATABASE"] = m.Database.Driver
		env["VELA_HOST"] = m.Vela.Address
		env["VELA_NETRC_MACHINE"] = m.Source.Host
		env["VELA_QUEUE"] = m.Queue.Driver
		env["VELA_SOURCE"] = m.Source.Driver
		env["VELA_OPEN_ID_ISSUER"] = m.Vela.OpenIDIssuer
		env["VELA_ID_TOKEN_REQUEST_URL"] = fmt.Sprintf("%s/api/v1/repos/%s/builds/%d/id_token", m.Vela.Address, r.GetFullName(), b.GetNumber())
		workspace = fmt.Sprintf("%s/%s/%s/%s", workspace, m.Source.Host, r.GetOrg(), r.GetName())
	}

	if netrc != nil {
		env["VELA_NETRC_PASSWORD"] = *netrc
	}

	env["VELA_WORKSPACE"] = workspace

	// populate environment variables from repo api
	env = appendMap(env, r.Environment())
	// populate environment variables from build api
	env = appendMap(env, b.Environment(workspace))
	// populate environment variables from user api
	env = appendMap(env, u.Environment())

	return env
}
