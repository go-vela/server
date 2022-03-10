// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
func (c *client) Compile(v interface{}) (*pipeline.Build, error) {
	p, err := c.Parse(v, c.repo.GetPipelineType(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, err
	}

	// create map of templates for easy lookup
	templates := mapFromTemplates(p.Templates)

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

	switch {
	case p.Metadata.RenderInline:
		newPipeline, err := c.compileInline(p)
		if err != nil {
			return nil, err
		}
		// validate the yaml configuration
		err = c.Validate(newPipeline)
		if err != nil {
			return nil, err
		}

		if len(newPipeline.Stages) > 0 {
			return c.compileStages(newPipeline, map[string]*yaml.Template{}, r)
		}

		return c.compileSteps(newPipeline, map[string]*yaml.Template{}, r)
	case len(p.Stages) > 0:
		return c.compileStages(p, templates, r)
	default:
		return c.compileSteps(p, templates, r)
	}
}

// CompileLite produces a partial of an executable pipeline from a yaml configuration.
func (c *client) CompileLite(v interface{}, template, substitute bool) (*yaml.Build, error) {
	p, err := c.Parse(v, c.repo.GetPipelineType(), map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	if p.Metadata.RenderInline {
		newPipeline, err := c.compileInline(p)
		if err != nil {
			return nil, err
		}
		// validate the yaml configuration
		err = c.Validate(newPipeline)
		if err != nil {
			return nil, err
		}

		p = newPipeline
	}

	if template {
		// create map of templates for easy lookup
		templates := mapFromTemplates(p.Templates)

		switch {
		case len(p.Stages) > 0:
			// inject the templates into the steps
			p, err = c.ExpandStages(p, templates)
			if err != nil {
				return nil, err
			}

			if substitute {
				// inject the substituted environment variables into the steps
				p.Stages, err = c.SubstituteStages(p.Stages)
				if err != nil {
					return nil, err
				}
			}
		case len(p.Steps) > 0:
			// inject the templates into the steps
			p, err = c.ExpandSteps(p, templates)
			if err != nil {
				return nil, err
			}

			if substitute {
				// inject the substituted environment variables into the steps
				p.Steps, err = c.SubstituteSteps(p.Steps)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// compileInline parses and expands out inline pipelines.
func (c *client) compileInline(p *yaml.Build) (*yaml.Build, error) {
	newPipeline := *p
	newPipeline.Templates = yaml.TemplateSlice{}

	for _, template := range p.Templates {
		bytes, err := c.getTemplate(template, template.Name)
		if err != nil {
			return nil, err
		}

		parsed, err := c.Parse(bytes, template.Format, template.Variables)
		if err != nil {
			return nil, err
		}

		switch {
		case len(parsed.Environment) > 0:
			for key, value := range parsed.Environment {
				newPipeline.Environment[key] = value
			}

			fallthrough
		case len(parsed.Stages) > 0:
			// ensure all templated steps inside stages have template prefix
			for stgIndex, newStage := range parsed.Stages {
				parsed.Stages[stgIndex].Name = fmt.Sprintf("%s_%s", template.Name, newStage.Name)

				for index, newStep := range newStage.Steps {
					parsed.Stages[stgIndex].Steps[index].Name = fmt.Sprintf("%s_%s", template.Name, newStep.Name)
				}
			}

			newPipeline.Stages = append(newPipeline.Stages, parsed.Stages...)

			fallthrough
		case len(parsed.Steps) > 0:
			// ensure all templated steps have template prefix
			for index, newStep := range parsed.Steps {
				parsed.Steps[index].Name = fmt.Sprintf("%s_%s", template.Name, newStep.Name)
			}

			newPipeline.Steps = append(newPipeline.Steps, parsed.Steps...)

			fallthrough
		case len(parsed.Services) > 0:
			newPipeline.Services = append(newPipeline.Services, parsed.Services...)
			fallthrough
		case len(parsed.Secrets) > 0:
			newPipeline.Secrets = append(newPipeline.Secrets, parsed.Secrets...)
		default:
			// nolint: lll // ignore long line length due to error message
			return nil, fmt.Errorf("empty template %s provided: template must contain secrets, services, stages or steps", template.Name)
		}

		if len(newPipeline.Stages) > 0 && len(newPipeline.Steps) > 0 {
			// nolint: lll // ignore long line length due to error message
			return nil, fmt.Errorf("invalid template %s provided: templates cannot mix stages and steps", template.Name)
		}
	}

	// validate the yaml configuration
	err := c.Validate(&newPipeline)
	if err != nil {
		return nil, err
	}

	return &newPipeline, nil
}

// compileSteps executes the workflow for converting a YAML pipeline into an executable struct.
//
// nolint:dupl,lll // linter thinks the steps and stages workflows are identical
func (c *client) compileSteps(p *yaml.Build, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, error) {
	var err error

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
	p, err = c.ExpandSteps(p, tmpls)
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

// compileStages executes the workflow for converting a YAML pipeline into an executable struct.
//
// nolint:dupl,lll // linter thinks the steps and stages workflows are identical
func (c *client) compileStages(p *yaml.Build, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, error) {
	var err error

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
	p, err = c.ExpandStages(p, tmpls)
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

// errorHandler ensures the error contains the number of request attempts.
func errorHandler(resp *http.Response, err error, attempts int) (*http.Response, error) {
	if err != nil {
		err = fmt.Errorf("giving up connecting to modification endpoint after %d attempts due to: %w", attempts, err)
	}

	return resp, err
}

// modifyConfig sends the configuration to external http endpoint for modification.
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
