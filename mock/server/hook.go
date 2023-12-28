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
	// HookResp represents a JSON return for a single hook.
	HookResp = `{
  "id": 1,
  "repo_id": 1,
  "build_id": 1,
  "number": 1,
  "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
  "created": 1563475419,
  "host": "github.com",
  "event": "push",
  "branch": "main",
  "error": "",
  "status": "success",
  "link": "https://github.com/github/octocat/settings/hooks/1"
}`

	// HooksResp represents a JSON return for one to many hooks.
	HooksResp = `[
  {
    "id": 2,
    "repo_id": 1,
    "build_id": 1,
    "number": 2,
    "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
    "created": 1563475420,
    "host": "github.com",
    "event": "push",
    "branch": "main",
    "error": "",
    "status": "success",
    "link": "https://github.com/github/octocat/settings/hooks/1"
  },
  {
    "id": 1,
    "repo_id": 1,
    "build_id": 1,
    "number": 1,
    "source_id": "c8da1302-07d6-11ea-882f-4893bca275b8",
    "created": 1563475419,
    "host": "github.com",
    "event": "push",
    "branch": "main",
    "error": "",
    "status": "success",
    "link": "https://github.com/github/octocat/settings/hooks/1"
  }
]`
)

// getHooks returns mock JSON for a http GET.
func getHooks(c *gin.Context) {
	data := []byte(HooksResp)

	var body []library.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// getHook has a param :hook returns mock JSON for a http GET.
//
// Pass "0" to :hook to test receiving a http 404 response.
func getHook(c *gin.Context) {
	s := c.Param("hook")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Hook %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	data := []byte(HookResp)

	var body library.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// addHook returns mock JSON for a http POST.
func addHook(c *gin.Context) {
	data := []byte(HookResp)

	var body library.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusCreated, body)
}

// updateHook has a param :hook returns mock JSON for a http PUT.
//
// Pass "0" to :hook to test receiving a http 404 response.
func updateHook(c *gin.Context) {
	if !strings.Contains(c.FullPath(), "admin") {
		s := c.Param("hook")

		if strings.EqualFold(s, "0") {
			msg := fmt.Sprintf("Hook %s does not exist", s)

			c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

			return
		}
	}

	data := []byte(HookResp)

	var body library.Hook
	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// removeHook has a param :hook returns mock JSON for a http DELETE.
//
// Pass "0" to :hook to test receiving a http 404 response.
func removeHook(c *gin.Context) {
	s := c.Param("hook")

	if strings.EqualFold(s, "0") {
		msg := fmt.Sprintf("Hook %s does not exist", s)

		c.AbortWithStatusJSON(http.StatusNotFound, types.Error{Message: &msg})

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("Hook %s removed", s))
}
