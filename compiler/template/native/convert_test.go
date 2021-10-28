package native

import (
	"github.com/go-vela/types/raw"
	"reflect"
	"testing"
)

func Test_convertPlatformVars(t *testing.T) {
	tests := []struct {
		name         string
		slice        raw.StringSliceMap
		templateName string
		want         raw.StringSliceMap
	}{
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
			name: "with combination of vela and user vars",
			slice: raw.StringSliceMap{
				"VELA_BUILD_AUTHOR":   "octocat",
				"VELA_REPO_FULL_NAME": "go-vela/hello-world",
				"FOO_VAR1":            "test1",
				"BAR_VAR1":            "test2",
			},
			templateName: "foo",
			want:         raw.StringSliceMap{"build_author": "octocat", "repo_full_name": "go-vela/hello-world", "template_name": "foo"},
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
