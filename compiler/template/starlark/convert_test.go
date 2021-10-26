// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
	"go.starlark.net/starlark"
)

func TestStarlark_Render_convertTemplateVars(t *testing.T) {
	// setup types
	tags := starlark.Tuple(nil)
	tags = append(tags, starlark.String("latest"))
	tags = append(tags, starlark.String("1.14"))
	tags = append(tags, starlark.String("1.15"))

	commands := starlark.NewDict(16)
	err := commands.SetKey(starlark.String("test"), starlark.String("go test ./..."))
	if err != nil {
		t.Error(err)
	}

	strWant := starlark.NewDict(0)
	err = strWant.SetKey(starlark.String("pull"), starlark.String("always"))
	if err != nil {
		t.Error(err)
	}

	arrayWant := starlark.NewDict(0)
	err = arrayWant.SetKey(starlark.String("tags"), tags)
	if err != nil {
		t.Error(err)
	}

	mapWant := starlark.NewDict(0)
	err = mapWant.SetKey(starlark.String("commands"), commands)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		args map[string]interface{}
		want *starlark.Dict
	}{
		{
			name: "test for a user passed string",
			args: map[string]interface{}{"pull": "always"},
			want: strWant,
		},
		{
			name: "test for a user passed array",
			args: map[string]interface{}{"tags": []string{"latest", "1.14", "1.15"}},
			want: arrayWant,
		},
		{
			name: "test for a user passed map",
			args: map[string]interface{}{"commands": map[string]string{"test": "go test ./..."}},
			want: mapWant,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTemplateVars(tt.args)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTemplateVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStarlark_Render_velaEnvironmentData(t *testing.T) {
	// setup types
	build := starlark.NewDict(1)
	err := build.SetKey(starlark.String("author"), starlark.String("octocat"))
	if err != nil {
		t.Error(err)
	}

	repo := starlark.NewDict(1)
	err = repo.SetKey(starlark.String("full_name"), starlark.String("go-vela/hello-world"))
	if err != nil {
		t.Error(err)
	}

	user := starlark.NewDict(1)
	err = user.SetKey(starlark.String("admin"), starlark.String("true"))
	if err != nil {
		t.Error(err)
	}

	system := starlark.NewDict(2)
	err = system.SetKey(starlark.String("template_name"), starlark.String("foo"))
	if err != nil {
		t.Error(err)
	}
	err = system.SetKey(starlark.String("workspace"), starlark.String("/vela/src/github.com/go-vela/hello-world"))
	if err != nil {
		t.Error(err)
	}

	withAllPre := starlark.NewDict(0)
	err = withAllPre.SetKey(starlark.String("build"), build)
	if err != nil {
		t.Error(err)
	}
	err = withAllPre.SetKey(starlark.String("repo"), repo)
	if err != nil {
		t.Error(err)
	}
	err = withAllPre.SetKey(starlark.String("user"), user)
	if err != nil {
		t.Error(err)
	}
	err = withAllPre.SetKey(starlark.String("system"), system)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name         string
		slice        raw.StringSliceMap
		templateName string
		want         *starlark.Dict
		wantErr      bool
	}{
		{
			name: "with all vela prefixed var",
			slice: raw.StringSliceMap{
				"VELA_BUILD_AUTHOR":   "octocat",
				"VELA_REPO_FULL_NAME": "go-vela/hello-world",
				"VELA_USER_ADMIN":     "true",
				"VELA_WORKSPACE":      "/vela/src/github.com/go-vela/hello-world",
			},
			templateName: "foo",
			want:         withAllPre,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertPlatformVars(tt.slice, tt.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertPlatformVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPlatformVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
