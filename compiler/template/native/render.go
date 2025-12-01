// SPDX-License-Identifier: Apache-2.0

package native

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/go-vela/server/compiler/types/raw"
	types "github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/internal"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, name string, tName string, environment raw.StringSliceMap, variables map[string]interface{}) (*types.Build, []string, error) {
	buffer := new(bytes.Buffer)

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
		return nil, nil, fmt.Errorf("unable to parse template %s: %w", tName, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, variables)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to execute template %s: %w", tName, err)
	}

	// unmarshal the template to the pipeline
	config, warnings, err := internal.ParseYAML(buffer.Bytes(), tName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", name, newStep.Name)
	}

	return &types.Build{
			Metadata:    config.Metadata,
			Steps:       config.Steps,
			Secrets:     config.Secrets,
			Services:    config.Services,
			Environment: config.Environment,
			Templates:   config.Templates,
			Deployment:  config.Deployment,
		},
		warnings,
		nil
}

// RenderBuild renders the templated build.
func RenderBuild(tmpl string, b string, envs map[string]string, variables map[string]interface{}) (*types.Build, []string, error) {
	buffer := new(bytes.Buffer)

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
		return nil, nil, err
	}

	// execute the template
	err = t.Execute(buffer, variables)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to execute template: %w", err)
	}

	// unmarshal the template to the pipeline
	config, warnings, err := internal.ParseYAML(buffer.Bytes(), "")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	return config, warnings, nil
}
