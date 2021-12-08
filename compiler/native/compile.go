// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	yml "github.com/buildkite/yaml"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

// ModifyRequest contains the payload passed to the modification endpoint.
type ModifyRequest struct {
	Pipeline string `json:"pipeline,omitempty"`
	Build    int    `json:"build,omitempty"`
	Repo     string `json:"repo,omitempty"`
	Org      string `json:"org,omitempty"`
	User     string `json:"user,omitempty"`
}

// ModifyResponse contains the payload returned by the modification endpoint.
type ModifyResponse struct {
	Pipeline string `json:"pipeline,omitempty"`
}

// Compile produces an executable pipeline from a yaml configuration.
//
// nolint: gocyclo,funlen // ignore function length due to comments
func (c *client) Compile(v interface{}) (*pipeline.Build, error) {
	p, err := c.Parse(v)
	if err != nil {
		return nil, err
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, err
	}

	// create map of templates for easy lookup
	tmpls := mapFromTemplates(p.Templates)

	// create the ruledata to purge steps
	r := &pipeline.RuleData{
		Branch:  c.build.GetBranch(),
		Comment: c.comment,
		Event:   c.build.GetEvent(),
		Path:    c.files,
		Repo:    c.repo.GetFullName(),
		Tag:     strings.TrimPrefix(c.build.GetRef(), "refs/tags/"),
		Target:  c.build.GetDeploy(),
	}

	if len(p.Stages) > 0 {
		// check if the pipeline disabled the clone
		if p.Metadata.Clone == nil || *p.Metadata.Clone {
			// inject the clone stage
			p, err = c.CloneStage(p)
			if err != nil {
				return nil, err
			}
		}

		// inject the init stage
		p, err = c.InitStage(p)
		if err != nil {
			return nil, err
		}

		// inject the templates into the stages
		p.Stages, p.Secrets, p.Services, p.Environment, err = c.ExpandStages(p, tmpls)
		if err != nil {
			return nil, err
		}

		if c.ModificationService.Endpoint != "" {
			// send config to external endpoint for modification
			p, err = c.modifyConfig(p, c.build, c.repo)
			if err != nil {
				return nil, err
			}
		}

		// validate the yaml configuration
		err = c.Validate(p)
		if err != nil {
			return nil, err
		}

		// Create some default global environment inject vars
		// these are used below to overwrite to an empty
		// map if they should not be injected into a container
		envGlobalServices, envGlobalSecrets, envGlobalSteps := p.Environment, p.Environment, p.Environment

		if !p.Metadata.HasEnvironment("services") {
			envGlobalServices = make(raw.StringSliceMap)
		}

		if !p.Metadata.HasEnvironment("secrets") {
			envGlobalSecrets = make(raw.StringSliceMap)
		}

		if !p.Metadata.HasEnvironment("steps") {
			envGlobalSteps = make(raw.StringSliceMap)
		}

		// inject the environment variables into the services
		p.Services, err = c.EnvironmentServices(p.Services, envGlobalServices)
		if err != nil {
			return nil, err
		}

		// inject the environment variables into the secrets
		p.Secrets, err = c.EnvironmentSecrets(p.Secrets, envGlobalSecrets)
		if err != nil {
			return nil, err
		}

		// inject the environment variables into the stages
		p.Stages, err = c.EnvironmentStages(p.Stages, envGlobalSteps)
		if err != nil {
			return nil, err
		}

		// inject the substituted environment variables into the stages
		p.Stages, err = c.SubstituteStages(p.Stages)
		if err != nil {
			return nil, err
		}

		// inject the scripts into the stages
		p.Stages, err = c.ScriptStages(p.Stages)
		if err != nil {
			return nil, err
		}

		// return executable representation
		return c.TransformStages(r, p)
	}

	// check if the pipeline disabled the clone
	if p.Metadata.Clone == nil || *p.Metadata.Clone {
		// inject the clone step
		p, err = c.CloneStep(p)
		if err != nil {
			return nil, err
		}
	}

	// inject the init step
	p, err = c.InitStep(p)
	if err != nil {
		return nil, err
	}

	// inject the templates into the steps
	p.Steps, p.Secrets, p.Services, p.Environment, err = c.ExpandSteps(p, tmpls)
	if err != nil {
		return nil, err
	}

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		p, err = c.modifyConfig(p, c.build, c.repo)
		if err != nil {
			return nil, err
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, err
	}

	// Create some default global environment inject vars
	// these are used below to overwrite to an empty
	// map if they should not be injected into a container
	envGlobalServices, envGlobalSecrets, envGlobalSteps := p.Environment, p.Environment, p.Environment

	if !p.Metadata.HasEnvironment("services") {
		envGlobalServices = make(raw.StringSliceMap)
	}

	if !p.Metadata.HasEnvironment("secrets") {
		envGlobalSecrets = make(raw.StringSliceMap)
	}

	if !p.Metadata.HasEnvironment("steps") {
		envGlobalSteps = make(raw.StringSliceMap)
	}

	// inject the environment variables into the services
	p.Services, err = c.EnvironmentServices(p.Services, envGlobalServices)
	if err != nil {
		return nil, err
	}

	// inject the environment variables into the secrets
	p.Secrets, err = c.EnvironmentSecrets(p.Secrets, envGlobalSecrets)
	if err != nil {
		return nil, err
	}

	// inject the environment variables into the steps
	p.Steps, err = c.EnvironmentSteps(p.Steps, envGlobalSteps)
	if err != nil {
		return nil, err
	}

	// inject the substituted environment variables into the steps
	p.Steps, err = c.SubstituteSteps(p.Steps)
	if err != nil {
		return nil, err
	}

	// inject the scripts into the steps
	p.Steps, err = c.ScriptSteps(p.Steps)
	if err != nil {
		return nil, err
	}

	// return executable representation
	return c.TransformSteps(r, p)
}

