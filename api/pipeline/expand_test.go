// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"go.yaml.in/yaml/v3"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/native"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal"
	pipelineMiddleware "github.com/go-vela/server/router/middleware/pipeline"
	repoMiddleware "github.com/go-vela/server/router/middleware/repo"
	userMiddleware "github.com/go-vela/server/router/middleware/user"
)

func TestExpandPipelineStagesReturnsCleanYAML(t *testing.T) {
	gin.SetMode(gin.TestMode)

	assertCleanStages := func(t *testing.T, w *httptest.ResponseRecorder) {
		t.Helper()

		if w.Code != http.StatusOK {
			t.Fatalf("ExpandPipeline returned status %d, want %d", w.Code, http.StatusOK)
		}

		var got struct {
			Stages map[string]any `yaml:"stages"`
		}

		if err := yaml.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatalf("unable to parse yaml response: %v", err)
		}

		if len(got.Stages) == 0 {
			t.Fatalf("expected stages in response, got none")
		}

		if _, ok := got.Stages["build"]; !ok {
			t.Fatalf("expected build stage in response, got keys %#v", got.Stages)
		}

		if _, hasKind := got.Stages["kind"]; hasKind {
			t.Fatalf("unexpected raw yaml node output: %#v", got.Stages)
		}
	}

	setup := func(t *testing.T, requestPath string) (*httptest.ResponseRecorder, *gin.Context) {
		t.Helper()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodPost, requestPath, nil)
		c.Request = req

		logger := logrus.New()
		logger.SetOutput(io.Discard)
		c.Set("logger", logrus.NewEntry(logger))

		c.Set("metadata", &internal.Metadata{})

		repo := new(api.Repo)
		repo.SetOrg("foo")
		repo.SetName("bar")
		repo.SetFullName("foo/bar")
		repo.SetPipelineType(constants.PipelineTypeYAML)
		repoMiddleware.ToContext(c, repo)

		pipeline := new(api.Pipeline)
		pipeline.SetCommit("123")
		pipeline.SetType(constants.PipelineTypeYAML)
		pipeline.SetData([]byte(`version: "1"
stages:
  build:
    steps:
      - name: test
        image: alpine
        commands:
          - echo hello
`))
		pipelineMiddleware.ToContext(c, pipeline)

		user := new(api.User)
		userMiddleware.ToContext(c, user)

		engine := newTestCompiler(t)
		compiler.WithGinContext(c, engine)

		return w, c
	}

	t.Run("explicit yaml", func(t *testing.T) {
		w, c := setup(t, "/api/v1/pipelines/foo/bar/123/expand?output=yaml")

		ExpandPipeline(c)

		assertCleanStages(t, w)
	})

	t.Run("default yaml", func(t *testing.T) {
		w, c := setup(t, "/api/v1/pipelines/foo/bar/123/expand")

		ExpandPipeline(c)

		assertCleanStages(t, w)
	})
}

func newTestCompiler(t *testing.T) compiler.Engine {
	t.Helper()

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "clone-image",
				Value: "target/vela-git:latest",
			},
			&cli.IntFlag{
				Name:  "max-template-depth",
				Value: 5,
			},
			&cli.BoolFlag{
				Name:  "github-driver",
				Value: false,
			},
			&cli.StringFlag{
				Name:  "github-url",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "github-token",
				Value: "",
			},
			&cli.Int64Flag{
				Name:  "compiler-starlark-exec-limit",
				Value: 0,
			},
		},
	}

	engine, err := native.FromCLICommand(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unable to create compiler: %v", err)
	}

	return engine
}
