// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package server

import (
	"fmt"
	"net/http"
	"strings"

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

	// PipelineResp represents a YAML return for a single pipeline.
	PipelineResp = `---
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
  - name: go
    template:
      name: sample

  - name: non-template-echo
    image: golang:latest
    commands:
      - echo hello

templates:
  - name: sample
    source: github.com/go-vela/vela-tutorials/templates/sample.yml
    type: github
`

	// TemplateResp represents a YAML return for templates in a pipeline.
	TemplateResp = `---
sample:
  name: sample
  source: github.com/go-vela/vela-tutorials/templates/sample.yml
  type: github
`
)

// getPipeline has a param :repo returns mock YAML for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func getPipeline(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(PipelineResp)

	var body yaml.Build
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// compilePipeline has a param :repo returns mock YAML for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func compilePipeline(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(CompileResp)

	var body yaml.Build
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// expandPipeline has a param :repo returns mock YAML for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func expandPipeline(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(ExpandResp)

	var body yaml.Build
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// getTemplates has a param :repo returns mock YAML for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func getTemplates(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(TemplateResp)

	body := make(map[string]*yaml.Template)
	_ = yml.Unmarshal(data, &body)

	c.YAML(http.StatusOK, body)
}

// validatePipeline has a param :repo returns mock YAML for a http GET.
//
// Pass "not-found" to :repo to test receiving a http 404 response.
func validatePipeline(c *gin.Context) {
	r := c.Param("repo")

	if strings.Contains(r, "not-found") {
		msg := fmt.Sprintf("Repo %s does not exist", r)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, "pipeline is valid")
}
