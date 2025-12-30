// SPDX-License-Identifier: Apache-2.0

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

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	yml "go.yaml.in/yaml/v3"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
)

// ModifyRequest contains the payload passed to the modification endpoint.
type ModifyRequest struct {
	Pipeline string `json:"pipeline,omitempty"`
	Build    int64  `json:"build,omitempty"`
	Repo     string `json:"repo,omitempty"`
	Org      string `json:"org,omitempty"`
	User     string `json:"user,omitempty"`
}

// ModifyResponse contains the payload returned by the modification endpoint.
type ModifyResponse struct {
	Pipeline string `json:"pipeline,omitempty"`
}

// Compile produces an executable pipeline from a yaml configuration.
func (c *Client) Compile(ctx context.Context, v interface{}) (*pipeline.Build, *api.Pipeline, error) {
	p, data, warnings, err := c.Parse(v, c.repo.GetPipelineType(), new(yaml.Template))
	if err != nil {
		return nil, nil, err
	}

	// create the netrc using the scm
	// this has to occur after Parse because the scm configurations might be set in yaml
	// netrc can be provided directly using WithNetrc for situations like local exec
	if c.netrc == nil && c.scm != nil {
		// get the netrc password from the scm
		netrc, exp, err := c.scm.GetNetrcPassword(ctx, c.db, c.cache, c.repo, c.user, p.Git)
		if err != nil {
			return nil, nil, err
		}

		c.WithNetrc(netrc)
		c.netrcExp = exp
	}

	// create the API pipeline object from the yaml configuration
	_pipeline := p.ToPipelineAPI()
	_pipeline.SetData(data)
	_pipeline.SetType(c.repo.GetPipelineType())
	_pipeline.SetWarnings(warnings)

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
		Sender:  c.build.GetSender(),
		Tag:     strings.TrimPrefix(c.build.GetRef(), "refs/tags/"),
		Target:  c.build.GetDeploy(),
		Label:   c.labels,
		Status:  c.build.GetStatus(),
		Env:     make(raw.StringSliceMap),
	}

	// add instance when we have the metadata (local exec will not)
	if c.metadata != nil && c.metadata.Vela != nil {
		r.Instance = c.metadata.Vela.Address
	}

	switch {
	case p.Metadata.RenderInline:
		newPipeline, err := c.compileInline(ctx, p, c.GetTemplateDepth())
		if err != nil {
			return nil, _pipeline, err
		}
		// validate the yaml configuration
		err = c.ValidateYAML(newPipeline)
		if err != nil {
			return nil, _pipeline, err
		}

		if len(newPipeline.Stages) > 0 {
			return c.compileStages(ctx, newPipeline, _pipeline, map[string]*yaml.Template{}, r)
		}

		return c.compileSteps(ctx, newPipeline, _pipeline, map[string]*yaml.Template{}, r)
	case len(p.Stages) > 0:
		return c.compileStages(ctx, p, _pipeline, templates, r)
	default:
		return c.compileSteps(ctx, p, _pipeline, templates, r)
	}
}

