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

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, name string, tName string, environment raw.StringSliceMap, variables map[string]interface{}) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(environment, name)}
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
	t, err := template.New(name).Funcs(sf).Funcs(templateFuncMap).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %s: %w", tName, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, variables)
	if err != nil {
		return nil, fmt.Errorf("unable to execute template %s: %w", tName, err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", name, newStep.Name)
	}

	return &types.Build{Steps: config.Steps, Secrets: config.Secrets, Services: config.Services, Environment: config.Environment}, nil
}

// RenderBuild renders the templated build.
func RenderBuild(tmpl string, b string, envs map[string]string, variables map[string]interface{}) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(envs, tmpl)}
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
	t, err := template.New(tmpl).Funcs(sf).Funcs(templateFuncMap).Parse(b)
	if err != nil {
		return nil, err
	}

	// execute the template
	err = t.Execute(buffer, variables)
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
