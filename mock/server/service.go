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
	// ServiceResp represents a JSON return for a single service.
	ServiceResp = `{
  "id": 1,
  "build_id": 1,
  "repo_id": 1,
  "number": 1,
  "name": "clone",
  "status": "success",
  "error": "",
  "exit_code": 0,
  "created": 1563475419,
  "started": 1563475420,
  "finished": 1563475421
}`

	// ServicesResp represents a JSON return for one to many services.
	ServicesResp = `[
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
    "started": 1563475420,
    "finished": 1563475421
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
    "started": 1563475420,
    "finished": 1563475421
  }
]`
)

// getServices returns mock JSON for a http GET.
func getServices(c *gin.Context) {
	data := []byte(ServicesResp)

	var body []library.Service
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getService has a param :service returns mock JSON for a http GET.
//
// Pass "0" to :service to test receiving a http 404 response.
func getService(c *gin.Context) {
	s := c.Param("service")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Service %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(ServiceResp)

	var body library.Service
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addService returns mock JSON for a http POST.
func addService(c *gin.Context) {
	data := []byte(ServiceResp)

	var body library.Service
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateService has a param :service returns mock JSON for a http PUT.
//
// Pass "0" to :service to test receiving a http 404 response.
func updateService(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		s := c.Param("service")

		if strings.EqualFold(s, "0") {
			msg := fmt.Sprintf("Service %s does not exist", s)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(ServiceResp)

	var body library.Service
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeService has a param :service returns mock JSON for a http DELETE.
//
// Pass "0" to :service to test receiving a http 404 response.
func removeService(c *gin.Context) {
	s := c.Param("service")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Service %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Service %s removed", s))
}
