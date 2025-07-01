// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"os"
	"reflect"
	"testing"

	"go.yaml.in/yaml/v3"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Origin_MergeEnv(t *testing.T) {
	// setup tests
	tests := []struct {
		origin      *Origin
		environment map[string]string
		failure     bool
	}{
		{
			origin: &Origin{
				Name:        "vault",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-vault:latest",
				Parameters: map[string]interface{}{
					"addr":        "vault.example.com",
					"auth_method": "token",
					"items": []interface{}{
						map[string]string{"source": "secret/docker", "path": "docker"},
					},
				},
				Pull: "always",
				Secrets: StepSecretSlice{
					{
						Source: "vault_token",
						Target: "vault_token",
					},
				},
			},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			origin:      &Origin{},
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			origin:      nil,
			environment: map[string]string{"BAR": "baz"},
			failure:     false,
		},
		{
			origin: &Origin{
				Name:        "vault",
				Environment: map[string]string{"FOO": "bar"},
				Image:       "target/vela-vault:latest",
				Parameters: map[string]interface{}{
					"addr":        "vault.example.com",
					"auth_method": "token",
					"items": []interface{}{
						map[string]string{"source": "secret/docker", "path": "docker"},
					},
				},
				Pull: "always",
				Secrets: StepSecretSlice{
					{
						Source: "vault_token",
						Target: "vault_token",
					},
				},
			},
			environment: nil,
			failure:     true,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.origin.MergeEnv(test.environment)

		if test.failure {
			if err == nil {
				t.Errorf("MergeEnv should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("MergeEnv returned err: %v", err)
		}
	}
}

func TestYaml_SecretSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		secrets *SecretSlice
		want    *pipeline.SecretSlice
	}{
		{
			secrets: &SecretSlice{
				{
					Name:   "docker_username",
					Key:    "github/octocat/docker/username",
					Engine: "native",
					Type:   "repo",
					Origin: Origin{},
					Pull:   "build_start",
				},
				{
					Name:   "docker_username",
					Key:    "",
					Engine: "",
					Type:   "",
					Origin: Origin{
						Name:        "vault",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-vault:latest",
						Parameters: map[string]interface{}{
							"addr": "vault.company.com",
						},
						Pull: "always",
						Ruleset: Ruleset{
							If: Rules{
								Event:    []string{"push"},
								Operator: "and",
							},
						},
						Secrets: StepSecretSlice{
							{
								Source: "foo",
								Target: "foo",
							},
							{
								Source: "foobar",
								Target: "foobar",
							},
						},
					},
					Pull: "build_start",
				},
			},
			want: &pipeline.SecretSlice{
				{
					Name:   "docker_username",
					Key:    "github/octocat/docker/username",
					Engine: "native",
					Type:   "repo",
					Origin: &pipeline.Container{},
					Pull:   "build_start",
				},
				{
					Name:   "docker_username",
					Key:    "",
					Engine: "",
					Type:   "",
					Origin: &pipeline.Container{
						Name:        "vault",
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-vault:latest",
						Pull:        "always",
						Ruleset: pipeline.Ruleset{
							If: pipeline.Rules{
								Event:    []string{"push"},
								Operator: "and",
							},
						},
						Secrets: pipeline.StepSecretSlice{
							{
								Source: "foo",
								Target: "foo",
							},
							{
								Source: "foobar",
								Target: "foobar",
							},
						},
					},
					Pull: "build_start",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.secrets.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_SecretSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *SecretSlice
	}{
		{
			failure: false,
			file:    "testdata/secret.yml",
			want: &SecretSlice{
				{
					Name:   "foo",
					Key:    "bar",
					Engine: "native",
					Type:   "repo",
					Pull:   "build_start",
				},
				{
					Name:   "noKey",
					Key:    "noKey",
					Engine: "native",
					Type:   "repo",
					Pull:   "build_start",
				},
				{
					Name:   "noType",
					Key:    "bar",
					Engine: "native",
					Type:   "repo",
					Pull:   "build_start",
				},
				{
					Name:   "noEngine",
					Key:    "bar",
					Engine: "native",
					Type:   "repo",
					Pull:   "build_start",
				},
				{
					Name:   "noKeyEngineAndType",
					Key:    "noKeyEngineAndType",
					Engine: "native",
					Type:   "repo",
					Pull:   "build_start",
				},
				{
					Name:   "externalSecret",
					Key:    "",
					Engine: "",
					Type:   "",
					Origin: Origin{
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-vault:latest",
						Parameters: map[string]interface{}{
							"addr": "vault.company.com",
						},
						Pull: "always",
						Ruleset: Ruleset{
							If: Rules{
								Event:    []string{"push"},
								Operator: "and",
								Matcher:  "filepath",
							},
						},
						Secrets: StepSecretSlice{
							{
								Source: "foo",
								Target: "FOO",
							},
							{
								Source: "foobar",
								Target: "FOOBAR",
							},
						},
					},
					Pull: "",
				},
				{
					Name:   "",
					Key:    "",
					Engine: "",
					Type:   "",
					Origin: Origin{
						Environment: map[string]string{"FOO": "bar"},
						Image:       "target/vela-vault:latest",
						Parameters: map[string]interface{}{
							"addr": "vault.company.com",
						},
						Pull: "always",
						Ruleset: Ruleset{
							If: Rules{
								Event:    []string{"push"},
								Operator: "and",
								Matcher:  "filepath",
							},
						},
						Secrets: StepSecretSlice{
							{
								Source: "foo",
								Target: "FOO",
							},
							{
								Source: "foobar",
								Target: "FOOBAR",
							},
						},
					},
					Pull: "",
				},
			},
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(SecretSlice)

		// run test
		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
		}

		err = yaml.Unmarshal(b, got)

		if test.failure {
			if err == nil {
				t.Errorf("UnmarshalYAML should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalYAML is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_StepSecretSlice_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		secrets *StepSecretSlice
		want    *pipeline.StepSecretSlice
	}{
		{
			secrets: &StepSecretSlice{
				{
					Source: "docker_username",
					Target: "plugin_username",
				},
			},
			want: &pipeline.StepSecretSlice{
				{
					Source: "docker_username",
					Target: "plugin_username",
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.secrets.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}

func TestYaml_StepSecretSlice_UnmarshalYAML(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		file    string
		want    *StepSecretSlice
	}{
		{
			failure: false,
			file:    "testdata/step_secret_slice.yml",
			want: &StepSecretSlice{
				{
					Source: "foo",
					Target: "BAR",
				},
				{
					Source: "hello",
					Target: "WORLD",
				},
			},
		},
		{
			failure: false,
			file:    "testdata/step_secret_string.yml",
			want: &StepSecretSlice{
				{
					Source: "foo",
					Target: "FOO",
				},
				{
					Source: "hello",
					Target: "HELLO",
				},
			},
		},
		{
			failure: true,
			file:    "testdata/step_secret_slice_invalid_no_source.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/step_secret_slice_invalid_no_target.yml",
			want:    nil,
		},
		{
			failure: true,
			file:    "testdata/invalid.yml",
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := new(StepSecretSlice)

		// run test
		b, err := os.ReadFile(test.file)
		if err != nil {
			t.Errorf("unable to read file: %v", err)
		}

		err = yaml.Unmarshal(b, got)

		if test.failure {
			if err == nil {
				t.Errorf("UnmarshalYAML should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UnmarshalYAML returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("UnmarshalYAML is %v, want %v", got, test.want)
		}
	}
}
