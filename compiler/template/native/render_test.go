// SPDX-License-Identifier: Apache-2.0

package native

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	goyaml "go.yaml.in/yaml/v3"

	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

func TestNative_Render(t *testing.T) {
	type args struct {
		velaFile     string
		templateFile string
	}

	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{"basic", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/basic/tmpl.yml"}, "testdata/step/basic/want.yml", false},
		{"multiline", args{velaFile: "testdata/step/multiline/step.yml", templateFile: "testdata/step/multiline/tmpl.yml"}, "testdata/step/multiline/want.yml", false},
		{"conditional match", args{velaFile: "testdata/step/conditional/step.yml", templateFile: "testdata/step/conditional/tmpl.yml"}, "testdata/step/conditional/want.yml", false},
		{"loop map", args{velaFile: "testdata/step/loop_map/step.yml", templateFile: "testdata/step/loop_map/tmpl.yml"}, "testdata/step/loop_map/want.yml", false},
		{"loop slice", args{velaFile: "testdata/step/loop_slice/step.yml", templateFile: "testdata/step/loop_slice/tmpl.yml"}, "testdata/step/loop_slice/want.yml", false},
		{"platform vars", args{velaFile: "testdata/step/with_vars_plat/step.yml", templateFile: "testdata/step/with_vars_plat/tmpl.yml"}, "testdata/step/with_vars_plat/want.yml", false},
		{"to yaml", args{velaFile: "testdata/step/to_yaml/step.yml", templateFile: "testdata/step/to_yaml/tmpl.yml"}, "testdata/step/to_yaml/want.yml", false},
		{"invalid template", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/invalid_template.yml"}, "", true},
		{"invalid variable", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/invalid_variables.yml"}, "", true},
		{"invalid yml", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/invalid.yml"}, "", true},
		{"disallowed env func", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/disallowed/tmpl_env.yml"}, "", true},
		{"disallowed expandenv func", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/disallowed/tmpl_expandenv.yml"}, "", true},
		{"bad ruleset formatting", args{velaFile: "testdata/step/basic/step.yml", templateFile: "testdata/step/bad_ruleset_format.yml"}, "", true},
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

			tmpl, err := os.ReadFile(tt.args.templateFile)
			if err != nil {
				t.Error(err)
			}

			tmplBuild, _, err := Render(string(tmpl), b.Steps[0].Name, b.Steps[0].Template.Name, b.Steps[0].Environment, b.Steps[0].Template.Variables)
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
	type args struct {
		velaFile string
	}

	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{"steps", args{velaFile: "testdata/build/basic/build.yml"}, "testdata/build/basic/want.yml", false},
		{"stages", args{velaFile: "testdata/build/basic_stages/build.yml"}, "testdata/build/basic_stages/want.yml", false},
		{"conditional match", args{velaFile: "testdata/build/conditional/build.yml"}, "testdata/build/conditional/want.yml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := os.ReadFile(tt.args.velaFile)
			if err != nil {
				t.Error(err)
			}

			got, _, err := RenderBuild("build", string(sFile), map[string]string{
				"VELA_REPO_FULL_NAME": "octocat/hello-world",
				"VELA_BUILD_BRANCH":   "main",
			}, map[string]interface{}{})
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderBuild() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true {
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
