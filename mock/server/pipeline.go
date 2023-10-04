// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/yaml"

	yml "github.com/buildkite/yaml"
)

const (
	// CompileResp represents a YAML return for an compiled pipeline.
	CompileResp = `---
version: "1"

secrets:
  - name: docker_username
    key: go-vela/docker/username
    engine: native
    type: org

  - name: docker_password
    key: go-vela/docker/password
    engine: native
    type: org

steps:
  - name: go_test
    image: golang:latest
    environment:
      CGO_ENABLED: "0"
      GOOS: linux
    commands:
      - go test ./...

  - name: go_build
    image: golang:latest
    environment:
      CGO_ENABLED: "0"
      GOOS: linux
    commands:
      - go build

  - name: non-template-echo
    image: golang:latest
    commands:
      - echo hello

templates:
  - name: sample
    source: github.com/go-vela/vela-tutorials/templates/sample.yml
    type: github
`

	// ExpandResp represents a YAML return for an expanded pipeline.
	ExpandResp = `---
version: "1"

secrets:
  - name: docker_username
    key: go-vela/docker/username
    engine: native
    type: org

  - name: docker_password
    key: go-vela/docker/password
    engine: native
    type: org

steps:
  - name: go_test
    image: golang:latest
    environment:
      CGO_ENABLED: "0"
      GOOS: linux
    commands:
      - go test ./...

  - name: go_build
    image: golang:latest
    environment:
      CGO_ENABLED: "0"
      GOOS: linux
    commands:
      - go build

  - name: non-template-echo
    image: golang:latest
    commands:
      - echo hello

templates:
  - name: sample
    source: github.com/go-vela/vela-tutorials/templates/sample.yml
    type: github
`
	// PipelineResp represents a JSON return for a single pipeline.
	PipelineResp = `{
  "id": 1,
  "repo_id": 1,
  "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
  "flavor": "",
  "platform": "",
  "ref": "refs/heads/master",
  "type": "yaml",
  "version": "1",
  "external_secrets": false,
  "internal_secrets": false,
  "services": false,
  "stages": false,
  "steps": true,
  "templates": false,
  "data": "LS0tCnZlcnNpb246ICIxIgoKc3RlcHM6CiAgLSBuYW1lOiBlY2hvCiAgICBpbWFnZTogYWxwaW5lOmxhdGVzdAogICAgY29tbWFuZHM6IFtlY2hvIGZvb10="
}`

	// PipelinesResp represents a JSON return for one to many hooks.
	PipelinesResp = `[
  {
    "id": 2
    "repo_id": 1,
    "commit": "a49aaf4afae6431a79239c95247a2b169fd9f067",
    "flavor": "",
    "platform": "",
    "ref": "refs/heads/master",
    "type": "yaml",
    "version": "1",
    "external_secrets": false,
    "internal_secrets": false,
    "services": false,
    "stages": false,
    "steps": true,
    "templates": false,
    "data": "LS0tCnZlcnNpb246ICIxIgoKc3RlcHM6CiAgLSBuYW1lOiBlY2hvCiAgICBpbWFnZTogYWxwaW5lOmxhdGVzdAogICAgY29tbWFuZHM6IFtlY2hvIGZvb10="
  },
  {
    "id": 1,
    "repo_id": 1,
    "commit": "48afb5bdc41ad69bf22588491333f7cf71135163",
    "flavor": "",
    "platform": "",
    "ref": "refs/heads/master",
    "type": "yaml",
    "version": "1",
    "external_secrets": false,
    "internal_secrets": false,
    "services": false,
    "stages": false,
    "steps": true,
    "templates": false,
    "data": "LS0tCnZlcnNpb246ICIxIgoKc3RlcHM6CiAgLSBuYW1lOiBlY2hvCiAgICBpbWFnZTogYWxwaW5lOmxhdGVzdAogICAgY29tbWFuZHM6IFtlY2hvIGZvb10="
  }
]`

	// TemplateResp represents a YAML return for templates in a pipeline.
	TemplateResp = `---
sample:
  name: sample
  source: github.com/go-vela/vela-tutorials/templates/sample.yml
  type: github
`
)

// getPipelines returns mock JSON for a http GET.
func getPipelines(c *gin.Context) {
	data := []byte(PipelinesResp)

	var body []library.Pipeline
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getPipeline has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func getPipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(PipelineResp)

	var body library.Pipeline
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addPipeline returns mock JSON for a http POST.
func addPipeline(c *gin.Context) {
	data := []byte(PipelineResp)

	var body library.Pipeline
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updatePipeline has a param :pipeline returns mock JSON for a http PUT.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func updatePipeline(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		p := c.Param("pipeline")

		if strings.EqualFold(p, "0") {
			msg := fmt.Sprintf("Pipeline %s does not exist", p)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(PipelineResp)

	var body library.Pipeline
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removePipeline has a param :pipeline returns mock JSON for a http DELETE.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func removePipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Pipeline %s removed", p))
}

// compilePipeline has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func compilePipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(CompileResp)

	var body yaml.Build
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// expandPipeline has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func expandPipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(ExpandResp)

	var body yaml.Build
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// getTemplates has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func getTemplates(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(TemplateResp)

	body := make(map[string]*yaml.Template)
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// validatePipeline has a param :pipeline returns mock YAML for a http GET.
//
// Pass "0" to :pipeline to test receiving a http 404 response.
func validatePipeline(c *gin.Context) {
	p := c.Param("pipeline")

	if strings.EqualFold(p, "0") {
		msg := fmt.Sprintf("Pipeline %s does not exist", p)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, "pipeline is valid")
}
