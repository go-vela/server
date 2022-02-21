// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/go-vela/types/raw"

	types "github.com/go-vela/types/yaml"

	"github.com/Masterminds/sprig/v3"

	"github.com/buildkite/yaml"
)

// RenderStep combines the template with the step in the yaml pipeline.

func RenderStep(tmpl string, s *types.Step) (types.StepSlice, types.SecretSlice, types.ServiceSlice, raw.StringSliceMap, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(s.Environment, s.Name)}
	templateFuncMap := map[string]interface{}{
		"vela":   velaFuncs.returnPlatformVar,
		"toYaml": toYAML,
	}
	// modify Masterminds/sprig functions
	// to remove OS functions
	//
	// https://masterminds.github.io/sprig/os.html
	sf := sprig.TxtFuncMap()
	delete(sf, "env")
	delete(sf, "expandenv")

	// parse the template with Masterminds/sprig functions
	//
	// https://pkg.go.dev/github.com/Masterminds/sprig?tab=doc#TxtFuncMap
	t, err := template.New(s.Name).Funcs(sf).Funcs(templateFuncMap).Parse(tmpl)
	if err != nil {
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("unable to parse template %s: %v", s.Template.Name, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, s.Template.Variables)
	if err != nil {
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("unable to execute template %s: %v", s.Template.Name, err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, raw.StringSliceMap{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	return config.Steps, config.Secrets, config.Services, config.Environment, nil
}

// RenderBuild renders the templated build.
func RenderBuild(b string, envs map[string]string) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(envs, "")}
	templateFuncMap := map[string]interface{}{
		"vela":   velaFuncs.returnPlatformVar,
		"toYaml": toYAML,
	}
	// modify Masterminds/sprig functions
	// to remove OS functions
	//
	// https://masterminds.github.io/sprig/os.html
	sf := sprig.TxtFuncMap()
	delete(sf, "env")
	delete(sf, "expandenv")

	// parse the template with Masterminds/sprig functions
	//
	// https://pkg.go.dev/github.com/Masterminds/sprig?tab=doc#TxtFuncMap
	t, err := template.New("build").Funcs(sf).Funcs(templateFuncMap).Parse(b)
	if err != nil {
		return nil, err
	}

	// execute the template
	err = t.Execute(buffer, "")
	if err != nil {
		return nil, fmt.Errorf("unable to execute template: %w", err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	return config, nil
}
