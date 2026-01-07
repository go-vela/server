// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/go-vela/server/compiler/registry"
	"github.com/go-vela/server/compiler/template/native"
	"github.com/go-vela/server/compiler/template/starlark"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
)

// ExpandStages injects the template for each
// templated step in every stage in a yaml configuration.
func (c *Client) ExpandStages(ctx context.Context, s *yaml.Build, tmpls map[string]*yaml.Template, r *pipeline.RuleData, warnings []string) (*yaml.Build, []string, error) {
	var (
		p   *yaml.Build
		err error
	)

	if len(tmpls) == 0 {
		return s, warnings, nil
	}

	// iterate through all stages
	for _, stage := range s.Stages {
		// inject the templates into the steps for the stage
		p, warnings, err = c.ExpandSteps(ctx, &yaml.Build{Steps: stage.Steps, Secrets: s.Secrets, Services: s.Services, Environment: s.Environment}, tmpls, r, warnings, c.GetTemplateDepth())
		if err != nil {
			return nil, warnings, err
		}

		stage.Steps = p.Steps
		s.Secrets = p.Secrets
		s.Services = p.Services
		s.Environment = p.Environment
	}

	return s, warnings, nil
}

// ExpandSteps injects the template for each
// templated step in a yaml configuration.
//
//nolint:funlen,gocyclo // ignore function length
func (c *Client) ExpandSteps(ctx context.Context, s *yaml.Build, tmpls map[string]*yaml.Template, r *pipeline.RuleData, warnings []string, depth int) (*yaml.Build, []string, error) {
	if len(tmpls) == 0 {
		return s, warnings, nil
	}

	// return if max template depth has been reached
	if depth == 0 {
		retErr := fmt.Errorf("max template depth of %d exceeded", c.GetTemplateDepth())

		return s, warnings, retErr
	}

	steps := yaml.StepSlice{}
	secrets := s.Secrets
	services := s.Services
	environment := s.Environment
	templates := s.Templates

	if len(environment) == 0 {
		environment = make(raw.StringSliceMap)
	}

	// iterate through each step
	for _, step := range s.Steps {
		// skip if no template is provided for the step
		if len(step.Template.Name) == 0 {
			// add existing step if no template
			steps = append(steps, step)
			continue
		}

		// lookup step template name
		tmpl, ok := tmpls[step.Template.Name]
		if !ok {
			return s, warnings, fmt.Errorf("missing template source for template %s in pipeline for step %s", step.Template.Name, step.Name)
		}

		// if ruledata is nil (CompileLite), continue with expansion
		if r != nil {
			// form a one-step pipeline to prep for purge check
			check := &yaml.StepSlice{step}
			pipeline := &pipeline.Build{
				Steps: *check.ToPipeline(),
			}

			pipeline, err := pipeline.Purge(r)
			if err != nil {
				return nil, warnings, fmt.Errorf("unable to purge pipeline: %w", err)
			}

			// if step purged, do not proceed with expansion
			if len(pipeline.Steps) == 0 {
				continue
			}
		}

		// Create some default global environment inject vars
		// these are used below to overwrite to an empty
		// map if they should not be injected into a container
		envGlobalSteps := s.Environment

		if !s.Metadata.HasEnvironment("steps") {
			envGlobalSteps = make(raw.StringSliceMap)
		}

		// inject environment information for template
		step, err := c.EnvironmentStep(step, envGlobalSteps)
		if err != nil {
			return s, warnings, err
		}

		var (
			bytes []byte
			found bool
		)

		if bytes, found = c.TemplateCache[tmpl.Source]; !found {
			bytes, err = c.getTemplate(ctx, tmpl, step.Template.Name)
			if err != nil {
				return s, warnings, err
			}
		}

		// initialize variable map if not parsed from config
		if len(step.Template.Variables) == 0 {
			step.Template.Variables = make(map[string]interface{})
		}

		// inject template name into variables
		step.Template.Variables["VELA_TEMPLATE_NAME"] = step.Template.Name

		tmplBuild, tmplWarnings, err := c.mergeTemplate(bytes, tmpl, step)
		if err != nil {
			return s, warnings, err
		}

		warnings = append(warnings, tmplWarnings...)

		// if template references other templates, expand again
		if len(tmplBuild.Templates) != 0 {
			// if the tmplBuild has render_inline but the parent build does not, abort
			if tmplBuild.Metadata.RenderInline && !s.Metadata.RenderInline {
				return s, warnings, fmt.Errorf("cannot use render_inline inside a called template (%s)", step.Template.Name)
			}

			templates = append(templates, tmplBuild.Templates...)

			tmplBuild, warnings, err = c.ExpandSteps(ctx, tmplBuild, mapFromTemplates(tmplBuild.Templates), r, warnings, depth-1)
			if err != nil {
				return s, warnings, err
			}
		}

		// loop over secrets within template
		for _, secret := range tmplBuild.Secrets {
			found := false
			// loop over secrets within base configuration
			for _, sec := range secrets {
				// check if the template secret and base secret name match
				if sec.Name == secret.Name {
					found = true
				}
			}

			// only append template secret if it does not exist within base configuration
			if !secret.Origin.Empty() || !found {
				secrets = append(secrets, secret)
			}
		}

		// loop over services within template
		for _, service := range tmplBuild.Services {
			found := false

			for _, serv := range services {
				if serv.Name == service.Name {
					found = true
				}
			}

			// only append template service if it does not exist within base configuration
			if !found {
				services = append(services, service)
			}
		}

		// loop over environment within template
		for key, value := range tmplBuild.Environment {
			found := false

			for env := range environment {
				if key == env {
					found = true
				}
			}

			// only append template environment if it does not exist within base configuration
			if !found {
				environment[key] = value
			}
		}

		// add templated steps
		steps = append(steps, tmplBuild.Steps...)
	}

	s.Steps = steps
	s.Secrets = secrets
	s.Services = services
	s.Environment = environment
	s.Templates = templates

	return s, warnings, nil
}