// CompileLite produces a partial of an executable pipeline from a yaml configuration.
func (c *Client) CompileLite(ctx context.Context, v interface{}, ruleData *pipeline.RuleData, substitute bool) (*yaml.Build, *api.Pipeline, error) {
	p, data, warnings, err := c.Parse(v, c.repo.GetPipelineType(), new(yaml.Template))
	if err != nil {
		return nil, nil, err
	}

	// create the API pipeline object from the yaml configuration
	_pipeline := p.ToPipelineAPI()
	_pipeline.SetData(data)
	_pipeline.SetType(c.repo.GetPipelineType())
	_pipeline.SetWarnings(warnings)

	if p.Metadata.RenderInline {
		newPipeline, err := c.compileInline(ctx, p, c.GetTemplateDepth())
		if err != nil {
			return nil, _pipeline, err
		}
		// validate the yaml configuration
		err = c.ValidateYAML(newPipeline)
		if err != nil {
			return nil, _pipeline, err
		}

		p = newPipeline
	}

	// create map of templates for easy lookup
	templates := mapFromTemplates(p.Templates)

	// expand deployment config
	p, err = c.ExpandDeployment(ctx, p, templates)
	if err != nil {
		return nil, _pipeline, err
	}

	switch {
	case len(p.Stages) > 0:
		// inject the templates into the steps
		p, warnings, err = c.ExpandStages(ctx, p, templates, ruleData, _pipeline.GetWarnings())
		if err != nil {
			return nil, _pipeline, err
		}

		_pipeline.SetWarnings(warnings)

		if substitute {
			// inject the substituted environment variables into the steps
			p.Stages, err = c.SubstituteStages(p.Stages)
			if err != nil {
				return nil, _pipeline, err
			}
		}

		if ruleData != nil {
			purgedStages := yaml.StageSlice{}

			for _, stg := range p.Stages {
				stg.Steps = purgeStepsLite(stg.Steps, ruleData)

				if len(stg.Steps) > 0 {
					purgedStages = append(purgedStages, stg)
				}
			}

			p.Secrets = purgeSecretsLite(p.Secrets, ruleData)
			p.Services = purgeServicesLite(p.Services, ruleData)
			p.Stages = purgedStages
		}

	case len(p.Steps) > 0:
		// inject the templates into the steps
		p, warnings, err = c.ExpandSteps(ctx, p, templates, ruleData, _pipeline.GetWarnings(), c.GetTemplateDepth())
		if err != nil {
			return nil, _pipeline, err
		}

		_pipeline.SetWarnings(warnings)

		if substitute {
			// inject the substituted environment variables into the steps
			p.Steps, err = c.SubstituteSteps(p.Steps)
			if err != nil {
				return nil, _pipeline, err
			}
		}

		if ruleData != nil {
			p.Secrets = purgeSecretsLite(p.Secrets, ruleData)
			p.Services = purgeServicesLite(p.Services, ruleData)
			p.Steps = purgeStepsLite(p.Steps, ruleData)
		}
	}

	// validate the yaml configuration
	err = c.ValidateYAML(p)
	if err != nil {
		return nil, _pipeline, err
	}

	return p, _pipeline, nil
}

