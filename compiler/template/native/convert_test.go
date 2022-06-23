// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
)

func Test_convertPlatformVars(t *testing.T) {
	tests := []struct {
		name         string
		slice        raw.StringSliceMap
		templateName string
		want         raw.StringSliceMap
	}{
		{
			name: "with all deployment parameter prefixed vars",
			slice: raw.StringSliceMap{
				"DEPLOYMENT_PARAMETER_IMAGE":   "alpine",
				"DEPLOYMENT_PARAMETER_VERSION": "3.14",
			},
			templateName: "foo",
			want:         raw.StringSliceMap{"image": "alpine", "version": "3.14", "template_name": "foo"},
		},
		{
			name: "with all vela prefixed vars",
			slice: raw.StringSliceMap{
				"VELA_BUILD_AUTHOR":   "octocat",
				"VELA_REPO_FULL_NAME": "go-vela/hello-world",
				"VELA_USER_ADMIN":     "true",
				"VELA_WORKSPACE":      "/vela/src/github.com/go-vela/hello-world",
			},
			templateName: "foo",
			want:         raw.StringSliceMap{"build_author": "octocat", "repo_full_name": "go-vela/hello-world", "user_admin": "true", "workspace": "/vela/src/github.com/go-vela/hello-world", "template_name": "foo"},
		},
		{
			name: "with combination of deployment parameter, vela, and user vars",
			slice: raw.StringSliceMap{
				"DEPLOYMENT_PARAMETER_IMAGE":   "alpine",
				"DEPLOYMENT_PARAMETER_VERSION": "3.14",
				"VELA_BUILD_AUTHOR":            "octocat",
				"VELA_REPO_FULL_NAME":          "go-vela/hello-world",
				"FOO_VAR1":                     "test1",
				"BAR_VAR1":                     "test2",
			},
			templateName: "foo",
			want:         raw.StringSliceMap{"image": "alpine", "version": "3.14", "build_author": "octocat", "repo_full_name": "go-vela/hello-world", "template_name": "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertPlatformVars(tt.slice, tt.templateName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPlatformVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_funcHandler_returnPlatformVar(t *testing.T) {
	type fields struct {
		envs raw.StringSliceMap
	}

	type args struct {
		input string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "existing deployment parameter without prefix (lowercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"image": "alpine",
				},
			},
			args: args{input: "image"},
			want: "alpine",
		},
		{
			name: "existing deployment parameter without prefix (uppercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"image": "alpine",
				},
			},
			args: args{input: "IMAGE"},
			want: "alpine",
		},
		{
			name: "existing deployment parameter with prefix (lowercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"image": "alpine",
				},
			},
			args: args{input: "deployment_parameter_image"},
			want: "alpine",
		},
		{
			name: "existing deployment parameter with prefix (uppercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"image": "alpine",
				},
			},
			args: args{input: "DEPLOYMENT_PARAMETER_IMAGE"},
			want: "alpine",
		},
		{
			name: "existing platform var without prefix (lowercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"build_author": "octocat",
				},
			},
			args: args{input: "build_author"},
			want: "octocat",
		},
		{
			name: "existing platform var without prefix (uppercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"build_author": "octocat",
				},
			},
			args: args{input: "BUILD_AUTHOR"},
			want: "octocat",
		},
		{
			name: "existing platform var with prefix (lowercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"build_author": "octocat",
				},
			},
			args: args{input: "vela_build_author"},
			want: "octocat",
		},
		{
			name: "existing platform var with prefix (uppercase)",
			fields: fields{
				envs: raw.StringSliceMap{
					"build_author": "octocat",
				},
			},
			args: args{input: "VELA_BUILD_AUTHOR"},
			want: "octocat",
		},
		{
			name: "non existent var",
			fields: fields{
				envs: raw.StringSliceMap{
					"build_author": "octocat",
				},
			},
			args: args{input: "foo"},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := funcHandler{
				envs: tt.fields.envs,
			}
			if got := h.returnPlatformVar(tt.args.input); got != tt.want {
				t.Errorf("returnPlatformVar() = %v, want %v", got, tt.want)
			}
		})
	}
}
