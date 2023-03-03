// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
	// LogResp represents a JSON return for a single log.
	LogResp = `{
  "id": 1,
  "build_id": 1,
  "repo_id": 1,
  "service_id": 1,
  "step_id": 1,
  "init_step_id": 1,
  "data": "SGVsbG8sIFdvcmxkIQ=="
}`
)

// getServiceLog has a param :service returns mock JSON for a http GET.
//
// Pass "0" to :step to test receiving a http 404 response.
func getServiceLog(c *gin.Context) {
	s := c.Param("service")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addServiceLog returns mock JSON for a http GET.
func addServiceLog(c *gin.Context) {
	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateServiceLog has a param :service returns mock JSON for a http PUT.
//
// Pass "0" to :step to test receiving a http 404 response.
func updateServiceLog(c *gin.Context) {
	s := c.Param("service")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeServiceLog has a param :service returns mock JSON for a http DELETE.
//
// Pass "0" to :step to test receiving a http 404 response.
func removeServiceLog(c *gin.Context) {
	s := c.Param("service")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Log %s removed", s))
}

// getStepLog has a param :step returns mock JSON for a http GET.
//
// Pass "0" to :step to test receiving a http 404 response.
func getStepLog(c *gin.Context) {
	s := c.Param("step")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addStepLog returns mock JSON for a http GET.
func addStepLog(c *gin.Context) {
	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateStepLog has a param :step returns mock JSON for a http PUT.
//
// Pass "0" to :step to test receiving a http 404 response.
func updateStepLog(c *gin.Context) {
	s := c.Param("step")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeStepLog has a param :step returns mock JSON for a http DELETE.
//
// Pass "0" to :step to test receiving a http 404 response.
func removeStepLog(c *gin.Context) {
	s := c.Param("step")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Log %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Log %s removed", s))
}

// getInitStepLog has a param :initstep returns mock JSON for a http GET.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func getInitStepLog(c *gin.Context) {
	i := c.Param("initstep")

	if strings.EqualFold(i, "0") {
		msg := fmt.Sprintf("Log %s does not exist", i)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addInitStepLog returns mock JSON for a http GET.
func addInitStepLog(c *gin.Context) {
	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateInitStepLog has a param :initstep returns mock JSON for a http PUT.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func updateInitStepLog(c *gin.Context) {
	i := c.Param("initstep")

	if strings.EqualFold(i, "0") {
		msg := fmt.Sprintf("Log %s does not exist", i)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(LogResp)

	var body library.Log
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeInitStepLog has a param :initstep returns mock JSON for a http DELETE.
//
// Pass "0" to :initstep to test receiving a http 404 response.
func removeInitStepLog(c *gin.Context) {
	i := c.Param("initstep")

	if strings.EqualFold(i, "0") {
		msg := fmt.Sprintf("Log %s does not exist", i)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Log %s removed", i))
}