// errorHandler ensures the error contains the number of request attempts.
func errorHandler(resp *http.Response, err error, attempts int) (*http.Response, error) {
	if err != nil {
		// nolint:lll // detailed error message
		err = fmt.Errorf("giving up connecting to modification endpoint after %d attempts due to: %v", attempts, err)
	}

	return resp, err
}

// modifyConfig sends the configuration to external http endpoint for modification.
// nolint:lll // parameter struct references push line limit
func (c *client) modifyConfig(build *yaml.Build, libraryBuild *library.Build, repo *library.Repo) (*yaml.Build, error) {
	// create request to send to endpoint
	data, err := yml.Marshal(build)
	if err != nil {
		return nil, err
	}

	modReq := &ModifyRequest{
		Pipeline: string(data),
		Build:    libraryBuild.GetNumber(),
		Repo:     repo.GetName(),
		Org:      repo.GetOrg(),
		User:     libraryBuild.GetAuthor(),
	}

	// marshal json to send in request
	b, err := json.Marshal(modReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modify payload")
	}

	// setup http client
	retryClient := retryablehttp.Client{
		HTTPClient:   cleanhttp.DefaultPooledClient(),
		RetryWaitMin: 500 * time.Millisecond,
		RetryWaitMax: 1 * time.Second,
		RetryMax:     c.ModificationService.Retries,
		CheckRetry:   retryablehttp.DefaultRetryPolicy,
		ErrorHandler: errorHandler,
		Backoff:      retryablehttp.DefaultBackoff,
	}

	// create POST request
	req, err := retryablehttp.NewRequest("POST", c.ModificationService.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	// ensure the overall request(s) do not take over the defined timeout
	ctx, cancel := context.WithTimeout(req.Request.Context(), c.ModificationService.Timeout)
	defer cancel()
	req.WithContext(ctx)

	// add content-type and auth headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ModificationService.Secret))

	// send the request
	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// fail if the response code was not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("modification endpoint returned status code %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	response := new(ModifyResponse)
	// unmarshal the response into the ModifyResponse struct
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON modification payload: %w", err)
	}

	newBuild := new(yaml.Build)
	// unmarshal the response into the yaml.Build struct
	err = yml.Unmarshal([]byte(response.Pipeline), &newBuild)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML modification payload: %w", err)
	}

	return newBuild, nil
}
