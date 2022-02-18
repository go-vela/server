// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/template/native"
	"github.com/go-vela/server/compiler/template/starlark"
	"github.com/spf13/afero"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/sirupsen/logrus"
)

// ExpandStages injects the template for each
// templated step in every stage in a yaml configuration.
func (c *client) ExpandStages(s *yaml.Build, tmpls map[string]*yaml.Template) (yaml.StageSlice, yaml.SecretSlice, yaml.ServiceSlice, raw.StringSliceMap, error) {
	// iterate through all stages
	for _, stage := range s.Stages {
		// inject the templates into the steps for the stage
		steps, secrets, services, environment, err := c.ExpandSteps(&yaml.Build{Steps: stage.Steps, Secrets: s.Secrets, Services: s.Services, Environment: s.Environment}, tmpls)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		stage.Steps = steps
		s.Secrets = secrets
		s.Services = services
		s.Environment = environment
	}

	return s.Stages, s.Secrets, s.Services, s.Environment, nil
}

// ExpandSteps injects the template for each
// templated step in a yaml configuration.
//
// nolint: lll,funlen,gocyclo // ignore long line length due to variable names
func (c *client) ExpandSteps(s *yaml.Build, tmpls map[string]*yaml.Template) (yaml.StepSlice, yaml.SecretSlice, yaml.ServiceSlice, raw.StringSliceMap, error) {
	steps := yaml.StepSlice{}
	secrets := s.Secrets
	services := s.Services
	environment := s.Environment

	if len(environment) == 0 {
		environment = make(raw.StringSliceMap)
	}

	// iterate through each step
	for _, step := range s.Steps {
		// nolint: ineffassign,staticcheck // ignore ineffectual assignment
		bytes := []byte{}

		// skip if no template is provided for the step
		if len(step.Template.Name) == 0 {
			// add existing step if no template
			steps = append(steps, step)
			continue
		}

		// lookup step template name
		tmpl, ok := tmpls[step.Template.Name]
		if !ok {
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("missing template source for template %s in pipeline for step %s", step.Template.Name, step.Name)
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
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
		}

		switch {
		case c.local:
			a := &afero.Afero{
				Fs: afero.NewOsFs(),
			}

			bytes, err = a.ReadFile(tmpl.Source)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
			}

		case strings.EqualFold(tmpl.Type, "github"):
			// parse source from template
			src, err := c.Github.Parse(tmpl.Source)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("invalid template source provided for %s: %v", step.Template.Name, err)
			}

			// pull from github without auth when the host isn't provided or is set to github.com
			if !c.UsePrivateGithub && (len(src.Host) == 0 || strings.Contains(src.Host, "github.com")) {
				logrus.WithFields(logrus.Fields{
					"org":  src.Org,
					"repo": src.Repo,
					"path": src.Name,
					"host": src.Host,
				}).Tracef("Using GitHub client to pull template")

				bytes, err = c.Github.Template(nil, src)
				if err != nil {
					return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
				}
			} else {
				logrus.WithFields(logrus.Fields{
					"org":  src.Org,
					"repo": src.Repo,
					"path": src.Name,
					"host": src.Host,
				}).Tracef("Using authenticated GitHub client to pull template")
				// use private (authenticated) github instance to pull from
				bytes, err = c.PrivateGithub.Template(c.user, src)
				if err != nil {
					return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
				}
			}

		default:
			logrus.Errorf("Unsupported template type: %v", tmpl.Type)
			continue
		}

		var (
			tmplSteps       yaml.StepSlice
			tmplSecrets     yaml.SecretSlice
			tmplServices    yaml.ServiceSlice
			tmplEnvironment raw.StringSliceMap
		)

		// TODO: provide friendlier error messages with file type mismatches
		switch tmpl.Format {
		case "go", "golang", "":
			// render template for steps
			tmplSteps, tmplSecrets, tmplServices, tmplEnvironment, err = native.RenderStep(string(bytes), step)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
			}
		case "starlark":
			// render template for steps
			tmplSteps, tmplSecrets, tmplServices, tmplEnvironment, err = starlark.RenderStep(string(bytes), step)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, err
			}
		default:
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("format of %s is unsupported", tmpl.Format)
		}

		// loop over secrets within template
		for _, secret := range tmplSecrets {
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
		for _, service := range tmplServices {
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
		for key, value := range tmplEnvironment {
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
		steps = append(steps, tmplSteps...)
	}

	return steps, secrets, services, environment, nil
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
