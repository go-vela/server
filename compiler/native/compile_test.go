// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"

	"github.com/go-vela/types"
	"github.com/go-vela/types/pipeline"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"

	yml "github.com/buildkite/yaml"
)

func TestNative_Compile_StagesPipeline(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	initEnv := environment(nil, m, nil, nil)
	initEnv["HELLO"] = "Hello, Global Environment"

	cloneEnv := environment(nil, m, nil, nil)
	cloneEnv["HELLO"] = "Hello, Global Environment"

	installEnv := environment(nil, m, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})
	installEnv["HELLO"] = "Hello, Global Environment"

	testEnv := environment(nil, m, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})
	testEnv["HELLO"] = "Hello, Global Environment"

	buildEnv := environment(nil, m, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})
	buildEnv["HELLO"] = "Hello, Global Environment"

	dockerEnv := environment(nil, m, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"
	dockerEnv["HELLO"] = "Hello, Global Environment"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_init_init",
						Directory:   "/vela/src/foo//",
						Environment: initEnv,
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
			&pipeline.Stage{
				Name: "clone",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_clone_clone",
						Directory:   "/vela/src/foo//",
						Environment: cloneEnv,
						Image:       "target/vela-git:v0.4.0",
						Name:        "clone",
						Number:      2,
						Pull:        "not_present",
					},
				},
			},
			&pipeline.Stage{
				Name:  "install",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_install_install",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: installEnv,
						Image:       "openjdk:latest",
						Name:        "install",
						Number:      3,
						Pull:        "always",
					},
				},
			},
			&pipeline.Stage{
				Name:  "test",
				Needs: []string{"install", "clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_test_test",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: testEnv,
						Image:       "openjdk:latest",
						Name:        "test",
						Number:      4,
						Pull:        "always",
					},
				},
			},
			&pipeline.Stage{
				Name:  "build",
				Needs: []string{"install", "clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_build_build",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: buildEnv,
						Image:       "openjdk:latest",
						Name:        "build",
						Number:      5,
						Pull:        "always",
					},
				},
			},
			&pipeline.Stage{
				Name:  "publish",
				Needs: []string{"build", "clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_publish_publish",
						Directory:   "/vela/src/foo//",
						Image:       "plugins/docker:18.09",
						Environment: dockerEnv,
						Name:        "publish",
						Number:      6,
						Pull:        "always",
						Secrets: pipeline.StepSecretSlice{
							&pipeline.StepSecret{
								Source: "docker_username",
								Target: "registry_username",
							},
							&pipeline.StepSecret{
								Source: "docker_password",
								Target: "registry_password",
							},
						},
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/stages_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	compiler.WithMetadata(m)

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	// WARNING: hack to compare stages
	//
	// Channel values can only be compared for equality.
	// Two channel values are considered equal if they
	// originated from the same make call meaning they
	// refer to the same channel value in memory.
	for i, stage := range got.Stages {
		tmp := want.Stages

		tmp[i].Done = stage.Done
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Compile_StagesPipeline_Modification(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	engine.POST("/config/bad", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"foo": "bar"})
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	name := "foo"
	author := "author"
	number := 1

	// run test
	yaml, err := ioutil.ReadFile("testdata/stages_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	type args struct {
		endpoint     string
		libraryBuild *library.Build
		repo         *library.Repo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bad url", args{
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     "bad",
		}, true},
		{"invalid return", args{
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/bad"),
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := client{
				ModificationService: ModificationConfig{
					Timeout:  1 * time.Second,
					Endpoint: tt.args.endpoint,
				},
				repo:  &library.Repo{Name: &author},
				build: &library.Build{Author: &name, Number: &number},
			}
			_, err := compiler.Compile(yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNative_Compile_StepsPipeline_Modification(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	engine.POST("/config/bad", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, gin.H{"foo": "bar"})
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	name := "foo"
	author := "author"
	number := 1

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	type args struct {
		endpoint     string
		libraryBuild *library.Build
		repo         *library.Repo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"bad url", args{
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     "bad",
		}, true},
		{"invalid return", args{
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/bad"),
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := client{
				ModificationService: ModificationConfig{
					Timeout:  1 * time.Second,
					Endpoint: tt.args.endpoint,
				},
				repo:  tt.args.repo,
				build: tt.args.libraryBuild,
			}
			_, err := compiler.Compile(yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNative_Compile_StepsPipeline(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	initEnv := environment(nil, m, nil, nil)
	initEnv["HELLO"] = "Hello, Global Environment"

	cloneEnv := environment(nil, m, nil, nil)
	cloneEnv["HELLO"] = "Hello, Global Environment"

	installEnv := environment(nil, m, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})
	installEnv["HELLO"] = "Hello, Global Environment"

	testEnv := environment(nil, m, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})
	testEnv["HELLO"] = "Hello, Global Environment"

	buildEnv := environment(nil, m, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})
	buildEnv["HELLO"] = "Hello, Global Environment"

	dockerEnv := environment(nil, m, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"
	dockerEnv["HELLO"] = "Hello, Global Environment"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: initEnv,
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: cloneEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_install",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/vela/src/foo//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Number:      3,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_test",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/vela/src/foo//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Number:      4,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_build",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Directory:   "/vela/src/foo//",
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: buildEnv,
				Image:       "openjdk:latest",
				Name:        "build",
				Number:      5,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_publish",
				Directory:   "/vela/src/foo//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "publish",
				Number:      6,
				Pull:        "always",
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_pipeline.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	compiler.WithMetadata(m)

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Compile mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Compile_StagesPipelineTemplate(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:name/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	setupEnv := environment(nil, m, nil, nil)
	setupEnv["bar"] = "test4"
	setupEnv["star"] = "test3"

	installEnv := environment(nil, m, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})
	installEnv["bar"] = "test4"
	installEnv["star"] = "test3"

	testEnv := environment(nil, m, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})
	testEnv["bar"] = "test4"
	testEnv["star"] = "test3"

	buildEnv := environment(nil, m, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})
	buildEnv["bar"] = "test4"
	buildEnv["star"] = "test3"

	dockerEnv := environment(nil, m, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"
	dockerEnv["bar"] = "test4"
	dockerEnv["star"] = "test3"

	serviceEnv := environment(nil, m, nil, nil)
	serviceEnv["bar"] = "test4"
	serviceEnv["star"] = "test3"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: "init",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_init_init",
						Directory:   "/vela/src/foo//",
						Environment: setupEnv,
						Image:       "#init",
						Name:        "init",
						Number:      1,
						Pull:        "not_present",
					},
				},
			},
			&pipeline.Stage{
				Name: "clone",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_clone_clone",
						Directory:   "/vela/src/foo//",
						Environment: setupEnv,
						Image:       "target/vela-git:v0.4.0",
						Name:        "clone",
						Number:      2,
						Pull:        "not_present",
					},
				},
			},
			&pipeline.Stage{
				Name:  "gradle",
				Needs: []string{"clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_gradle_sample_install",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: installEnv,
						Image:       "openjdk:latest",
						Name:        "sample_install",
						Number:      3,
						Pull:        "always",
					},
					&pipeline.Container{
						ID:          "__0_gradle_sample_test",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: testEnv,
						Image:       "openjdk:latest",
						Name:        "sample_test",
						Number:      4,
						Pull:        "always",
					},
					&pipeline.Container{
						ID:          "__0_gradle_sample_build",
						Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
						Directory:   "/vela/src/foo//",
						Entrypoint:  []string{"/bin/sh", "-c"},
						Environment: buildEnv,
						Image:       "openjdk:latest",
						Name:        "sample_build",
						Number:      5,
						Pull:        "always",
					},
				},
			},
			&pipeline.Stage{
				Name:  "publish",
				Needs: []string{"gradle", "clone"},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "__0_publish_publish",
						Directory:   "/vela/src/foo//",
						Image:       "plugins/docker:18.09",
						Environment: dockerEnv,
						Name:        "publish",
						Number:      6,
						Pull:        "always",
						Secrets: pipeline.StepSecretSlice{
							&pipeline.StepSecret{
								Source: "docker_username",
								Target: "registry_username",
							},
							&pipeline.StepSecret{
								Source: "docker_password",
								Target: "registry_password",
							},
						},
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "foo_password",
				Key:    "org/repo/foo/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service___0_postgres",
				Detach:      true,
				Image:       "postgres:12",
				Name:        "postgres",
				Number:      1,
				Pull:        "not_present",
				Environment: serviceEnv,
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/stages_pipeline_template.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	compiler.WithMetadata(m)

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	// WARNING: hack to compare stages
	//
	// Channel values can only be compared for equality.
	// Two channel values are considered equal if they
	// originated from the same make call meaning they
	// refer to the same channel value in memory.
	for i, stage := range got.Stages {
		tmp := want.Stages

		tmp[i].Done = stage.Done
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Compile_StepsPipelineTemplate(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	setupEnv := environment(nil, m, nil, nil)
	setupEnv["bar"] = "test4"
	setupEnv["star"] = "test3"

	installEnv := environment(nil, m, nil, nil)
	installEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	installEnv["GRADLE_USER_HOME"] = ".gradle"
	installEnv["HOME"] = "/root"
	installEnv["SHELL"] = "/bin/sh"
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})
	installEnv["bar"] = "test4"
	installEnv["star"] = "test3"

	testEnv := environment(nil, m, nil, nil)
	testEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	testEnv["GRADLE_USER_HOME"] = ".gradle"
	testEnv["HOME"] = "/root"
	testEnv["SHELL"] = "/bin/sh"
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})
	testEnv["bar"] = "test4"
	testEnv["star"] = "test3"

	buildEnv := environment(nil, m, nil, nil)
	buildEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	buildEnv["GRADLE_USER_HOME"] = ".gradle"
	buildEnv["HOME"] = "/root"
	buildEnv["SHELL"] = "/bin/sh"
	buildEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew build"})
	buildEnv["bar"] = "test4"
	buildEnv["star"] = "test3"

	dockerEnv := environment(nil, m, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"
	dockerEnv["bar"] = "test4"
	dockerEnv["star"] = "test3"

	serviceEnv := environment(nil, m, nil, nil)
	serviceEnv["bar"] = "test4"
	serviceEnv["star"] = "test3"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: setupEnv,
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: setupEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_sample_install",
				Directory:   "/vela/src/foo//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "sample_install",
				Number:      3,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_sample_test",
				Directory:   "/vela/src/foo//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "sample_test",
				Number:      4,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_sample_build",
				Directory:   "/vela/src/foo//",
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{"/bin/sh", "-c"},
				Environment: buildEnv,
				Image:       "openjdk:latest",
				Name:        "sample_build",
				Number:      5,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_docker",
				Directory:   "/vela/src/foo//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "docker",
				Number:      6,
				Pull:        "always",
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "foo_password",
				Key:    "org/repo/foo/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "service___0_postgres",
				Detach:      true,
				Environment: serviceEnv,
				Image:       "postgres:12",
				Name:        "postgres",
				Number:      1,
				Pull:        "not_present",
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_pipeline_template.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	compiler.WithMetadata(m)

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_Compile_InvalidType(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", s.URL, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	gradleEnv := environment(nil, m, nil, nil)
	gradleEnv["GRADLE_OPTS"] = "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false"
	gradleEnv["GRADLE_USER_HOME"] = ".gradle"

	dockerEnv := environment(nil, m, nil, nil)
	dockerEnv["PARAMETER_REGISTRY"] = "index.docker.io"
	dockerEnv["PARAMETER_REPO"] = "github/octocat"
	dockerEnv["PARAMETER_TAGS"] = "latest,dev"

	want := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_docker",
				Directory:   "/vela/src/foo//",
				Image:       "plugins/docker:18.09",
				Environment: dockerEnv,
				Name:        "docker",
				Number:      3,
				Pull:        "always",
				Secrets: pipeline.StepSecretSlice{
					&pipeline.StepSecret{
						Source: "docker_username",
						Target: "registry_username",
					},
					&pipeline.StepSecret{
						Source: "docker_password",
						Target: "registry_password",
					},
				},
			},
		},
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Name:   "docker_username",
				Key:    "org/repo/docker/username",
				Engine: "native",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
			&pipeline.Secret{
				Name:   "docker_password",
				Key:    "org/repo/docker/password",
				Engine: "vault",
				Type:   "repo",
				Origin: &pipeline.Container{},
			},
		},
	}

	// run test
	yaml, err := ioutil.ReadFile("testdata/invalid_type.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	compiler.WithMetadata(m)

	got, err := compiler.Compile(yaml)
	if err != nil {
		t.Errorf("Compile returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compile is %v, want %v", got, want)
	}
}

