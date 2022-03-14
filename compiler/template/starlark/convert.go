// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"strings"

	"github.com/go-vela/types/raw"
	"go.starlark.net/starlark"
)

// convertTemplateVars takes template variables and converts
// them to a starlark string dictionary for template reference.
//
// Example Usage within template: ctx["vars"]["message"] = "Hello, World!"
//
// Explanation of type "starlark.StringDict":
// https://pkg.go.dev/go.starlark.net/starlark#StringDict
func convertTemplateVars(m map[string]interface{}) (*starlark.Dict, error) {
	dict := starlark.NewDict(0)

	// loop through user vars converting provided types to starlark primitives
	for key, value := range m {
		val, err := toStarlark(value)
		if err != nil {
			return nil, err
		}

		err = dict.SetKey(starlark.String(key), val)
		if err != nil {
			return nil, err
		}
	}

	return dict, nil
}

// convertPlatformVars takes the platform injected variables
// within the step environment block and converts them to a
// starlark string dictionary.
//
// Example Usage within template: ctx["vela"]["build"]["number"] = 1
//
// Explanation of type "starlark.StringDict":
// https://pkg.go.dev/go.starlark.net/starlark#StringDict
func convertPlatformVars(slice raw.StringSliceMap, name string) (*starlark.Dict, error) {
	build := starlark.NewDict(0)
	repo := starlark.NewDict(0)
	user := starlark.NewDict(0)
	system := starlark.NewDict(0)
	dict := starlark.NewDict(0)

	err := dict.SetKey(starlark.String("build"), build)
	if err != nil {
		return nil, err
	}

	err = dict.SetKey(starlark.String("repo"), repo)
	if err != nil {
		return nil, err
	}

	err = dict.SetKey(starlark.String("user"), user)
	if err != nil {
		return nil, err
	}

	err = dict.SetKey(starlark.String("system"), system)
	if err != nil {
		return nil, err
	}

	err = system.SetKey(starlark.String("template_name"), starlark.String(name))
	if err != nil {
		return nil, err
	}

	for key, value := range slice {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "vela_") {
			key = strings.TrimPrefix(key, "vela_")

			switch {
			case strings.HasPrefix(key, "build_"):
				err := build.SetKey(starlark.String(strings.TrimPrefix(key, "build_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(key, "repo_"):
				err := repo.SetKey(starlark.String(strings.TrimPrefix(key, "repo_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(key, "user_"):
				err := user.SetKey(starlark.String(strings.TrimPrefix(key, "user_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			default:
				err := system.SetKey(starlark.String(key), starlark.String(value))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return dict, nil
}