// ExpandDeployment injects the template for a
// templated deployment config in a yaml configuration.
func (c *Client) ExpandDeployment(ctx context.Context, b *yaml.Build, tmpls map[string]*yaml.Template) (*yaml.Build, error) {
	if len(tmpls) == 0 {
		return b, nil
	}

	if len(b.Deployment.Template.Name) == 0 {
		return b, nil
	}

	// lookup step template name
	tmpl, ok := tmpls[b.Deployment.Template.Name]
	if !ok {
		return b, fmt.Errorf("missing template source for template %s in pipeline for deployment config", b.Deployment.Template.Name)
	}

	bytes, err := c.getTemplate(ctx, tmpl, b.Deployment.Template.Name)
	if err != nil {
		return b, err
	}

	// initialize variable map if not parsed from config
	if len(b.Deployment.Template.Variables) == 0 {
		b.Deployment.Template.Variables = make(map[string]interface{})
	}

	tmplBuild, _, err := c.mergeDeployTemplate(bytes, tmpl, &b.Deployment)
	if err != nil {
		return b, err
	}

	b.Deployment = tmplBuild.Deployment

	return b, nil
}

func (c *Client) getTemplate(ctx context.Context, tmpl *yaml.Template, name string) ([]byte, error) {
	var (
		bytes []byte
		err   error
	)

	switch {
	case c.local:
		a := &afero.Afero{
			Fs: afero.NewOsFs(),
		}

		// iterate over locally provided templates
		for _, t := range c.localTemplates {
			parts := strings.Split(t, ":")
			if len(parts) != 2 {
				return nil, fmt.Errorf("local templates must be provided in the form <name>:<path>, got %s", t)
			}

			// if local template has a match, read file path provided
			if strings.EqualFold(tmpl.Name, parts[0]) {
				bytes, err = a.ReadFile(parts[1])
				if err != nil {
					return bytes, err
				}

				return bytes, nil
			}
		}

		// file type templates can be retrieved locally using `source`
		if strings.EqualFold(tmpl.Type, "file") {
			bytes, err = a.ReadFile(tmpl.Source)
			if err != nil {
				return nil, fmt.Errorf("unable to read file for template %s. `File` type templates must be located at `source` or supplied to local template files", tmpl.Name)
			}

			return bytes, nil
		}

		// local exec may still request remote templates
		if !strings.EqualFold(tmpl.Type, "github") {
			return nil, fmt.Errorf("unable to find template %s: not supplied in list %s", tmpl.Name, c.localTemplates)
		}

		fallthrough

	case strings.EqualFold(tmpl.Type, "github"):
		// parse source from template
		src, err := c.Github.Parse(tmpl.Source)
		if err != nil {
			return bytes, fmt.Errorf("invalid template source provided for %s: %w", name, err)
		}

		// pull from github without auth when the host isn't provided or is set to github.com
		if c.UsePrivateGithub {
			logrus.WithFields(logrus.Fields{
				"org":  src.Org,
				"repo": src.Repo,
				"path": src.Name,
				"host": src.Host,
			}).Tracef("Using authenticated GitHub client to pull template")

			// verify private GitHub is actually set up
			if c.PrivateGithub == nil {
				return nil, fmt.Errorf("unable to fetch template %s: missing credentials", src.Name)
			}

			// use private (authenticated) github instance to pull from
			bytes, err = c.PrivateGithub.Template(ctx, c.user, src)
			if err != nil {
				return bytes, err
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"org":  src.Org,
				"repo": src.Repo,
				"path": src.Name,
				"host": src.Host,
			}).Tracef("Using GitHub client to pull template")

			bytes, err = c.Github.Template(ctx, nil, src)
			if err != nil {
				return bytes, err
			}
		}

	case strings.EqualFold(tmpl.Type, "file"):
		src := &registry.Source{
			Org:  c.repo.GetOrg(),
			Repo: c.repo.GetName(),
			Name: tmpl.Source,
			Ref:  c.commit,
		}

		if !c.UsePrivateGithub {
			logrus.WithFields(logrus.Fields{
				"org":  src.Org,
				"repo": src.Repo,
				"path": src.Name,
			}).Tracef("Using GitHub client to pull template")

			bytes, err = c.Github.Template(ctx, nil, src)
			if err != nil {
				return bytes, err
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"org":  src.Org,
				"repo": src.Repo,
				"path": src.Name,
			}).Tracef("Using authenticated GitHub client to pull template")

			if c.PrivateGithub == nil {
				return nil, fmt.Errorf("unable to fetch template %s: missing credentials", src.Name)
			}

			// use private (authenticated) github instance to pull from
			bytes, err = c.PrivateGithub.Template(ctx, c.user, src)
			if err != nil {
				return bytes, err
			}
		}

	default:
		return bytes, fmt.Errorf("unsupported template type: %v", tmpl.Type)
	}

	c.TemplateCache[tmpl.Source] = bytes

	return bytes, nil
}

