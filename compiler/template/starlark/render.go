// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"bytes"
	"errors"
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"

	"github.com/go-vela/server/compiler/types/raw"
	types "github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/internal"
)

var (
	// ErrMissingMainFunc defines the error type when the
	// main function does not exist in the provided template.
	ErrMissingMainFunc = errors.New("unable to find main function in template")

	// ErrInvalidMainFunc defines the error type when the
	// main function is invalid within the provided template.
	ErrInvalidMainFunc = errors.New("invalid main function (main must be a function) in template")

	// ErrInvalidPipelineReturn defines the error type when the
	// return type is not a pipeline within the provided template.
	ErrInvalidPipelineReturn = errors.New("invalid pipeline return in template")
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, name string, tName string, environment raw.StringSliceMap, variables map[string]interface{}, limit int64) (*types.Build, []string, error) {
	thread := &starlark.Thread{Name: name}

	if limit < 0 {
		return nil, nil, fmt.Errorf("starlark exec limit must be non-negative")
	}

	thread.SetMaxExecutionSteps(uint64(limit))

	predeclared := starlark.StringDict{"struct": starlark.NewBuiltin("struct", starlarkstruct.Make)}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), thread, "templated-base", tmpl, predeclared)
	if err != nil {
		return nil, nil, err
	}

	// check the provided template has a main function
	mainVal, ok := globals["main"]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrMissingMainFunc, tName)
	}

	// check the provided main is a function
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidMainFunc, tName)
	}

	// load the user provided vars into a starlark type
	userVars, err := convertTemplateVars(variables)
	if err != nil {
		return nil, nil, err
	}

	// load the platform provided vars into a starlark type
	velaVars, err := convertPlatformVars(environment, name)
	if err != nil {
		return nil, nil, err
	}

	// add the user and platform vars to a context to be used
	// within the template caller i.e. ctx["vela"] or ctx["vars"]
	context := starlark.NewDict(0)

	err = context.SetKey(starlark.String("vela"), velaVars)
	if err != nil {
		return nil, nil, err
	}

	err = context.SetKey(starlark.String("vars"), userVars)
	if err != nil {
		return nil, nil, err
	}

	args := starlark.Tuple([]starlark.Value{context})

	// execute Starlark program from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	// extract the pipeline from the starlark program
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)

			buf.WriteString("---\n")

			err = writeJSON(buf, item)
			if err != nil {
				return nil, nil, err
			}

			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")

		err = writeJSON(buf, v)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidPipelineReturn, mainVal.Type())
	}

	// unmarshal the template to the pipeline
	config, warnings, err := internal.ParseYAML(buf.Bytes(), tName)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", name, newStep.Name)
	}

	return &types.Build{
			Steps:       config.Steps,
			Secrets:     config.Secrets,
			Services:    config.Services,
			Environment: config.Environment,
			Deployment:  config.Deployment,
		},
		warnings,
		nil
}

// RenderBuild renders the templated build.
//
//nolint:lll // ignore function length due to input args
func RenderBuild(tmpl string, b string, envs map[string]string, variables map[string]interface{}, limit int64) (*types.Build, []string, error) {
	thread := &starlark.Thread{Name: "templated-base"}

	if limit < 0 {
		return nil, nil, fmt.Errorf("starlark exec limit must be non-negative")
	}

	thread.SetMaxExecutionSteps(uint64(limit))

	predeclared := starlark.StringDict{"struct": starlark.NewBuiltin("struct", starlarkstruct.Make)}

	globals, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), thread, "templated-base", b, predeclared)
	if err != nil {
		return nil, nil, err
	}

	// check the provided template has a main function
	mainVal, ok := globals["main"]
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrMissingMainFunc, "templated-base")
	}

	// check the provided main is a function
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidMainFunc, "templated-base")
	}

	// load the user provided vars into a starlark type
	userVars, err := convertTemplateVars(variables)
	if err != nil {
		return nil, nil, err
	}

	// load the platform provided vars into a starlark type
	velaVars, err := convertPlatformVars(envs, tmpl)
	if err != nil {
		return nil, nil, err
	}

	// add the user and platform vars to a context to be used
	// within the template caller i.e. ctx["vela"] or ctx["vars"]
	context := starlark.NewDict(0)

	err = context.SetKey(starlark.String("vela"), velaVars)
	if err != nil {
		return nil, nil, err
	}

	err = context.SetKey(starlark.String("vars"), userVars)
	if err != nil {
		return nil, nil, err
	}

	args := starlark.Tuple([]starlark.Value{context})

	// execute Starlark program from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	// extract the pipeline from the starlark program
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)

			buf.WriteString("---\n")

			err = writeJSON(buf, item)
			if err != nil {
				return nil, nil, err
			}

			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")

		err = writeJSON(buf, v)
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("%w: %s", ErrInvalidPipelineReturn, mainVal.Type())
	}

	// unmarshal the template to the pipeline
	config, warnings, err := internal.ParseYAML(buf.Bytes(), "")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	return config, warnings, nil
}
