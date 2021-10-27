// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/go-vela/types/raw"

	yaml "github.com/buildkite/yaml"
	types "github.com/go-vela/types/yaml"
	"go.starlark.net/starlark"
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

// RenderStep combines the template with the step in the yaml pipeline.
//
// nolint: funlen,lll // ignore function length due to comments
func RenderStep(tmpl string, s *types.Step) (types.StepSlice, types.SecretSlice, types.ServiceSlice, raw.StringSliceMap, error) {
	config := new(types.Build)

	thread := &starlark.Thread{Name: s.Name}
	// arbitrarily limiting the steps of the thread to 5000 to help prevent infinite loops
	// may need to further investigate spawning a separate POSIX process if user input is problematic
	// see https://github.com/google/starlark-go/issues/160#issuecomment-466794230 for further details
	//
	// nolint: gomnd // ignore magic number
	thread.SetMaxExecutionSteps(5000)
	globals, err := starlark.ExecFile(thread, s.Template.Name, tmpl, nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// check the provided template has a main function
	mainVal, ok := globals["main"]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("%s: %s", ErrMissingMainFunc, s.Template.Name)
	}

	// check the provided main is a function
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("%s: %s", ErrInvalidMainFunc, s.Template.Name)
	}

	// load the user provided vars into a starlark type
	userVars, err := convertTemplateVars(s.Template.Variables)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// load the platform provided vars into a starlark type
	velaVars, err := convertPlatformVars(s.Environment, s.Name)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// add the user and platform vars to a context to be used
	// within the template caller i.e. ctx["vela"] or ctx["vars"]
	context := starlark.NewDict(0)
	err = context.SetKey(starlark.String("vela"), velaVars)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = context.SetKey(starlark.String("vars"), userVars)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	args := starlark.Tuple([]starlark.Value{context})

	// execute Starlark program from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, nil, nil, nil, err
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
				return nil, nil, nil, nil, err
			}
			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")
		err = writeJSON(buf, v)
		if err != nil {
			return nil, nil, nil, nil, err
		}
	default:
		return nil, nil, nil, nil, fmt.Errorf("%s: %s", ErrInvalidPipelineReturn, mainVal.Type())
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buf.Bytes(), config)
	if err != nil {
		// nolint: lll // ignore long line length due to return args
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
	config := new(types.Build)

	thread := &starlark.Thread{Name: "templated-base"}
	// arbitrarily limiting the steps of the thread to 5000 to help prevent infinite loops
	// may need to further investigate spawning a separate POSIX process if user input is problematic
	// see https://github.com/google/starlark-go/issues/160#issuecomment-466794230 for further details
	//
	// nolint: gomnd // ignore magic number
	thread.SetMaxExecutionSteps(5000)
	globals, err := starlark.ExecFile(thread, "templated-base", b, nil)
	if err != nil {
		return nil, err
	}

	// check the provided template has a main function
	mainVal, ok := globals["main"]
	if !ok {
		return nil, fmt.Errorf("%s: %s", ErrMissingMainFunc, "templated-base")
	}

	// check the provided main is a function
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, fmt.Errorf("%s: %s", ErrInvalidMainFunc, "templated-base")
	}

	// load the platform provided vars into a starlark type
	velaVars, err := convertPlatformVars(envs, "")
	if err != nil {
		return nil, err
	}

	// add the user and platform vars to a context to be used
	// within the template caller i.e. ctx["vela"] or ctx["vars"]
	context := starlark.NewDict(0)
	err = context.SetKey(starlark.String("vela"), velaVars)
	if err != nil {
		return nil, err
	}

	args := starlark.Tuple([]starlark.Value{context})

	// execute Starlark program from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, err
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
				return nil, err
			}
			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")
		err = writeJSON(buf, v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s: %s", ErrInvalidPipelineReturn, mainVal.Type())
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buf.Bytes(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	return config, nil
}