//nolint:lll // ignore long line length due to input arguments
func (c *Client) mergeTemplate(bytes []byte, tmpl *yaml.Template, step *yaml.Step) (*yaml.Build, []string, error) {
	switch tmpl.Format {
	case constants.PipelineTypeGo, constants.PipelineTypeGoAlt, "":
		//nolint:lll // ignore long line length due to return
		return native.Render(string(bytes), step.Name, step.Template.Name, step.Environment, step.Template.Variables)
	case constants.PipelineTypeStarlark:
		//nolint:lll // ignore long line length due to return
		return starlark.Render(string(bytes), step.Name, step.Template.Name, step.Environment, step.Template.Variables, c.GetStarlarkExecLimit())
	default:
		//nolint:lll // ignore long line length due to return
		return &yaml.Build{}, nil, fmt.Errorf("format of %s is unsupported", tmpl.Format)
	}
}

func (c *Client) mergeDeployTemplate(bytes []byte, tmpl *yaml.Template, d *yaml.Deployment) (*yaml.Build, []string, error) {
	switch tmpl.Format {
	case constants.PipelineTypeGo, constants.PipelineTypeGoAlt, "":
		//nolint:lll // ignore long line length due to return
		return native.Render(string(bytes), "", d.Template.Name, make(raw.StringSliceMap), d.Template.Variables)
	case constants.PipelineTypeStarlark:
		//nolint:lll // ignore long line length due to return
		return starlark.Render(string(bytes), "", d.Template.Name, make(raw.StringSliceMap), d.Template.Variables, c.GetStarlarkExecLimit())
	default:
		//nolint:lll // ignore long line length due to return
		return &yaml.Build{}, nil, fmt.Errorf("format of %s is unsupported", tmpl.Format)
	}
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
