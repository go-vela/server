// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

func TestNative_ExpandStages(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	tmpls := map[string]*yaml.Template{
		"gradle": {
			Name:   "gradle",
			Source: "github.example.com/foo/bar/long_template.yml",
			Type:   "github",
		},
	}

	stages := yaml.StageSlice{
		&yaml.Stage{
			Name: "foo",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Name: "sample",
					Template: yaml.StepTemplate{
						Name: "gradle",
						Variables: map[string]interface{}{
							"image":       "openjdk:latest",
							"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
							"pull_policy": "pull: true",
						},
					},
				},
			},
		},
	}

	wantStages := yaml.StageSlice{
		&yaml.Stage{
			Name: "foo",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Commands: []string{"./gradlew downloadDependencies"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_install",
					Pull:  "always",
				},
				&yaml.Step{
					Commands: []string{"./gradlew check"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_test",
					Pull:  "always",
				},
				&yaml.Step{
					Commands: []string{"./gradlew build", "echo gradle"},
					Environment: raw.StringSliceMap{
						"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
						"GRADLE_USER_HOME": ".gradle",
					},
					Image: "openjdk:latest",
					Name:  "sample_build",
					Pull:  "always",
				},
			},
		},
	}

	wantSecrets := yaml.SecretSlice{
		&yaml.Secret{
			Name:   "docker_username",
			Key:    "org/repo/foo/bar",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "foo_password",
			Key:    "org/repo/foo/password",
			Engine: "vault",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
	}

	wantServices := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres:12",
			Name:  "postgres",
			Pull:  "not_present",
		},
	}

	wantEnvironment := raw.StringSliceMap{
		"star": "test3",
		"bar":  "test4",
	}

	// run test -- missing private github
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.PrivateGithub = nil
	_, _, err = compiler.ExpandStages(
		context.Background(),
		&yaml.Build{
			Stages:      stages,
			Services:    yaml.ServiceSlice{},
			Environment: raw.StringSliceMap{},
		},
		tmpls,
		new(pipeline.RuleData),
		nil,
	)

	if err == nil {
		t.Errorf("ExpandStages should have returned error with empty PrivateGitHub")
	}

	// run test
	compiler, err = FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	build, _, err := compiler.ExpandStages(
		context.Background(),
		&yaml.Build{
			Stages:      stages,
			Services:    yaml.ServiceSlice{},
			Environment: raw.StringSliceMap{},
		},
		tmpls,
		new(pipeline.RuleData),
		nil,
	)

	if err != nil {
		t.Errorf("ExpandStages returned err: %v", err)
	}

	if diff := cmp.Diff(build.Stages, wantStages); diff != "" {
		t.Errorf("ExpandStages() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
		t.Errorf("ExpandStages() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Services, wantServices); diff != "" {
		t.Errorf("ExpandStages() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
		t.Errorf("ExpandStages() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_ExpandSteps(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "GitHub",
			tmpls: map[string]*yaml.Template{
				"gradle": {
					Name:   "gradle",
					Source: "github.example.com/foo/bar/long_template.yml",
					Type:   "github",
				},
			},
		},
		{
			name: "File",
			tmpls: map[string]*yaml.Template{
				"gradle": {
					Name:   "gradle",
					Source: "long_template.yml",
					Type:   "file",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "gradle",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	wantSteps := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build", "echo gradle"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_build",
			Pull:  "always",
		},
	}

	wantSecrets := yaml.SecretSlice{
		&yaml.Secret{
			Name:   "docker_username",
			Key:    "org/repo/foo/bar",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "foo_password",
			Key:    "org/repo/foo/password",
			Engine: "vault",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
	}

	wantServices := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres:12",
			Name:  "postgres",
			Pull:  "not_present",
		},
	}

	wantEnvironment := raw.StringSliceMap{
		"foo":  "test1",
		"bar":  "test2",
		"star": "test3",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithCommit("123abc456def").WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			build, _, err := compiler.ExpandSteps(
				context.Background(),
				&yaml.Build{
					Steps:       steps,
					Services:    yaml.ServiceSlice{},
					Environment: globalEnvironment,
				},
				test.tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
			if err != nil {
				t.Errorf("ExpandSteps_Type%s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(build.Steps, wantSteps); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Services, wantServices); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}

func TestNative_ExpandStepsWarnings(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "warnings",
			tmpls: map[string]*yaml.Template{
				"warnings": {
					Name:   "steps_merge_anchor_1.yml",
					Source: "github.example.com/foo/bar/steps_merge_anchor_1.yml",
					Type:   "github",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "warnings",
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	wantWarnings := []string{
		"[warnings]:25:duplicate << keys in single YAML map",
		"[warnings]:32:duplicate << keys in single YAML map",
		"[warnings]:44:duplicate << keys in single YAML map",
		"[warnings]:43:duplicate << keys in single YAML map",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithCommit("123abc456def").WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, warnings, err := compiler.ExpandSteps(
				context.Background(),
				&yaml.Build{
					Steps:       steps,
					Services:    yaml.ServiceSlice{},
					Environment: globalEnvironment,
				},
				test.tmpls, new(pipeline.RuleData), []string{}, compiler.GetTemplateDepth())
			if err != nil {
				t.Errorf("ExpandSteps_Type%s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(warnings, wantWarnings) {
				t.Errorf("ExpandSteps()_Type%s returned incorrect warnings: %v", test.name, warnings)
			}
		})
	}
}

func TestNative_ExpandDeployment(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "GitHub",
			tmpls: map[string]*yaml.Template{
				"deploy": {
					Name:   "deploy",
					Source: "github.example.com/foo/bar/deploy_template.yml",
					Type:   "github",
				},
			},
		},
	}

	deployCfg := yaml.Deployment{
		Template: yaml.StepTemplate{
			Name: "deploy",
			Variables: map[string]interface{}{
				"regions": []string{"us-east-1", "us-west-1"},
			},
		},
	}

	wantDeployCfg := yaml.Deployment{
		Targets: []string{"dev", "prod", "stage"},
		Parameters: yaml.ParameterMap{
			"region": {
				Description: "cluster region to deploy",
				Type:        "string",
				Required:    true,
				Options:     []string{"us-east-1", "us-west-1"},
			},
			"cluster_count": {
				Description: "number of clusters to deploy to",
				Type:        "integer",
				Required:    false,
				Min:         1,
				Max:         10,
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithCommit("123abc456def").WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			build, err := compiler.ExpandDeployment(
				context.Background(),
				&yaml.Build{
					Deployment: deployCfg,
				},
				test.tmpls)
			if err != nil {
				t.Errorf("ExpandDeployment_Type%s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(wantDeployCfg, build.Deployment); diff != "" {
				t.Errorf("ExpandDeployment()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}

func TestNative_ExpandStepsMulti(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	tmpls := map[string]*yaml.Template{
		"gradle": {
			Name:   "gradle",
			Source: "github.example.com/foo/bar/gradle.yml",
			Type:   "github",
		},
		"maven": {
			Name:   "maven",
			Source: "github.example.com/bar/foo/maven.yml",
			Type:   "github",
		},
		"npm": {
			Name:   "npm",
			Source: "github.example.com/foo/bar/gradle.yml",
			Type:   "github",
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "gradle",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
		},
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "maven",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
			Ruleset: yaml.Ruleset{
				If: yaml.Rules{
					Branch: []string{"main"},
				},
				Operator: "and",
			},
		},
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "npm",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
			Ruleset: yaml.Ruleset{
				If: yaml.Rules{
					Branch: []string{"dev"},
				},
				Operator: "and",
			},
		},
	}

	wantSteps := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_build",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"mvn downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"mvn check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"mvn build"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_build",
			Pull:  "always",
		},
	}

	wantSecrets := yaml.SecretSlice{
		&yaml.Secret{
			Name:   "docker_username",
			Key:    "org/repo/foo/bar",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "foo_password",
			Key:    "org/repo/foo/password",
			Engine: "vault",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "vault_token",
			Key:    "vault_token",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Origin: yaml.Origin{
				Name:  "private vault",
				Image: "target/secret-vault:latest",
				Pull:  "always",
				Secrets: yaml.StepSecretSlice{
					{
						Source: "vault_token",
						Target: "VAULT_TOKEN",
					},
				},
				Parameters: map[string]interface{}{
					"addr":        "vault.example.com",
					"auth_method": "token",
					"username":    "octocat",
					"items": []interface{}{
						map[string]interface{}{"path": "docker", "source": "secret/docker"},
					},
				},
			},
		},
		&yaml.Secret{
			Origin: yaml.Origin{
				Name:  "private vault",
				Image: "target/secret-vault:latest",
				Pull:  "always",
				Secrets: yaml.StepSecretSlice{
					{
						Source: "vault_token",
						Target: "VAULT_TOKEN",
					},
				},
				Parameters: map[string]interface{}{
					"addr":        "vault.example.com",
					"auth_method": "token",
					"username":    "octocat",
					"items": []interface{}{
						map[string]interface{}{"path": "docker", "source": "secret/docker"},
					},
				},
			},
		},
	}

	wantServices := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres:12",
			Name:  "postgres",
			Pull:  "not_present",
		},
	}

	wantEnvironment := raw.StringSliceMap{}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	ruledata := new(pipeline.RuleData)
	ruledata.Branch = "main"

	build, _, err := compiler.ExpandSteps(context.Background(),
		&yaml.Build{
			Steps:       steps,
			Services:    yaml.ServiceSlice{},
			Environment: raw.StringSliceMap{},
		},
		tmpls, ruledata, nil, compiler.GetTemplateDepth())
	if err != nil {
		t.Errorf("ExpandSteps returned err: %v", err)
	}

	if diff := cmp.Diff(build.Steps, wantSteps); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Services, wantServices); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_ExpandStepsStarlark(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	tmpls := map[string]*yaml.Template{
		"go": {
			Name:   "go",
			Source: "github.example.com/foo/bar/template.star",
			Format: "starlark",
			Type:   "github",
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name:      "go",
				Variables: map[string]interface{}{},
			},
		},
	}

	wantSteps := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"go build", "go test"},
			Image:    "golang:latest",
			Name:     "sample_build",
			Pull:     "not_present",
		},
	}

	wantSecrets := yaml.SecretSlice{}
	wantServices := yaml.ServiceSlice{}
	wantEnvironment := raw.StringSliceMap{
		"star": "test3",
		"bar":  "test4",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	build, _, err := compiler.ExpandSteps(context.Background(),
		&yaml.Build{
			Steps:       steps,
			Secrets:     yaml.SecretSlice{},
			Services:    yaml.ServiceSlice{},
			Environment: raw.StringSliceMap{},
		},
		tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
	if err != nil {
		t.Errorf("ExpandSteps returned err: %v", err)
	}

	if diff := cmp.Diff(build.Steps, wantSteps); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Services, wantServices); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}

	if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
		t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
	}
}

func TestNative_ExpandSteps_TemplateCallTemplate(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testBuild := new(api.Build)

	testBuild.SetID(1)
	testBuild.SetCommit("123abc456def")

	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "Test 1",
			tmpls: map[string]*yaml.Template{
				"chain": {
					Name:   "chain",
					Source: "github.example.com/faz/baz/template_calls_template.yml",
					Type:   "github",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "chain",
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	wantSteps := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_call template_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_call template_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build", "echo test"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_call template_build",
			Pull:  "always",
		},
	}

	wantSecrets := yaml.SecretSlice{
		&yaml.Secret{
			Name:   "docker_username",
			Key:    "org/repo/foo/bar",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "foo_password",
			Key:    "org/repo/foo/password",
			Engine: "vault",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
	}

	wantServices := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres:12",
			Name:  "postgres",
			Pull:  "not_present",
		},
	}

	wantEnvironment := raw.StringSliceMap{
		"foo":  "test1",
		"bar":  "test2",
		"star": "test3",
	}

	wantTemplates := yaml.TemplateSlice{
		{
			Name:   "chain",
			Source: "github.example.com/faz/baz/template_calls_template.yml",
			Type:   "github",
		},
		{
			Name:   "test",
			Source: "github.example.com/foo/bar/long_template.yml",
			Type:   "github",
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithBuild(testBuild).WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			build, _, err := compiler.ExpandSteps(context.Background(),
				&yaml.Build{
					Steps: steps, Services: yaml.ServiceSlice{},
					Environment: globalEnvironment,
					Templates:   yaml.TemplateSlice{test.tmpls["chain"]},
				},
				test.tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
			if err != nil {
				t.Errorf("ExpandSteps_Type%s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(build.Steps, wantSteps); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Services, wantServices); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Templates, wantTemplates); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}

func TestNative_ExpandStepsDuplicateCalls(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	testCallsMap := make(map[string]bool)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		testCallKey := c.Param("path")

		if refQuery, exists := c.GetQuery("ref"); exists {
			testCallKey += refQuery
		}

		// this is the real test
		if testCallsMap[testCallKey] {
			t.Errorf("ExpandSteps() called the same template %s twice", c.Param("path"))
		}

		testCallsMap[c.Param("path")] = true
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "GitHub",
			tmpls: map[string]*yaml.Template{
				"gradle": {
					Name:   "gradle",
					Source: "github.example.com/foo/bar/long_template.yml",
					Type:   "github",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "gradle",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
		},
		&yaml.Step{
			Name: "sample-dup",
			Template: yaml.StepTemplate{
				Name: "gradle",
				Variables: map[string]interface{}{
					"image":       "openjdk:latest",
					"environment": "{ GRADLE_USER_HOME: .gradle, GRADLE_OPTS: -Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false }",
					"pull_policy": "pull: true",
				},
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	wantSteps := yaml.StepSlice{
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build", "echo gradle"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample_build",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample-dup_install",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample-dup_test",
			Pull:  "always",
		},
		&yaml.Step{
			Commands: []string{"./gradlew build", "echo gradle"},
			Environment: raw.StringSliceMap{
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
				"GRADLE_USER_HOME": ".gradle",
			},
			Image: "openjdk:latest",
			Name:  "sample-dup_build",
			Pull:  "always",
		},
	}

	wantSecrets := yaml.SecretSlice{
		&yaml.Secret{
			Name:   "docker_username",
			Key:    "org/repo/foo/bar",
			Engine: "native",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
		&yaml.Secret{
			Name:   "foo_password",
			Key:    "org/repo/foo/password",
			Engine: "vault",
			Type:   "repo",
			Origin: yaml.Origin{},
			Pull:   "build_start",
		},
	}

	wantServices := yaml.ServiceSlice{
		&yaml.Service{
			Image: "postgres:12",
			Name:  "postgres",
			Pull:  "not_present",
		},
	}

	wantEnvironment := raw.StringSliceMap{
		"foo":  "test1",
		"bar":  "test2",
		"star": "test3",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithCommit("123abc456def").WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			build, _, err := compiler.ExpandSteps(
				context.Background(),
				&yaml.Build{
					Steps:       steps,
					Services:    yaml.ServiceSlice{},
					Environment: globalEnvironment,
				},
				test.tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
			if err != nil {
				t.Errorf("ExpandSteps_Type%s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(build.Steps, wantSteps); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Secrets, wantSecrets); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Services, wantServices); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}

			if diff := cmp.Diff(build.Environment, wantEnvironment); diff != "" {
				t.Errorf("ExpandSteps()_Type%s mismatch (-want +got):\n%s", test.name, diff)
			}
		})
	}
}

func TestNative_ExpandSteps_TemplateCallTemplate_CircularFail(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testBuild := new(api.Build)

	testBuild.SetID(1)
	testBuild.SetCommit("123abc456def")

	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "Test 1",
			tmpls: map[string]*yaml.Template{
				"circle": {
					Name:   "circle",
					Source: "github.example.com/bad/design/template_calls_itself.yml",
					Type:   "github",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "circle",
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithBuild(testBuild).WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := compiler.ExpandSteps(context.Background(),
				&yaml.Build{
					Steps: steps, Services: yaml.ServiceSlice{}, Environment: globalEnvironment,
				},
				test.tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
			if err == nil {
				t.Errorf("ExpandSteps_Type%s should have returned an error", test.name)
			}
		})
	}
}

func TestNative_ExpandSteps_CallTemplateWithRenderInline(t *testing.T) {
	// setup context
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	_, engine := gin.CreateTestContext(resp)

	// setup mock server
	engine.GET("/api/v3/repos/:org/:repo/contents/:path", func(c *gin.Context) {
		body, err := convertFileToGithubResponse(c.Param("path"))
		if err != nil {
			t.Error(err)
		}
		c.JSON(http.StatusOK, body)
	})

	s := httptest.NewServer(engine)
	defer s.Close()

	// setup types
	testBuild := new(api.Build)

	testBuild.SetID(1)
	testBuild.SetCommit("123abc456def")

	testRepo := new(api.Repo)

	testRepo.SetID(1)
	testRepo.SetOrg("foo")
	testRepo.SetName("bar")

	tests := []struct {
		name  string
		tmpls map[string]*yaml.Template
	}{
		{
			name: "Test 1",
			tmpls: map[string]*yaml.Template{
				"render_inline": {
					Name:   "render_inline",
					Source: "github.example.com/github/octocat/nested.yml",
					Type:   "github",
				},
			},
		},
	}

	steps := yaml.StepSlice{
		&yaml.Step{
			Name: "sample",
			Template: yaml.StepTemplate{
				Name: "render_inline",
			},
		},
	}

	globalEnvironment := raw.StringSliceMap{
		"foo": "test1",
		"bar": "test2",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, s.URL))
	if err != nil {
		t.Errorf("Creating new compiler returned err: %v", err)
	}

	compiler.WithBuild(testBuild).WithRepo(testRepo)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := compiler.ExpandSteps(context.Background(),
				&yaml.Build{
					Steps:       steps,
					Services:    yaml.ServiceSlice{},
					Environment: globalEnvironment,
				},
				test.tmpls, new(pipeline.RuleData), nil, compiler.GetTemplateDepth())
			if err == nil {
				t.Errorf("ExpandSteps_Type%s should have returned an error", test.name)
			}
		})
	}
}

func TestNative_mapFromTemplates(t *testing.T) {
	// setup types
	str := "foo"

	tmpl := []*yaml.Template{
		{
			Name:   str,
			Source: str,
			Type:   str,
		},
	}

	want := map[string]*yaml.Template{
		str: {
			Name:   str,
			Source: str,
			Type:   str,
		},
	}

	// run test
	got := mapFromTemplates(tmpl)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapFromTemplates is %v, want %v", got, want)
	}
}