func TestNative_Compile_Clone(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	fooEnv := environment(nil, m, nil, nil)
	fooEnv["PARAMETER_REGISTRY"] = "foo"

	cloneEnv := environment(nil, m, nil, nil)
	cloneEnv["PARAMETER_DEPTH"] = "5"

	wantFalse := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       false,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: fooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      2,
				Pull:        "always",
			},
		},
	}

	wantTrue := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       true,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: fooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      3,
				Pull:        "always",
			},
		},
	}

	wantReplace := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:       false,
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: cloneEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "always",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: fooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      3,
				Pull:        "always",
			},
		},
	}

	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *pipeline.Build
		wantErr bool
	}{
		{"false", args{
			file: "testdata/clone_false.yml",
		}, wantFalse, false},
		{"true", args{
			file: "testdata/clone_true.yml",
		}, wantTrue, false},
		{"replace", args{
			file: "testdata/clone_replace.yml",
		}, wantReplace, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// run test
			yaml, err := ioutil.ReadFile(tt.args.file)
			if err != nil {
				t.Errorf("Reading yaml file return err: %v", err)
			}

			compiler, err := New(c)
			if err != nil {
				t.Errorf("Creating compiler returned err: %v", err)
			}

			compiler.WithMetadata(m)

			got, err := compiler.Compile(yaml)
			if err != nil {
				t.Errorf("Compile returned err: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNative_Compile_Pipeline_Type(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-token", "", "doc")
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	defaultFooEnv := environment(nil, m, nil, nil)
	defaultFooEnv["PARAMETER_REGISTRY"] = "foo"

	defaultEnv := environment(nil, m, nil, nil)
	wantDefault := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:    true,
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: defaultEnv,
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: defaultEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: defaultFooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      3,
				Pull:        "not_present",
			},
		},
	}

	goPipelineType := "go"

	goFooEnv := environment(nil, m, &library.Repo{PipelineType: &goPipelineType}, nil)
	goFooEnv["PARAMETER_REGISTRY"] = "foo"

	defaultGoEnv := environment(nil, m, &library.Repo{PipelineType: &goPipelineType}, nil)
	wantGo := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:    true,
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: defaultGoEnv,
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: defaultGoEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: goFooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      3,
				Pull:        "not_present",
			},
		},
	}

	starPipelineType := "starlark"

	starlarkFooEnv := environment(nil, m, &library.Repo{PipelineType: &starPipelineType}, nil)
	starlarkFooEnv["PARAMETER_REGISTRY"] = "foo"

	defaultStarlarkEnv := environment(nil, m, &library.Repo{PipelineType: &starPipelineType}, nil)
	wantStarlark := &pipeline.Build{
		Version: "1",
		ID:      "__0",
		Metadata: pipeline.Metadata{
			Clone:    true,
			Template: false,
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:          "step___0_init",
				Directory:   "/vela/src/foo//",
				Environment: defaultStarlarkEnv,
				Image:       "#init",
				Name:        "init",
				Number:      1,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_clone",
				Directory:   "/vela/src/foo//",
				Environment: defaultStarlarkEnv,
				Image:       "target/vela-git:v0.4.0",
				Name:        "clone",
				Number:      2,
				Pull:        "not_present",
			},
			&pipeline.Container{
				ID:          "step___0_foo",
				Directory:   "/vela/src/foo//",
				Environment: starlarkFooEnv,
				Image:       "alpine",
				Name:        "foo",
				Number:      3,
				Pull:        "not_present",
			},
		},
	}

	type args struct {
		file         string
		pipelineType string
	}
	tests := []struct {
		name    string
		args    args
		want    *pipeline.Build
		wantErr bool
	}{
		{"default", args{file: "testdata/pipeline_type_default.yml", pipelineType: ""}, wantDefault, false},
		{"golang", args{file: "testdata/pipeline_type_go.yml", pipelineType: "go"}, wantGo, false},
		{"starlark", args{file: "testdata/pipeline_type.star", pipelineType: "starlark"}, wantStarlark, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// run test
			yaml, err := ioutil.ReadFile(tt.args.file)
			if err != nil {
				t.Errorf("Reading yaml file return err: %v", err)
			}

			compiler, err := New(c)
			if err != nil {
				t.Errorf("Creating compiler returned err: %v", err)
			}

			compiler.WithMetadata(m)
			compiler.WithRepo(&library.Repo{PipelineType: &tt.args.pipelineType})

			got, err := compiler.Compile(yaml)
			if err != nil {
				t.Errorf("Compile returned err: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Compile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNative_Compile_NoStepsorStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	name := "foo"
	author := "author"
	number := 1

	// run test
	yaml, err := ioutil.ReadFile("testdata/metadata.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}
	compiler.repo = &library.Repo{Name: &author}
	compiler.build = &library.Build{Author: &name, Number: &number}

	got, err := compiler.Compile(yaml)
	if err == nil {
		t.Errorf("Compile should have returned err")
	}

	if got != nil {
		t.Errorf("Compile is %v, want %v", got, nil)
	}
}

func TestNative_Compile_StepsandStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	name := "foo"
	author := "author"
	number := 1

	// run test
	yaml, err := ioutil.ReadFile("testdata/steps_and_stages.yml")
	if err != nil {
		t.Errorf("Reading yaml file return err: %v", err)
	}

	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}
	compiler.repo = &library.Repo{Name: &author}
	compiler.build = &library.Build{Author: &name, Number: &number}

	got, err := compiler.Compile(yaml)
	if err == nil {
		t.Errorf("Compile should have returned err")
	}

	if got != nil {
		t.Errorf("Compile is %v, want %v", got, nil)
	}
}

// convertResponse converts the build to the ModifyResponse.
func convertResponse(build *yaml.Build) (*ModifyResponse, error) {
	data, err := yml.Marshal(build)
	if err != nil {
		return nil, err
	}

	response := &ModifyResponse{
		Pipeline: string(data),
	}

	return response, nil
}

func Test_client_modifyConfig(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/foo/bar/contents/:path", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
		c.File("testdata/template.json")
	})

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	want := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Pull:        "not_present",
			},
			&yaml.Step{
				Environment: environment(nil, m, nil, nil),
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Pull:        "not_present",
			},
			&yaml.Step{
				Image:       "plugins/docker:18.09",
				Environment: nil,
				Name:        "docker",
				Pull:        "always",
				Parameters: map[string]interface{}{
					"init_options": map[interface{}]interface{}{
						"get_plugins": "true",
					},
				},
			},
		},
	}

	want2 := &yaml.Build{
		Version: "1",
		Metadata: yaml.Metadata{
			Template:    false,
			Environment: []string{"steps", "services", "secrets"},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Environment: environment(nil, m, nil, nil),
				Image:       "#init",
				Name:        "init",
				Pull:        "not_present",
			},
			&yaml.Step{
				Environment: environment(nil, m, nil, nil),
				Image:       "target/vela-git:v0.3.0",
				Name:        "clone",
				Pull:        "not_present",
			},
			&yaml.Step{
				Image:       "plugins/docker:18.09",
				Environment: nil,
				Name:        "docker",
				Pull:        "always",
				Parameters: map[string]interface{}{
					"init_options": map[interface{}]interface{}{
						"get_plugins": "true",
					},
				},
			},
			&yaml.Step{
				Image:       "alpine",
				Environment: nil,
				Name:        "modification",
				Pull:        "always",
				Commands:    []string{"echo hello from modification"},
			},
		},
	}

	engine.POST("/config/unmodified", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		response, err := convertResponse(want)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	engine.POST("/config/timeout", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.Header("Content-Type", "application/json")
		response, err := convertResponse(want)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	engine.POST("/config/modified", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		output := want
		var steps []*yaml.Step
		steps = append(steps, want.Steps...)
		steps = append(steps, &yaml.Step{
			Image:       "alpine",
			Environment: nil,
			Name:        "modification",
			Pull:        "always",
			Commands:    []string{"echo hello from modification"},
		})
		output.Steps = steps
		response, err := convertResponse(want)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	engine.POST("/config/empty", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	engine.POST("/config/unauthorized", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		response, err := convertResponse(want)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusForbidden, response)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	name := "foo"
	author := "author"
	number := 1

	type args struct {
		endpoint     string
		build        *yaml.Build
		libraryBuild *library.Build
		repo         *library.Repo
	}
	tests := []struct {
		name    string
		args    args
		want    *yaml.Build
		wantErr bool
	}{
		{"unmodified", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/unmodified"),
		}, want, false},
		{"modified", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/modified"),
		}, want2, false},
		{"invalid endpoint", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     "bad",
		}, nil, true},
		{"unauthorized endpoint", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/unauthorized"),
		}, nil, true},
		{"timeout endpoint", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/timeout"),
		}, nil, true},
		{"empty payload", args{
			build:        want,
			libraryBuild: &library.Build{Number: &number, Author: &author},
			repo:         &library.Repo{Name: &name},
			endpoint:     fmt.Sprintf("%s/%s", s.URL, "config/empty"),
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := client{
				ModificationService: ModificationConfig{
					Timeout:  2 * time.Second,
					Retries:  2,
					Endpoint: tt.args.endpoint,
				},
			}
			got, err := compiler.modifyConfig(tt.args.build, tt.args.libraryBuild, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("modifyConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("modifyConfig() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
