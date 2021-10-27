// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"io/ioutil"
	"testing"

	goyaml "github.com/buildkite/yaml"
	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestStarlark_RenderStep(t *testing.T) {
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
		{"basic", args{velaFile: "testdata/step/basic/step.yml", starlarkFile: "testdata/step/basic/template.py"}, "testdata/step/basic/want.yml", false},
		{"with method", args{velaFile: "testdata/step/with_method/step.yml", starlarkFile: "testdata/step/with_method/template.star"}, "testdata/step/with_method/want.yml", false},
		{"user vars", args{velaFile: "testdata/step/with_vars/step.yml", starlarkFile: "testdata/step/with_vars/template.star"}, "testdata/step/with_vars/want.yml", false},
		{"platform vars", args{velaFile: "testdata/step/with_vars_plat/step.yml", starlarkFile: "testdata/step/with_vars_plat/template.star"}, "testdata/step/with_vars_plat/want.yml", false},
		{"cancel due to complexity", args{velaFile: "testdata/step/cancel/step.yml", starlarkFile: "testdata/step/cancel/template.star"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := ioutil.ReadFile(tt.args.velaFile)
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

			tmpl, err := ioutil.ReadFile(tt.args.starlarkFile)
			if err != nil {
				t.Error(err)
			}

			steps, secrets, services, environment, err := RenderStep(string(tmpl), b.Steps[0])
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderStep() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true {
				wFile, err := ioutil.ReadFile(tt.wantFile)
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

				if diff := cmp.Diff(wantSteps, steps); diff != "" {
					t.Errorf("RenderStep() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantSecrets, secrets); diff != "" {
					t.Errorf("RenderStep() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantServices, services); diff != "" {
					t.Errorf("RenderStep() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantEnvironment, environment); diff != "" {
					t.Errorf("RenderStep() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestNative_RenderBuild(t *testing.T) {
	type args struct {
		velaFile string
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{"steps", args{velaFile: "testdata/build/basic/build.star"}, "testdata/build/basic/want.yml", false},
		{"stages", args{velaFile: "testdata/build/basic_stages/build.star"}, "testdata/build/basic_stages/want.yml", false},
		{"conditional match", args{velaFile: "testdata/build/conditional/build.star"}, "testdata/build/conditional/want.yml", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := ioutil.ReadFile(tt.args.velaFile)
			if err != nil {
				t.Error(err)
			}

			got, err := RenderBuild(string(sFile), map[string]string{
				"VELA_REPO_FULL_NAME": "octocat/hello-world",
				"VELA_BUILD_BRANCH":   "master",
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true {
				wFile, err := ioutil.ReadFile(tt.wantFile)
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
