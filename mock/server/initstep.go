// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore duplicate with user code
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

const (
	// InitStepResp represents a JSON return for a single step.
	InitStepResp = `{
  "id": 1,
  "repo_id": 1,
  "build_id": 1,
  "number": 1,
  "reporter": "Foobar Runtime",
  "name": "foobar",
  "mimetype": "text/plain"
}`

	// InitStepsResp represents a JSON return for one to many steps.
	InitStepsResp = `[
  {
    "id": 2,
    "repo_id": 1,
    "build_id": 1,
    "number": 2,
    "reporter": "Foobar Runtime",
    "name": "foobar",
    "mimetype": "text/plain"
  },
  {
    "id": 1,
    "repo_id": 1,
    "build_id": 1,
    "number": 1,
    "reporter": "Foobar Runtime",
    "name": "foobar",
    "mimetype": "text/plain"
  }
]`
)

// getInitSteps returns mock JSON for a http GET.
func getInitSteps(c *gin.Context) {
	data := []byte(InitStepsResp)

	var body []library.InitStep
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getInitStep has a param :initstep returns mock JSON for a http GET.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func getInitStep(c *gin.Context) {
	i := c.Param("initstep")

	if strings.EqualFold(i, "0") {
		msg := fmt.Sprintf("InitStep %s does not exist", i)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(InitStepResp)

	var body library.InitStep
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addInitStep returns mock JSON for a http POST.
func addInitStep(c *gin.Context) {
	data := []byte(InitStepResp)

	var body library.InitStep
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateInitStep has a param :initstep returns mock JSON for a http PUT.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func updateInitStep(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		i := c.Param("initstep")

		if strings.EqualFold(i, "0") {
			msg := fmt.Sprintf("InitStep %s does not exist", i)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(InitStepResp)

	var body library.InitStep
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeInitStep has a param :initstep returns mock JSON for a http DELETE.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func removeInitStep(c *gin.Context) {
	i := c.Param("initstep")

	if strings.EqualFold(i, "0") {
		msg := fmt.Sprintf("InitStep %s does not exist", i)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("InitStep %s removed", i))
}
