// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-vela/types/constants"

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
func (c *client) Compile(v interface{}) (*pipeline.Build, *library.Pipeline, error) {
	p, data, err := c.Parse(v, c.repo.GetPipelineType(), new(yaml.Template))
	if err != nil {
		c.log.AppendData([]byte("pipeline parse failed\n"))
		return nil, nil, err
	}

	// create the library pipeline object from the yaml configuration
	_pipeline := p.ToPipelineLibrary()
	_pipeline.SetData(data)
	_pipeline.SetType(c.repo.GetPipelineType())

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, _pipeline, err
	}

	// create map of templates for easy lookup
	templates := mapFromTemplates(p.Templates)

	event := c.build.GetEvent()
	action := c.build.GetEventAction()

	// if the build has an event action, concatenate event and event action for matching
	if !strings.EqualFold(action, "") {
		event = event + ":" + action
	}

	// create the ruledata to purge steps
	r := &pipeline.RuleData{
		Branch:  c.build.GetBranch(),
		Comment: c.comment,
		Event:   event,
		Path:    c.files,
		Repo:    c.repo.GetFullName(),
		Tag:     strings.TrimPrefix(c.build.GetRef(), "refs/tags/"),
		Target:  c.build.GetDeploy(),
	}

	switch {
	case p.Metadata.RenderInline:
		newPipeline, err := c.compileInline(p, nil)
		if err != nil {
			return nil, _pipeline, err
		}
		// validate the yaml configuration
		err = c.Validate(newPipeline)
		if err != nil {
			return nil, _pipeline, err
		}

		if len(newPipeline.Stages) > 0 {
			return c.compileStages(newPipeline, _pipeline, map[string]*yaml.Template{}, r)
		}

		return c.compileSteps(newPipeline, _pipeline, map[string]*yaml.Template{}, r)
	case len(p.Stages) > 0:
		return c.compileStages(p, _pipeline, templates, r)
	default:
		return c.compileSteps(p, _pipeline, templates, r)
	}
}

