// SPDX-License-Identifier: Apache-2.0

package starlark

import (
	"os"
	"testing"

	goyaml "github.com/buildkite/yaml"
	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestStarlark_Render(t *testing.T) {
	type args struct {
		velaFile     string
		starlarkFile string
	}

	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{
			name: "basic",
			args: args{
				velaFile:     "testdata/step/basic/step.yml",
				starlarkFile: "testdata/step/basic/template.py",
			},
			wantFile: "testdata/step/basic/want.yml",
			wantErr:  false,
		},
		{
			name: "with method",
			args: args{
				velaFile:     "testdata/step/with_method/step.yml",
				starlarkFile: "testdata/step/with_method/template.star",
			},
			wantFile: "testdata/step/with_method/want.yml",
			wantErr:  false,
		},
		{
			name: "user vars",
			args: args{
				velaFile:     "testdata/step/with_vars/step.yml",
				starlarkFile: "testdata/step/with_vars/template.star",
			},
			wantFile: "testdata/step/with_vars/want.yml",
			wantErr:  false,
		},
		{
			name: "platform vars",
			args: args{
				velaFile:     "testdata/step/with_vars_plat/step.yml",
				starlarkFile: "testdata/step/with_vars_plat/template.star",
			},
			wantFile: "testdata/step/with_vars_plat/want.yml",
			wantErr:  false,
		},
		{
			name: "cancel due to complexity",
			args: args{
				velaFile:     "testdata/step/cancel/step.yml",
				starlarkFile: "testdata/step/cancel/template.star",
			},
			wantFile: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := os.ReadFile(tt.args.velaFile)
			if err != nil {
				t.Error(err)
			}
			b := &yaml.Build{}
			err = goyaml.Unmarshal(sFile, b)
			if err != nil {
				t.Error(err)
			}
			b.Steps[0].Environment = raw.StringSliceMap{
				"VELA_REPO_FULL_NAME": "octocat/hello-world",
			}

			tmpl, err := os.ReadFile(tt.args.starlarkFile)
			if err != nil {
				t.Error(err)
			}

			tmplBuild, err := Render(string(tmpl), b.Steps[0].Name, b.Steps[0].Template.Name, b.Steps[0].Environment, b.Steps[0].Template.Variables, 7500)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true {
				wFile, err := os.ReadFile(tt.wantFile)
				if err != nil {
					t.Error(err)
				}
				w := &yaml.Build{}
				err = goyaml.Unmarshal(wFile, w)
				if err != nil {
					t.Error(err)
				}
				wantSteps := w.Steps
				wantSecrets := w.Secrets
				wantServices := w.Services
				wantEnvironment := w.Environment

				if diff := cmp.Diff(wantSteps, tmplBuild.Steps); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantSecrets, tmplBuild.Secrets); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantServices, tmplBuild.Services); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantEnvironment, tmplBuild.Environment); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestNative_RenderBuild(t *testing.T) {
	noWantFile := "none"

	type args struct {
		velaFile string
	}

	tests := []struct {
		name      string
		args      args
		wantFile  string
		wantErr   bool
		execLimit uint64
	}{
		{
			name: "steps",
			args: args{
				velaFile: "testdata/build/basic/build.star",
			},
			wantFile:  "testdata/build/basic/want.yml",
			wantErr:   false,
			execLimit: 7500,
		},
		{
			name: "stages",
			args: args{
				velaFile: "testdata/build/basic_stages/build.star",
			},
			wantFile:  "testdata/build/basic_stages/want.yml",
			wantErr:   false,
			execLimit: 7500,
		},
		{
			name: "conditional match",
			args: args{
				velaFile: "testdata/build/conditional/build.star",
			},
			wantFile:  "testdata/build/conditional/want.yml",
			wantErr:   false,
			execLimit: 7500,
		},
		{
			name: "steps, with structs",
			args: args{
				velaFile: "testdata/build/with_struct/build.star",
			},
			wantFile:  "testdata/build/with_struct/want.yml",
			wantErr:   false,
			execLimit: 7500,
		},
		{
			name: "large build - exec step limit good",
			args: args{
				velaFile: "testdata/build/large/build.star",
			},
			wantFile:  "testdata/build/large/want.yml",
			wantErr:   false,
			execLimit: 7500,
		},
		{
			name: "large build - exec step limit too low",
			args: args{
				velaFile: "testdata/build/large/build.star",
			},
			wantFile:  "",
			wantErr:   true,
			execLimit: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := os.ReadFile(tt.args.velaFile)
			if err != nil {
				t.Error(err)
			}

			got, err := RenderBuild("build", string(sFile), map[string]string{
				"VELA_REPO_FULL_NAME": "octocat/hello-world",
				"VELA_BUILD_BRANCH":   "master",
				"VELA_REPO_ORG":       "octocat",
			}, map[string]interface{}{}, tt.execLimit)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true && tt.wantFile != noWantFile {
				wFile, err := os.ReadFile(tt.wantFile)
				if err != nil {
					t.Error(err)
				}
				want := &yaml.Build{}
				err = goyaml.Unmarshal(wFile, want)
				if err != nil {
					t.Error(err)
				}

				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("RenderBuild() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
