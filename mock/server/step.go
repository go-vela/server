// SPDX-License-Identifier: Apache-2.0

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
	// StepResp represents a JSON return for a single step.
	StepResp = `{
  "id": 1,
  "build_id": 1,
  "repo_id": 1,
  "number": 1,
  "name": "clone",
  "status": "success",
  "error": "",
  "exit_code": 0,
  "created": 1563475419,
  "started": 0,
  "finished": 0,
  "host": "host.company.com",
  "runtime": "docker",
  "distribution": "linux"
}`

	// StepsResp represents a JSON return for one to many steps.
	StepsResp = `[
  {
    "id": 2,
    "build_id": 1,
    "repo_id": 1,
    "number": 2,
    "name": "build",
    "status": "success",
    "error": "",
    "exit_code": 0,
    "created": 1563475419,
    "started": 0,
    "finished": 0,
    "host": "host.company.com",
    "runtime": "docker",
    "distribution": "linux"
  },
  {
    "id": 1,
    "build_id": 1,
    "repo_id": 1,
    "number": 1,
    "name": "clone",
    "status": "success",
    "error": "",
    "exit_code": 0,
    "created": 1563475419,
    "started": 0,
    "finished": 0,
    "host": "host.company.com",
    "runtime": "docker",
    "distribution": "linux"
  }
]`
)

// getSteps returns mock JSON for a http GET.
func getSteps(c *gin.Context) {
	data := []byte(StepsResp)

	var body []library.Step
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getStep has a param :step returns mock JSON for a http GET.
//
// Pass "0" to :step to test receiving a http 404 response.
func getStep(c *gin.Context) {
	s := c.Param("step")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Step %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(StepResp)

	var body library.Step
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addStep returns mock JSON for a http POST.
func addStep(c *gin.Context) {
	data := []byte(StepResp)

	var body library.Step
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateStep has a param :step returns mock JSON for a http PUT.
//
// Pass "0" to :step to test receiving a http 404 response.
func updateStep(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		s := c.Param("step")

		if strings.EqualFold(s, "0") {
			msg := fmt.Sprintf("Step %s does not exist", s)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(StepResp)

	var body library.Step
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeStep has a param :step returns mock JSON for a http DELETE.
//
// Pass "0" to :step to test receiving a http 404 response.
func removeStep(c *gin.Context) {
	s := c.Param("step")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Step %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Step %s removed", s))
}
