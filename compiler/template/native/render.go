// SPDX-License-Identifier: Apache-2.0

package native

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/buildkite/yaml"

	"github.com/go-vela/server/compiler/types/raw"
	bkTypes "github.com/go-vela/server/compiler/types/yaml/buildkite"
	types "github.com/go-vela/server/compiler/types/yaml/yaml"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, name string, tName string, environment raw.StringSliceMap, variables map[string]interface{}) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(bkTypes.Build)

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

	return &types.Build{Metadata: *config.Metadata.ToYAML(), Steps: *config.Steps.ToYAML(), Secrets: *config.Secrets.ToYAML(), Services: *config.Services.ToYAML(), Environment: config.Environment, Templates: *config.Templates.ToYAML(), Deployment: *config.Deployment.ToYAML()}, nil
}

// RenderBuild renders the templated build.
func RenderBuild(tmpl string, b string, envs map[string]string, variables map[string]interface{}) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(bkTypes.Build)

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

	return config.ToYAML(), nil
}