// compileInline parses and expands out inline pipelines.
func (c *Client) compileInline(ctx context.Context, p *yaml.Build, depth int) (*yaml.Build, error) {
	newPipeline := *p

	// return if max template depth has been reached
	if depth == 0 {
		retErr := fmt.Errorf("max template depth of %d exceeded", c.GetTemplateDepth())

		return nil, retErr
	}

	for _, template := range p.Templates {
		var (
			bytes []byte
			found bool
			err   error
		)

		if bytes, found = c.TemplateCache[template.Source]; !found {
			bytes, err = c.getTemplate(ctx, template, template.Name)
			if err != nil {
				return nil, err
			}
		}

		format := template.Format

		// set the default format to golang if the user did not define anything
		if template.Format == "" {
			format = constants.PipelineTypeGo
		}

		// initialize variable map if not parsed from config
		if len(template.Variables) == 0 {
			template.Variables = make(map[string]interface{})
		}

		// inject template name into variables
		template.Variables["VELA_TEMPLATE_NAME"] = template.Name

		parsed, _, _, err := c.Parse(bytes, format, template)
		if err != nil {
			return nil, err
		}

		// if template parsed contains a template reference, recurse with decremented depth
		if len(parsed.Templates) > 0 && parsed.Metadata.RenderInline {
			parsed, err = c.compileInline(ctx, parsed, depth-1)
			if err != nil {
				return nil, err
			}

			newPipeline.Templates = append(newPipeline.Templates, parsed.Templates...)
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

	return &newPipeline, nil
}

// compileSteps executes the workflow for converting a YAML pipeline into an executable struct.
func (c *Client) compileSteps(ctx context.Context, p *yaml.Build, _pipeline *api.Pipeline, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, *api.Pipeline, error) {
	var (
		warnings []string
		err      error
	)

	// check if the pipeline disabled the clone
	if p.Metadata.Clone == nil || *p.Metadata.Clone {
		// inject the clone step
		p, err = c.CloneStep(p)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// inject the init step
	p, err = c.InitStep(p)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the template for deploy config if exists
	p, err = c.ExpandDeployment(ctx, p, tmpls)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the templates into the steps
	p, warnings, err = c.ExpandSteps(ctx, p, tmpls, r, _pipeline.GetWarnings(), c.GetTemplateDepth())
	if err != nil {
		return nil, _pipeline, err
	}

	_pipeline.SetWarnings(warnings)

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		//
		//nolint:contextcheck // modification service has its own context with a set timeout
		p, err = c.modifyConfig(p, c.build, c.repo)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// validate the yaml configuration
	err = c.ValidateYAML(p)
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

	// validate the yaml configuration
	err = c.ValidatePipeline(build)
	if err != nil {
		return nil, _pipeline, err
	}

	return build, _pipeline, nil
}

// compileStages executes the workflow for converting a YAML pipeline into an executable struct.
func (c *Client) compileStages(ctx context.Context, p *yaml.Build, _pipeline *api.Pipeline, tmpls map[string]*yaml.Template, r *pipeline.RuleData) (*pipeline.Build, *api.Pipeline, error) {
	var (
		warnings []string
		err      error
	)

	// check if the pipeline disabled the clone
	if p.Metadata.Clone == nil || *p.Metadata.Clone {
		// inject the clone stage
		p, err = c.CloneStage(p)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// inject the init stage
	p, err = c.InitStage(p)
	if err != nil {
		return nil, _pipeline, err
	}

	// inject the templates into the stages
	p, warnings, err = c.ExpandStages(ctx, p, tmpls, r, _pipeline.GetWarnings())
	if err != nil {
		return nil, _pipeline, err
	}

	_pipeline.SetWarnings(warnings)

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		//
		//nolint:contextcheck // modification service has its own context with a set timeout
		p, err = c.modifyConfig(p, c.build, c.repo)
		if err != nil {
			return nil, _pipeline, err
		}
	}

	// validate the yaml configuration
	err = c.ValidateYAML(p)
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

	// validate the final pipeline configuration
	err = c.ValidatePipeline(build)
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
func (c *Client) modifyConfig(build *yaml.Build, apiBuild *api.Build, repo *api.Repo) (*yaml.Build, error) {
	// create request to send to endpoint
	data, err := yml.Marshal(build)
	if err != nil {
		return nil, err
	}

	modReq := &ModifyRequest{
		Pipeline: string(data),
		Build:    repo.GetCounter() + 1, // this is an assumption
		Repo:     repo.GetName(),
		Org:      repo.GetOrg(),
		User:     apiBuild.GetAuthor(),
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

// purgeStepsLite is a helper function that uses ruledata to purge steps that do not meet criteria for execution.
//
// used solely for CompileLite.
func purgeStepsLite(steps yaml.StepSlice, ruleData *pipeline.RuleData) yaml.StepSlice {
	purgedSteps := yaml.StepSlice{}

	for _, s := range steps {
		cRuleset := s.Ruleset.ToPipeline()
		if match, err := ruleData.Match(*cRuleset); err == nil && match {
			purgedSteps = append(purgedSteps, s)
		}
	}

	return purgedSteps
}

// purgeServicesLite is a helper function that uses ruledata to purge services that do not meet criteria for execution.
//
// used solely for CompileLite.
func purgeServicesLite(services yaml.ServiceSlice, ruleData *pipeline.RuleData) yaml.ServiceSlice {
	purgedServices := yaml.ServiceSlice{}

	for _, s := range services {
		cRuleset := s.Ruleset.ToPipeline()
		if match, err := ruleData.Match(*cRuleset); err == nil && match {
			purgedServices = append(purgedServices, s)
		}
	}

	return purgedServices
}

// purgeSecretsLite is a helper function that uses ruledata to purge secrets that do not meet criteria for execution.
//
// used solely for CompileLite.
func purgeSecretsLite(secrets yaml.SecretSlice, ruleData *pipeline.RuleData) yaml.SecretSlice {
	purgedSecrets := yaml.SecretSlice{}

	for _, sec := range secrets {
		if sec.Origin.Empty() {
			purgedSecrets = append(purgedSecrets, sec)

			continue
		}

		cRuleset := sec.Origin.Ruleset.ToPipeline()
		if match, err := ruleData.Match(*cRuleset); err == nil && match {
			purgedSecrets = append(purgedSecrets, sec)
		}
	}

	return purgedSecrets
}
