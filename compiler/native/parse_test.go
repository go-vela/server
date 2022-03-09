// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"

	"github.com/urfave/cli/v2"
)

func TestNative_Parse_Metadata_Bytes(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Parse_Metadata_File(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	f, err := os.Open("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Opening file returned err: %v", err)
	}

	defer f.Close()

	got, err := client.Parse(f)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Metadata_Invalid(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))

	// run test
	got, err := client.Parse(nil)

	if err == nil {
		t.Error("Parse should have returned err")
	}

	if got != nil {
		t.Errorf("Parse is %v, want nil", got)
	}
}

func TestNative_Parse_Metadata_Path(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	got, err := client.Parse("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Metadata_Reader(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(bytes.NewReader(b))
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Metadata_String(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(string(b))
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Parameters(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "plugins/docker:18.09",
				Parameters: map[string]interface{}{
					"registry": "index.docker.io",
					"repo":     "github/octocat",
					"tags":     []interface{}{"latest", "dev"},
				},
				Name: "docker",
				Pull: "always",
				Secrets: yaml.StepSecretSlice{
					&yaml.StepSecret{
						Source: "docker_username",
						Target: "docker_username",
					},
					&yaml.StepSecret{
						Source: "docker_password",
						Target: "docker_password",
					},
				},
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/parameters.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_StagesPipeline(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
		Environment: map[string]string{
			"HELLO": "Hello, Global Environment",
		},
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name:  "install",
				Needs: raw.StringSlice{"clone"},
				Environment: map[string]string{
					"GRADLE_USER_HOME": ".gradle",
				},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew downloadDependencies"},
						Environment: map[string]string{
							"GRADLE_OPTS": "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						},
						Image: "openjdk:latest",
						Name:  "install",
						Pull:  "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "test",
				Needs: raw.StringSlice{"install", "clone"},
				Environment: map[string]string{
					"GRADLE_USER_HOME": "willBeOverwrittenInStep",
				},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew check"},
						Environment: map[string]string{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "test",
						Pull:  "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "build",
				Needs: raw.StringSlice{"install", "clone"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew build"},
						Environment: map[string]string{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "build",
						Pull:  "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "publish",
				Needs: raw.StringSlice{"build", "clone"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "plugins/docker:18.09",
						Parameters: map[string]interface{}{
							"registry": "index.docker.io",
							"repo":     "github/octocat",
							"tags":     []interface{}{"latest", "dev"},
						},
						Name: "publish",
						Pull: "always",
						Secrets: yaml.StepSecretSlice{
							&yaml.StepSecret{
								Source: "docker_username",
								Target: "registry_username",
							},
							&yaml.StepSecret{
								Source: "docker_password",
								Target: "registry_password",
							},
						},
					},
				},
			},
		},
		Secrets: yaml.SecretSlice{
			&yaml.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&yaml.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/stages_pipeline.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_StepsPipeline(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
		Environment: map[string]string{
			"HELLO": "Hello, Global Environment",
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: []string{"./gradlew downloadDependencies"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "install",
				Pull:  "always",
			},
			&yaml.Step{
				Commands: []string{"./gradlew check"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "test",
				Pull:  "always",
			},
			&yaml.Step{
				Commands: []string{"./gradlew build"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "build",
				Pull:  "always",
			},
			&yaml.Step{
				Image: "plugins/docker:18.09",
				Parameters: map[string]interface{}{
					"registry": "index.docker.io",
					"repo":     "github/octocat",
					"tags":     []interface{}{"latest", "dev"},
				},
				Name: "publish",
				Pull: "always",
				Secrets: yaml.StepSecretSlice{
					&yaml.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&yaml.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: yaml.SecretSlice{
			&yaml.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&yaml.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/steps_pipeline.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)
	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Parse_Secrets(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Secrets: yaml.SecretSlice{
			&yaml.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
			},
			&yaml.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
			},
			&yaml.Secret{
				Name:   "docker_username",
				Key:    "org/docker/username",
				Engine: "native",
				Type:   "org",
			},
			&yaml.Secret{
				Name:   "docker_password",
				Key:    "org/docker/password",
				Engine: "vault",
				Type:   "org",
			},
			&yaml.Secret{
				Name:   "docker_username",
				Key:    "org/team/docker/username",
				Engine: "native",
				Type:   "shared",
			},
			&yaml.Secret{
				Name:   "docker_password",
				Key:    "org/team/docker/password",
				Engine: "vault",
				Type:   "shared",
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/secrets.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)

	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Stages(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name:  "install",
				Needs: raw.StringSlice{"clone"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew downloadDependencies"},
						Environment: map[string]string{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "install",
						Pull:  "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "test",
				Needs: []string{"install", "clone"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew check"},
						Environment: map[string]string{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "test",
						Pull:  "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "build",
				Needs: []string{"install", "clone"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew build"},
						Environment: map[string]string{
							"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
							"GRADLE_USER_HOME": ".gradle",
						},
						Image: "openjdk:latest",
						Name:  "build",
						Pull:  "always",
					},
				},
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/stages.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)

	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_Parse_Steps(t *testing.T) {
	// setup types
	client, _ := New(cli.NewContext(nil, flag.NewFlagSet("test", 0), nil))
	want := &yaml.Build{
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: []string{"./gradlew downloadDependencies"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "install",
				Pull:  "always",
			},
			&yaml.Step{
				Commands: []string{"./gradlew check"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "test",
				Pull:  "always",
			},
			&yaml.Step{
				Commands: []string{"./gradlew build"},
				Environment: map[string]string{
					"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
					"GRADLE_USER_HOME": ".gradle",
				},
				Image: "openjdk:latest",
				Name:  "build",
				Pull:  "always",
			},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/steps.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := client.Parse(b)

	if err != nil {
		t.Errorf("Parse returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse is %v, want %v", got, want)
	}
}

func TestNative_ParseBytes_Metadata(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := ParseBytes(b)

	if err != nil {
		t.Errorf("ParseBytes returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseBytes is %v, want %v", got, want)
	}
}

func TestNative_ParseBytes_Invalid(t *testing.T) {
	// run test
	b, err := ioutil.ReadFile("testdata/invalid.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := ParseBytes(b)

	if err == nil {
		t.Error("ParseBytes should have returned err")
	}

	if got == new(yaml.Build) {
		t.Errorf("ParseBytes is %v, want %v", got, nil)
	}
}

func TestNative_ParseFile_Metadata(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	f, err := os.Open("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Opening file returned err: %v", err)
	}

	defer f.Close()

	got, err := ParseFile(f)

	if err != nil {
		t.Errorf("ParseFile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseFile is %v, want %v", got, want)
	}
}

func TestNative_ParseFile_Invalid(t *testing.T) {
	// run test
	f, err := os.Open("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Opening file returned err: %v", err)
	}

	f.Close()

	got, err := ParseFile(f)

	if err == nil {
		t.Error("ParseFile should have returned err")
	}

	if got != nil {
		t.Errorf("ParseFile is %v, want nil", got)
	}
}

func TestNative_ParsePath_Metadata(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	got, err := ParsePath("testdata/metadata.yml")

	if err != nil {
		t.Errorf("ParsePath returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParsePath is %v, want %v", got, want)
	}
}

func TestNative_ParsePath_Invalid(t *testing.T) {
	// run test
	got, err := ParsePath("testdata/foobar.yml")

	if err == nil {
		t.Error("ParsePath should have returned err")
	}

	if got != nil {
		t.Errorf("ParsePath is %v, want nil", got)
	}
}

func TestNative_ParseReader_Metadata(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := ParseReader(bytes.NewReader(b))

	if err != nil {
		t.Errorf("ParseReader returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseReader is %v, want %v", got, want)
	}
}

func TestNative_ParseReader_Invalid(t *testing.T) {
	// run test
	got, err := ParseReader(FailReader{})

	if err == nil {
		t.Error("ParseFile should have returned err")
	}

	if got != nil {
		t.Errorf("ParseFile is %v, want nil", got)
	}
}

func TestNative_ParseString_Metadata(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: []string{"steps", "services", "secrets"},
		},
	}

	// run test
	b, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := ParseString(string(b))

	if err != nil {
		t.Errorf("ParseString returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseString is %v, want %v", got, want)
	}
}

type FailReader struct{}

func (FailReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("this is a reader that fails when you try to read")
}

func Test_client_Parse(t *testing.T) {
	// setup types
	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Clone:       nil,
			Environment: nil,
		},
		Steps: yaml.StepSlice{
			{
				Name:  "foo",
				Image: "alpine",
				Pull:  "not_present",
				Parameters: map[string]interface{}{
					"registry": "foo",
				},
			},
		},
	}

	type args struct {
		pipelineType string
		file         string
	}

	tests := []struct {
		name    string
		args    args
		want    *yaml.Build
		wantErr bool
	}{
		{"yaml", args{pipelineType: constants.PipelineTypeYAML, file: "testdata/pipeline_type_default.yml"}, want, false},
		{"starlark", args{pipelineType: constants.PipelineTypeStarlark, file: "testdata/pipeline_type.star"}, want, false},
		{"go", args{pipelineType: constants.PipelineTypeGo, file: "testdata/pipeline_type_go.yml"}, want, false},
		{"empty", args{pipelineType: "", file: "testdata/pipeline_type_default.yml"}, want, false},
		{"nil", args{pipelineType: "nil", file: "testdata/pipeline_type_default.yml"}, want, false},
		{"invalid", args{pipelineType: "foo", file: "testdata/pipeline_type_default.yml"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ioutil.ReadFile(tt.args.file)
			if err != nil {
				t.Errorf("Reading file returned err: %v", err)
			}

			var c *client
			if tt.args.pipelineType == "nil" {
				c = &client{}
			} else {
				c = &client{
					repo: &library.Repo{PipelineType: &tt.args.pipelineType},
				}
			}

			got, err := c.Parse(content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_client_ParseRaw(t *testing.T) {
	expected, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	type args struct {
		kind string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"byte", args{kind: "byte"}, string(expected), false},
		{"file", args{kind: "file"}, string(expected), false},
		{"io reader", args{kind: "ioreader"}, string(expected), false},
		{"string", args{kind: "string"}, string(expected), false},
		{"path", args{kind: "path"}, string(expected), false},
		{"unexpected", args{kind: "foo"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content interface{}
			var err error
			switch tt.args.kind {
			case "byte":
				content, err = ioutil.ReadFile("testdata/metadata.yml")
				if err != nil {
					t.Errorf("Reading file returned err: %v", err)
				}
			case "file":
				content, err = os.Open("testdata/metadata.yml")
				if err != nil {
					t.Errorf("Reading file returned err: %v", err)
				}
			case "ioreader":
				b, err := ioutil.ReadFile("testdata/metadata.yml")
				if err != nil {
					t.Errorf("ParseReader returned err: %v", err)
				}

				content = bytes.NewReader(b)
				if err != nil {
					t.Errorf("Reading file returned err: %v", err)
				}
			case "path":
				content = "testdata/metadata.yml"
			case "string":
				content = tt.want
			}

			c := &client{}
			got, err := c.ParseRaw(content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseRaw() got = %v, want %v", got, tt.want)
			}
		})
	}
}