// CompileLite produces a partial of an executable pipeline from a yaml configuration.
func (c *client) CompileLite(v interface{}, template, substitute bool, localTemplates []string) (*yaml.Build, *library.Pipeline, error) {
	p, data, err := c.Parse(v, c.repo.GetPipelineType(), new(yaml.Template))
	if err != nil {
		return nil, nil, err
	}

	// create the library pipeline object from the yaml configuration
	_pipeline := p.ToPipelineLibrary()
	_pipeline.SetData(data)
	_pipeline.SetType(c.repo.GetPipelineType())

	if p.Metadata.RenderInline {
		newPipeline, err := c.compileInline(p, localTemplates)
		if err != nil {
			return nil, _pipeline, err
		}
		// validate the yaml configuration
		err = c.Validate(newPipeline)
		if err != nil {
			return nil, _pipeline, err
		}

		p = newPipeline
	}

	if template {
		// create map of templates for easy lookup
		templates := mapFromTemplates(p.Templates)

		if c.local {
			for _, file := range localTemplates {
				// local templates override format is <name>:<source>
				//
				// example: example:/path/to/template.yml
				parts := strings.Split(file, ":")

				// make sure the template was configured
				_, ok := templates[parts[0]]
				if !ok {
					return nil, _pipeline, fmt.Errorf("template with name %s is not configured", parts[0])
				}

				// override the source for the given template
				templates[parts[0]].Source = parts[1]
			}
		}

		switch {
		case len(p.Stages) > 0:
			// inject the templates into the steps
			p, err = c.ExpandStages(p, templates)
			if err != nil {
				return nil, _pipeline, err
			}

			if substitute {
				// inject the substituted environment variables into the steps
				p.Stages, err = c.SubstituteStages(p.Stages)
				if err != nil {
					return nil, _pipeline, err
				}
			}
		case len(p.Steps) > 0:
			// inject the templates into the steps
			p, err = c.ExpandSteps(p, templates)
			if err != nil {
				return nil, _pipeline, err
			}

			if substitute {
				// inject the substituted environment variables into the steps
				p.Steps, err = c.SubstituteSteps(p.Steps)
				if err != nil {
					return nil, _pipeline, err
				}
			}
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, _pipeline, err
	}

	return p, _pipeline, nil
}

// compileInline parses and expands out inline pipelines.
func (c *client) compileInline(p *yaml.Build, localTemplates []string) (*yaml.Build, error) {
	c.log.AppendData([]byte("rendering inline pipeline template...\n"))

	newPipeline := *p
	newPipeline.Templates = yaml.TemplateSlice{}

	for _, template := range p.Templates {
		if c.local {
			for _, file := range localTemplates {
				// local templates override format is <name>:<source>
				//
				// example: example:/path/to/template.yml
				parts := strings.Split(file, ":")

				// make sure we're referencing the proper template
				if parts[0] == template.Name {
					// override the source for the given template
					template.Source = parts[1]
				}
			}
		}

		bytes, err := c.getTemplate(template, template.Name)
		if err != nil {
			return nil, err
		}

		format := template.Format

		// set the default format to golang if the user did not define anything
		if template.Format == "" {
			format = constants.PipelineTypeGo
		}

		parsed, _, err := c.Parse(bytes, format, template)
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
			//nolint:lll // ignore long line length due to error message
			return nil, fmt.Errorf("empty template %s provided: template must contain secrets, services, stages or steps", template.Name)
		}

		if len(newPipeline.Stages) > 0 && len(newPipeline.Steps) > 0 {
			//nolint:lll // ignore long line length due to error message
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
//nolint:dupl,lll // linter thinks the steps and stages workflows are identical
func (c *client) compileSteps(p *yaml.Build, _pipeline *library.Pipeline, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, *library.Pipeline, error) {
	c.log.AppendData([]byte("compiling steps...\n"))

	var err error

	// check if the pipeline disabled the clone
	if p.Metadata.Clone == nil || *p.Metadata.Clone {
		// inject the clone step
		p, err = c.CloneStep(p)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// inject the init step // TODO: stop injecting the init step
	p, err = c.InitStep(p)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the templates into the steps
	p, err = c.ExpandSteps(p, tmpls)
	if err != nil {
		return nil, _pipeline, err
	}

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		p, err = c.modifyConfig(p, c.build, c.repo)
		if err != nil {
			c.log.AppendData([]byte("modification endpoint error\n"))
			return nil, _pipeline, err
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, _pipeline, err
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
		return nil, _pipeline, err
	}

	// inject the environment variables into the secrets
	p.Secrets, err = c.EnvironmentSecrets(p.Secrets, envGlobalSecrets)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the environment variables into the steps
	p.Steps, err = c.EnvironmentSteps(p.Steps, envGlobalSteps)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the substituted environment variables into the steps
	p.Steps, err = c.SubstituteSteps(p.Steps)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the scripts into the steps
	p.Steps, err = c.ScriptSteps(p.Steps)
	if err != nil {
		return nil, _pipeline, err
	}

	// create executable representation
	build, err := c.TransformSteps(r, p)
	if err != nil {
		return nil, _pipeline, err
	}

	return build, _pipeline, nil
}

// compileStages executes the workflow for converting a YAML pipeline into an executable struct.
//
//nolint:dupl,lll // linter thinks the steps and stages workflows are identical
func (c *client) compileStages(p *yaml.Build, _pipeline *library.Pipeline, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, *library.Pipeline, error) {
	c.log.AppendData([]byte("compiling stages...\n"))

	var err error

	// check if the pipeline disabled the clone
	if p.Metadata.Clone == nil || *p.Metadata.Clone {
		// inject the clone stage
		p, err = c.CloneStage(p)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// inject the init stage // TODO: stop injecting the init stage
	p, err = c.InitStage(p)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the templates into the stages
	p, err = c.ExpandStages(p, tmpls)
	if err != nil {
		return nil, _pipeline, err
	}

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		p, err = c.modifyConfig(p, c.build, c.repo)
		if err != nil {
			c.log.AppendData([]byte("modification endpoint error\n"))
			return nil, _pipeline, err
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, _pipeline, err
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
		return nil, _pipeline, err
	}

	// inject the environment variables into the secrets
	p.Secrets, err = c.EnvironmentSecrets(p.Secrets, envGlobalSecrets)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the environment variables into the stages
	p.Stages, err = c.EnvironmentStages(p.Stages, envGlobalSteps)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the substituted environment variables into the stages
	p.Stages, err = c.SubstituteStages(p.Stages)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the scripts into the stages
	p.Stages, err = c.ScriptStages(p.Stages)
	if err != nil {
		return nil, _pipeline, err
	}

	// create executable representation
	build, err := c.TransformStages(r, p)
	if err != nil {
		return nil, _pipeline, err
	}

	return build, _pipeline, nil
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
	c.log.AppendData([]byte("sending pipeline to modification endpoint\n"))
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

	// ensure the overall request(s) do not take over the defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.ModificationService.Timeout)
	defer cancel()

	// create POST request
	req, err := retryablehttp.NewRequestWithContext(ctx, "POST", c.ModificationService.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

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

	body, err := io.ReadAll(resp.Body)
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
